package main

import (
	"docker_data_collector/internal"
	"docker_data_collector/pkg/config"
	"docker_data_collector/pkg/di"
	"docker_data_collector/pkg/graceful"
	"docker_data_collector/pkg/logger"

	"go.uber.org/dig"
)

func main() {
	err := di.Start(internal.Run, internal.Module().
		Append(config.Module()).
		Append(logger.Module()).
		Append(
			di.NewModule(
				di.NewDependency(graceful.NewContext),
				di.NewDependency(func() string { return "hello" }, dig.Name("config_path")),
			),
		),
	)
	if err != nil {
		panic(err)
	}
}
