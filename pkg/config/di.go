package config

import (
	"docker_data_collector/pkg/di"

	"go.uber.org/dig"
)

func Module() di.Module {
	return di.NewModule(
		di.NewDependency(NewConfig),
		di.NewDependency(LoggerConfig),
	)
}

type (
	configIn struct {
		dig.In

		ConfigPath string `name:"config_path"`
	}
)
