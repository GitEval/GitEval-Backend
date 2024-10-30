package controller

import (
	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetUser(ctx *gin.Context) error
}
type UserServiceProxy interface {
	//待填,别急
}
type userController struct {
	userService UserServiceProxy
}

func NewUserController(userService UserServiceProxy) UserController {
	return &userController{userService: userService}
}

func (c *userController) GetUser(ctx *gin.Context) error {
	return nil
}
