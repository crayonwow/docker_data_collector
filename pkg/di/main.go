package di

import (
	"context"
	"fmt"
	"time"
)

func main(ctx context.Context, pool ApplicationPool) error {
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
