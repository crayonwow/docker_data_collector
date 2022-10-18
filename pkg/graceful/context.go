package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func NewContext() context.Context {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	return ctx
}
