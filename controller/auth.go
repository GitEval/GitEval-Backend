package controller

import (
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
// @Accept json
// @Produce json
// @Param email query string true "邮箱"
// @Param password query string true "密码"
// @Success 200 {object} LoginResponse "登录成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/auth/login [post]
func (c *authController) Login(ctx *gin.Context) error {
	url, err := c.authService.Login(ctx)
	if err != nil {
		// 处理错误，比如返回一个错误页面或重定向到错误页面
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil // 或根据需要返回其他值
	}

	// 重定向到 URL
	ctx.Redirect(http.StatusFound, url) // HTTP 302
	return nil                          // 重定向后通常不需要返回
}
