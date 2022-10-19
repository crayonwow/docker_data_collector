//go:generate mockgen -source=sender.go -destination=sender_mocks.go -package=sender

package sender

type (
	Sender interface {
		Send(message string) error
	}
)
