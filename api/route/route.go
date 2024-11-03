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

// 启动
func (a *App) Run() {
	a.r.Run(a.c.Addr)
}

type AuthControllerProxy interface {
	Login(ctx *gin.Context)
	CallBack(ctx *gin.Context)
}
type UserControllerProxy interface {
	GetUser(ctx *gin.Context)
	GetRanking(ctx *gin.Context)
	GetEvaluation(ctx *gin.Context)
}

func NewRouter(authController AuthControllerProxy, userController UserControllerProxy) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	g := r.Group("/api/v1")
	//认证服务
	authGroup := g.Group("/auth")
	authGroup.GET("/login", authController.Login)
	authGroup.GET("/callBack", authController.CallBack)

	//用户服务
	userGroup := g.Group("/auth")
	userGroup.GET("/get/info", userController.GetUser)
	userGroup.GET("/get/rank", userController.GetRanking)
	userGroup.GET("/get/evaluation", userController.GetEvaluation)

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
