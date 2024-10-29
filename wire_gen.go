// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/GitEval/GitEval-Backend/api/route"
	"github.com/GitEval/GitEval-Backend/conf"
)

// Injectors from wire.go:

func WireApp(confPath string) route.App {
	engine := route.NewRouter()
	vipperSetting := conf.NewVipperSetting(confPath)
	appConf := conf.NewAppConf(vipperSetting)
	app := route.NewApp(engine, appConf)
	return app
}
