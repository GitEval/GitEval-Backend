//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/GitEval/GitEval-Backend/api/route"
	"github.com/GitEval/GitEval-Backend/conf"
	"github.com/google/wire"
)

func WireApp(confPath string) route.App {
	panic(wire.Build(
		conf.ProviderSet,
		route.ProviderSet,
	))
}
