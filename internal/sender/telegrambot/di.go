package telegrambot

import (
	"docker_data_collector/pkg/di"

	"go.uber.org/dig"
)

func Module() di.Module {
	return di.NewModule(
		di.NewDependency(newBot, dig.Name("telegram-bot")),
	)
}
