package internal

import (
	"docker_data_collector/internal/config"
	"docker_data_collector/internal/sender"
	"docker_data_collector/internal/watcher"
	"docker_data_collector/pkg/di"

	"go.uber.org/dig"
)

func Module() di.Module {
	return di.NewModule(di.NewDependency(appAdapter)).
		Append(config.Module()).
		Append(sender.Module()).
		Append(watcher.Module())
}

type (
	appIn struct {
		dig.In

		Applications []di.Application `group:"applications"`
	}

	appOut struct {
		dig.Out

		Application di.ApplicationPool
	}
)

func appAdapter(in appIn) appOut {
	return appOut{
		Application: in.Applications,
	}
}
