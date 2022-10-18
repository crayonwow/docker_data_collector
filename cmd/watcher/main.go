package main

import (
	"docker_data_collector/internal"
	"docker_data_collector/pkg/di"
	"docker_data_collector/pkg/graceful"
)

func main() {
	err := di.Start(internal.Run, internal.Module().Append(di.NewModule(di.NewDependency(graceful.NewContext))))
	if err != nil {
		panic(err)
	}
}
