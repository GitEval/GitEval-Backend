package controller

import (
	"context"
	"github.com/GitEval/GitEval-Backend/api/request"
	"github.com/GitEval/GitEval-Backend/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthServiceProxy interface {
	Login(ctx context.Context) (url string, err error)
	CallBack(ctx context.Context, code string) (userId int64, err error)
}

type AuthController struct {
	authService AuthServiceProxy
}

func NewAuthController(authService AuthServiceProxy) *AuthController {
	return &AuthController{authService: authService}
}

// Login 用户登录
// @Summary github用户登录授权接口
// @Description github用户登录授权接口,会自动重定向到github的授权接口上
// @Tags Auth
// @Produce json
// @Success 200 {object} response.Success "登录成功"
// @Failure 400 {object} response.Err "请求参数错误"
// @Failure 500 {object} response.Err "内部错误"
// @Router /api/v1/auth/login [get]
func (c *AuthController) Login(ctx *gin.Context) {
	url, err := c.authService.Login(ctx)
	if err != nil {
		// 处理错误，比如返回一个错误页面或重定向到错误页面
		ctx.JSON(http.StatusInternalServerError, response.Err{Err: err})
		return // 或根据需要返回其他值
	}

	// 重定向到 URL
	ctx.Redirect(http.StatusFound, url) // HTTP 302
	return
}

// CallBack 用户在github授权登录之后会被重定向到这里。进行一个请求的发送进行最终验证登录
// CallBack github重定向
// @Summary github重定向
// @Description github重定向,用来初始化这个用户,会返回一个user_id,userid是用在后续的请求上的
// @Tags Auth
// @Param code query string true "github重定向的code"
// @Produce json
// @Success 200 {object} response.Success{Data=response.CallBack} "初始化成功!"
// @Failure 400 {object} response.Err "请求参数错误"
// @Failure 500 {object} response.Err "内部错误"
// @Router /api/v1/auth/login [get]
func (c *AuthController) CallBack(ctx *gin.Context) {

	// 绑定查询参数
	var req request.CallBackReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err{Err: err})
		return
	}

	userid, err := c.authService.CallBack(ctx, req.Code)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, response.Success{
		Data: response.CallBack{
			UserId: userid,
		},
		Msg: "success",
	})
	return
}

func (c *AuthController) Logout(ctx *gin.Context) error {
	//待完成...我觉得能改成jwt就很好了....
	//删除用户的当前jwt并将对应当前jwt列入黑名单
	return nil
}
