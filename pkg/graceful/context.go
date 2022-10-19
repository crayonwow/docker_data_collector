package graceful

import (
	"context"
	"os/signal"
	"syscall"
)

func NewContext() context.Context {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	return ctx
}
