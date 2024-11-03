//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/GitEval/GitEval-Backend/api/route"
	"github.com/GitEval/GitEval-Backend/conf"
	"github.com/GitEval/GitEval-Backend/controller"
	"github.com/GitEval/GitEval-Backend/model"
	"github.com/GitEval/GitEval-Backend/pkg"
	"github.com/GitEval/GitEval-Backend/service"
	"github.com/google/wire"
)

func WireApp(confPath string) route.App {
	panic(wire.Build(
		conf.ProviderSet,
		controller.ProviderSet,
		service.ProviderSet,
		model.ProviderSet,
		pkg.ProviderSet,
		route.ProviderSet,
	))
}
