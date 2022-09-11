package watcher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

type (
	watcher struct {
		*docker.Client
	}
	ContainerHandler  func(context.Context, *docker.Client, types.Container) error
	containerHandlers []ContainerHandler
)

func NewWatcher() (*watcher, error) {
	cli, err := docker.NewClientWithOpts()
	if err != nil {
		return nil, err
	}

	return &watcher{cli}, nil
}

func (w *watcher) WatchContainers(ctx context.Context, interval time.Duration, handlers ...ContainerHandler) {
	w.watchContainers(ctx, interval, handlers)
}

func (w *watcher) WatchContainersEvents(ctx context.Context, handlers ...func(events.Message) error) {
	if len(handlers) == 0 {
		return
	}
	eventsChan, errChan := w.Events(ctx, types.EventsOptions{
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
			err = handlers.handle(ctx, w.Client, container)
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
	go func() {
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
	containers, err := w.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	return containers, err
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
