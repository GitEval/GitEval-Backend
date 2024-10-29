package controller

import (
	"github.com/GitEval/GitEval-Backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController interface {
	Login(ctx *gin.Context) error
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{userService: userService}
}

func (c *userController) GetUser(ctx *gin.Context) error {
	url, err := c.UserService.Login(ctx)
	if err != nil {
		// 处理错误，比如返回一个错误页面或重定向到错误页面
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil // 或根据需要返回其他值
	}

	// 重定向到 URL
	ctx.Redirect(http.StatusFound, url) // HTTP 302
	return nil                          // 重定向后通常不需要返回
}
