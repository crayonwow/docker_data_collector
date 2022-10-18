package di

import (
	"context"

	"go.uber.org/dig"
)

type (
	dependency struct {
		constructor interface{}
		ops         []dig.ProvideOption
	}

	Module []dependency

	Main func(ctx context.Context, pool ApplicationPool) error
)

func (m Module) Append(n Module) Module {
	return append(m, n...)
}

func NewDependency(c interface{}, ops ...dig.ProvideOption) dependency {
	return dependency{
		constructor: c,
		ops:         ops,
	}
}

func NewModule(deps ...dependency) Module {
	return deps
}

func Start(main Main, mods ...Module) error {
	c := dig.New()
	for _, m := range mods {
		for _, d := range m {
			err := c.Provide(d.constructor, d.ops...)
			if err != nil {
				return err
			}
		}
	}

	return c.Invoke(main)
}
