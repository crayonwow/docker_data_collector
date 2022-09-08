package watcher

import (
	"context"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
)

func TestWatcher(t *testing.T) {
	t.Skip(`before run manually start "docker run hello-world"`)
	t.Run("watch", func(t *testing.T) {
		r := require.New(t)
		w, err := NewWatcher()
		r.NoError(err)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		w.WatchContainers(ctx, time.Millisecond*500, func(ctx context.Context, client *docker.Client, container types.Container) error {
			return nil
		})
	})

	t.Run("events", func(t *testing.T) {
		r := require.New(t)
		w, err := NewWatcher()
		r.NoError(err)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()
		w.WatchContainersEvents(ctx)
	})
}
