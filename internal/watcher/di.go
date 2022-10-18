package watcher

import (
	"docker_data_collector/pkg/di"

	"go.uber.org/dig"
)

func Module() di.Module {
	return di.NewModule(
		di.NewDependency(newWatcher),
		di.NewDependency(applicationAdapter),
	)
}

type (
	applicationAdapterIn struct {
		dig.In

		Watcher *watcher
	}

	applicationAdapterOut struct {
		dig.Out

		Watcher di.Application `group:"applications"`
	}
)

func applicationAdapter(in applicationAdapterIn) applicationAdapterOut {
	return applicationAdapterOut{
		Watcher: in.Watcher,
	}
}
