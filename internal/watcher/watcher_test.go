package watcher

import (
	"context"
	"testing"
	"time"

	senderpkg "docker_data_collector/internal/sender"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestWatcher(t *testing.T) {
	t.Skip(`before run manually start "docker run hello-world"`)

	t.Run("watch", func(t *testing.T) {
		//r := require.New(t)
		//ctl := gomock.NewController(t)
		//
		//senderMock := senderpkg.NewMockSender(ctl)
		//senderMock.EXPECT().Send("hehel")
		//
		//w, err := NewWatcher(senderMock)
		//r.NoError(err)
		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		//defer cancel()
		//w.watchContainers(ctx,
		//	time.Millisecond*500,
		//	func(ctx context.Context, client *docker.Client, container types.Container,
		//	) error {
		//		return nil
		//	})
	})

	t.Run("events", func(t *testing.T) {
		r := require.New(t)
		ctl := gomock.NewController(t)

		senderMock := senderpkg.NewMockSender(ctl)
		senderMock.EXPECT().Send("hehel")

		w, err := NewWatcher(senderMock)
		r.NoError(err)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()
		w.watchContainersEvents(ctx)
	})
}
