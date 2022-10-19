package telegrambot

import (
	"docker_data_collector/pkg/di"
)

func Module() di.Module {
	return di.NewModule(
		di.NewDependency(newBot),
	)
}
