package main

import (
	"docker_data_collector/internal"
	"docker_data_collector/pkg/config"
	"docker_data_collector/pkg/di"
	"docker_data_collector/pkg/graceful"
	"docker_data_collector/pkg/logger"
)

func main() {
	di.
		NewContainer().
		Provide(
			internal.Module(),
			config.Module(),
			di.NewModule(
				di.NewDependency(graceful.NewContext),
			)).
		Invoke(logger.NewLogger).
		Main(internal.Run)
}
