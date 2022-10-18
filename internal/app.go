package internal

import (
	"context"
	"fmt"

	"docker_data_collector/pkg/di"
)

func Run(ctx context.Context, pool di.ApplicationPool) error {
	if err := pool.Run(ctx); err != nil {
		return fmt.Errorf("pool run: %w", err)
	}
	<-ctx.Done()
	if err := pool.Stop(); err != nil {
		return fmt.Errorf("pool run: %w", err)
	}

	return nil
}
