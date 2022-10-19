package internal

import (
	"context"
	"fmt"
	"time"

	"docker_data_collector/pkg/di"
)

func Run(ctx context.Context, pool di.ApplicationPool) error {
	if err := pool.Run(ctx); err != nil {
		return fmt.Errorf("pool run: %w", err)
	}

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := pool.Stop(ctx); err != nil {
		return fmt.Errorf("pool run: %w", err)
	}

	return nil
}
