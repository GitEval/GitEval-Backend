package route

import (
	"github.com/GitEval/GitEval-Backend/conf"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewApp,
	NewRouter,
)

type App struct {
	r *gin.Engine
	c *conf.AppConf
}

func NewApp(r *gin.Engine, c *conf.AppConf) App {
	return App{
		r: r,
		c: c,
	}
}

//启动
func (a *App) Run() {
	a.r.Run(a.c.Addr)
}

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//后续的接口应该用group来管理
	//例如:
	/*
			UserGroup := r.Group("/api/user")
		{
			//conf也可以注入到user中
			user := NewUser()
			UserGroup.POST("/login", user.Login)
			UserGroup.GET("/getinfo", middleware.JWTAuthMiddleware(), user.GetUserInfo)
		}
	*/
	return r
}
