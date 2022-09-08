package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

var (
	once sync.Once
	ctx  context.Context
	// cancel context.CancelFunc
)

func Context() context.Context {
	once.Do(func() {
		ctx, _ = signal.NotifyContext(context.Background(), os.Interrupt)
	})

	return ctx
}
