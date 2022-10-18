package sender

import (
	"docker_data_collector/internal/sender/telegrambot"
	"docker_data_collector/pkg/di"
)

func Module() di.Module {
	return di.NewModule().
		Append(telegrambot.Module())
}
