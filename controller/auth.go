package controller

import (
	"github.com/GitEval/GitEval-Backend/api/request"
	"github.com/GitEval/GitEval-Backend/api/response"
	"github.com/GitEval/GitEval-Backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController interface {
	Login(ctx *gin.Context) error
}

type authController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &authController{authService: authService}
}

// Login 用户登录
// @Summary 用户登录接口
// @Description 用户登录接口
// @Tags Auth
// @Produce json
// @Success 200 {object} response.Success "登录成功"
// @Failure 400 {object} response.Err "请求参数错误"
// @Failure 500 {object} response.Err "内部错误"
// @Router /api/v1/auth/login [get]
func (c *authController) Login(ctx *gin.Context) error {
	url, err := c.authService.Login(ctx)
	if err != nil {
		// 处理错误，比如返回一个错误页面或重定向到错误页面
		ctx.JSON(http.StatusInternalServerError, response.Err{Err: err})
		return nil // 或根据需要返回其他值
	}

	// 重定向到 URL
	ctx.Redirect(http.StatusFound, url) // HTTP 302
	return nil
}

// CallBack 用户在github授权登录之后会被重定向到这里。进行一个请求的发送进行最终验证登录
// CallBack github重定向
// @Summary github重定向
// @Description github重定向,用来初始化这个用户
// @Tags Auth
// @Produce json
// @Success 200 {object} response.Success "初始化成功!"
// @Failure 400 {object} response.Err "请求参数错误"
// @Failure 500 {object} response.Err "内部错误"
// @Router /api/v1/auth/login [get]
func (c *authController) CallBack(ctx *gin.Context) error {

	// 绑定查询参数
	var req request.CallBackReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err{Err: err})
		return nil
	}

	userid, err := c.authService.CallBack(ctx, req.Code)
	if err != nil {
		return err
	}

	ctx.JSON(http.StatusOK, response.Success{
		Data: response.CallBack{
			UserId: userid,
		},
		Msg: "success",
	})
	return nil
}

func (c *authController) Logout(ctx *gin.Context) error {
	//待完成...我觉得能改成jwt就很好了....
	//删除用户的当前jwt并将对应当前jwt列入黑名单
	return nil
}
