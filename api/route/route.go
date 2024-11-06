package route

import (
	"github.com/GitEval/GitEval-Backend/conf"
	"github.com/GitEval/GitEval-Backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"time"
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
	GetNation(ctx *gin.Context)
	GetDomain(ctx *gin.Context)
	SearchUser(ctx *gin.Context)
	GetUserInfo(ctx *gin.Context)
}

func NewRouter(authController AuthControllerProxy, userController UserControllerProxy, m *middleware.Middleware) *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// 添加 CORS 中间件
	r.Use(corsHdl())
	g := r.Group("/api/v1")
	//认证服务
	authGroup := g.Group("/auth")
	authGroup.GET("/login", authController.Login)
	authGroup.GET("/callBack", authController.CallBack)

	//用户服务
	userGroup := g.Group("/user")
	userGroup.GET("/getInfo", m.AuthMiddleware(), userController.GetUser)
	userGroup.GET("/getRank", m.AuthMiddleware(), userController.GetRanking)
	userGroup.GET("/getEvaluation", m.AuthMiddleware(), userController.GetEvaluation)
	userGroup.GET("/getNation", m.AuthMiddleware(), userController.GetNation)
	userGroup.GET("/getDomain", m.AuthMiddleware(), userController.GetDomain)
	userGroup.GET("/search", m.AuthMiddleware(), userController.SearchUser)
	userGroup.GET("/getUserInfo", m.AuthMiddleware(), userController.GetUserInfo)
	return r
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 允许的请求头
		AllowHeaders: []string{"Content-ContentType", "Authorization", "Origin"},
		// 是否允许携带凭证（如 Cookies）
		AllowCredentials: true,
		// 解决跨域问题,这个地方允许所有请求跨域了,之后要改成允许前端的请求,比如localhost
		AllowOriginFunc: func(origin string) bool {
			//检测请求来源是否以localhost开头
			//if strings.HasPrefix(origin, "localhost") {
			//	return true
			//}else {
			//	return false
			//}
			return true
		},

		// 预检请求的缓存时间
		MaxAge: 12 * time.Hour,
	})
}
