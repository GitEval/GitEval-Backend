// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/GitEval/GitEval-Backend/api/route"
	"github.com/GitEval/GitEval-Backend/conf"
	"github.com/GitEval/GitEval-Backend/controller"
	"github.com/GitEval/GitEval-Backend/middleware"
	"github.com/GitEval/GitEval-Backend/model"
	"github.com/GitEval/GitEval-Backend/pkg/github"
	"github.com/GitEval/GitEval-Backend/pkg/github/expireMap"
	"github.com/GitEval/GitEval-Backend/pkg/llm"
	"github.com/GitEval/GitEval-Backend/service"
)

// Injectors from wire.go:

func WireApp(confPath string) (route.App, func()) {
	vipperSetting := conf.NewVipperSetting(confPath)
	dataConfig := conf.NewDataConfig(vipperSetting)
	db := model.NewDB(dataConfig)
	data := model.NewData(db)
	gormUserDAO := model.NewGormUserDAO(data)
	gormContactDAO := model.NewGormContactDAO(data)
	gormDomainDAO := model.NewGormDomainDAO(data)
	gitHubConfig := conf.NewGitHubConfig(vipperSetting)
	expireMapExpireMap, cleanup := expireMap.NewExpireMap()
	gitHubAPI := github.NewGitHubAPI(gitHubConfig, expireMapExpireMap)
	llmConfig := conf.NewLLMConfig(vipperSetting)
	llmClient := llm.NewLLMClient(llmConfig)
	userService := service.NewUserService(gormUserDAO, gormContactDAO, gormDomainDAO, data, gitHubAPI, llmClient)
	authService := service.NewAuthService(userService, gitHubAPI, llmClient)
	jwtConfig := conf.NewJWTConfig(vipperSetting)
	jwtClient := middleware.NewJWTClient(jwtConfig)
	authController := controller.NewAuthController(authService, jwtClient)
	userController := controller.NewUserController(userService)
	middlewareMiddleware := middleware.NewMiddleware(jwtClient)
	engine := route.NewRouter(authController, userController, middlewareMiddleware)
	appConf := conf.NewAppConf(vipperSetting)
	app := route.NewApp(engine, appConf)
	return app, func() {
		cleanup()
	}
}
