package watcher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	senderpkg "docker_data_collector/internal/sender"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

type (
	sender interface {
		Send(string) error
	}

	watcher struct {
		c *docker.Client
		s sender

		wg sync.WaitGroup
	}

	ContainerHandler func(context.Context, *docker.Client, types.Container) error

	containerHandlers []ContainerHandler

	statsEntry struct {
		Container        string
		Name             string
		ID               string
		NetworkIO        string
		CPUPercentage    float64
		Memory           float64
		MemoryLimit      float64
		MemoryPercentage float64
		NetworkRx        float64
		NetworkTx        float64
		BlockRead        float64
		BlockWrite       float64
		PidsCurrent      uint64 // Not used on Windows
		IsInvalid        bool
	}
)

func newWatcher(s senderpkg.Sender) (*watcher, error) {
	cli, err := docker.NewClientWithOpts()
	if err != nil {
		return nil, err
	}

	return &watcher{
		c: cli,
		s: s,
	}, nil
}

func (w *watcher) Run(ctx context.Context) error {
	w.watchContainers(ctx, time.Hour*24, containerHandlers{
		w.statsHandler(),
		w.createdHandler(),
	})

	return nil
}

func (w *watcher) Stop() error {
	w.wg.Wait()
	return nil
}

func (w *watcher) watchContainersEvents(ctx context.Context, handlers ...func(events.Message) error) {
	if len(handlers) == 0 {
		return
	}
	eventsChan, errChan := w.c.Events(ctx, types.EventsOptions{
		Filters: filters.NewArgs(filters.Arg("type", events.ContainerEventType)),
	})

	started := make(chan struct{})

	go func() {
		close(started)
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-eventsChan:
				for _, handler := range handlers {
					hErr := handler(e)
					if hErr != nil {
						logrus.WithError(hErr).Error("handle event")
						continue
					}
				}
			case err := <-errChan:
				if !errors.Is(err, io.EOF) {
					return
				}
				logrus.WithError(err).Error("handle err chan")
			}
		}
	}()
	<-started
}

func (w *watcher) watchContainers(ctx context.Context, interval time.Duration, handlers containerHandlers) {
	f := func() error {
		containers, err := w.containersAll(ctx)
		if err != nil {
			return fmt.Errorf("containers all: %w", err)
		}

		for _, container := range containers {
			err = handlers.handle(ctx, w.c, container)
			if err != nil {
				logrus.WithError(err).Error("handle")
				continue
			}
		}
		return nil
	}

	err := f()
	if err != nil {
		logrus.WithError(err).Error("initial run")
	}

	logrus.WithField("interval", interval.String()).Debug("starting watch containers")

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			case <-time.After(interval):
			}
			err := f()
			if err != nil {
				logrus.WithError(err).Error("tick")
			}
		}
	}()
}

func (w *watcher) containersAll(ctx context.Context) ([]types.Container, error) {
	return w.c.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
}

func (c containerHandlers) handle(ctx context.Context, cli *docker.Client, container types.Container) (errors error) {
	for _, handler := range c {
		err := handler(ctx, cli, container)
		if err != nil {
			errors = multierr.Append(errors, err)
		}
	}
	return
}
