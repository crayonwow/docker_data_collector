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

	Container struct {
		c *dig.Container
	}
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

func NewContainer() *Container {
	return &Container{
		c: dig.New(),
	}
}

func (c *Container) Provide(mods ...Module) *Container {
	for _, m := range mods {
		for _, d := range m {
			err := c.c.Provide(d.constructor, d.ops...)
			if err != nil {
				panic(err)
			}
		}
	}

	return c
}

func (c *Container) Invoke(fs ...interface{}) *Container {
	for _, f := range fs {
		err := c.c.Invoke(f)
		if err != nil {
			panic(err)
		}
	}
	return c
}

func (c *Container) Main(main Main) {
	c.Invoke(main)
}
