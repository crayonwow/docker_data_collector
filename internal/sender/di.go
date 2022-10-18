package sender

import (
	"docker_data_collector/internal/sender/telegrambot"
	"docker_data_collector/pkg/di"

	"go.uber.org/dig"
)

func Module() di.Module {
	return di.NewModule(di.NewDependency(senderAdapter)).
		Append(telegrambot.Module())
}

type (
	senderAdapterIn struct {
		dig.In

		TGBot *telegrambot.Bot
	}

	senderAdapterOut struct {
		dig.Out

		TGBot Sender
	}
)

func senderAdapter(in senderAdapterIn) senderAdapterOut {
	return senderAdapterOut{
		TGBot: in.TGBot,
	}
}
