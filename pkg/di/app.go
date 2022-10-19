package di

import (
	"context"
	"sync"

	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

type (
	Runner interface {
		Run(ctx context.Context) error
	}

	Stopper interface {
		Stop(context.Context) error
	}

	Application interface {
		Runner
		Stopper
	}

	ApplicationPool []Application
)

func (a ApplicationPool) Run(ctx context.Context) error {
	pb := progressbar.Default(int64(len(a)), "starting applications...")
	defer func() {
		err := pb.Close()
		if err != nil {
			logrus.WithError(err).Error("close progress bar")
		}
	}()
	for _, application := range a {
		logrus.Infof("run %T", application)
		if err := application.Run(ctx); err != nil {
			logrus.WithError(err).Errorf("%T failed to run", application)
			return err
		}
		_ = pb.Add(1)
	}
	logrus.Info("all applications started")
	return nil
}

func (a ApplicationPool) Stop(ctx context.Context) error {
	wg := sync.WaitGroup{}
	done := make(chan struct{})
	wg.Add(len(a))
	go func() {
		wg.Wait()
		close(done)
	}()

	errorList := make([]error, 0, len(a))
	mu := sync.Mutex{}

	pb := progressbar.Default(int64(len(a)), "stopping applications...")
	defer func() {
		err := pb.Close()
		if err != nil {
			logrus.WithError(err).Error("close progress bar")
		}
	}()
	for _, application := range a {
		application := application
		go func() {
			defer wg.Done()
			if err := application.Stop(ctx); err != nil {
				logrus.WithError(err).Errorf("%T failed to stop", application)
				mu.Lock()
				errorList = append(errorList, err)
				mu.Unlock()
			} else {
				_ = pb.Add(1)
			}
		}()
	}

	select {
	case <-ctx.Done():
		logrus.Error("failed to stop all applications")
	case <-done:
		logrus.Info("all applications stopped")
	}

	if len(errorList) != 0 {
		return multierr.Combine(errorList...)
	}

	return nil
}
