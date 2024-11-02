package controller

import (
	"context"
	"fmt"
	"github.com/GitEval/GitEval-Backend/api/request"
	"github.com/GitEval/GitEval-Backend/api/response"
	"github.com/GitEval/GitEval-Backend/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController interface {
	GetUser(ctx *gin.Context)
}
type UserServiceProxy interface {
	GetUserById(ctx context.Context, id int64) (model.User, error)
	GetLeaderboard(ctx context.Context, userId int64) ([]model.Leaderboard, error)
}
type userController struct {
	userService UserServiceProxy
}

func NewUserController(userService UserServiceProxy) UserController {
	return &userController{userService: userService}
}

// GetUser 获取用户
// @Summary 从userid获取用户
// @Tags Auth
// @Accept json request.GetUserInfo
// @Produce json
// @Success 200 {object} response.Success "登录成功"
// @Failure 400 {object} response.Err "请求参数错误"
// @Router /api/v1/user/get/info [get]
func (c *userController) GetUser(ctx *gin.Context) {
	var req request.GetUserInfo
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err{
			Err: err,
		})
		return
	}
	user, err := c.userService.GetUserById(ctx, req.UserID)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Err{
			Err: fmt.Errorf("GetUserById: %w", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Success{
		Data: user,
		Msg:  "success",
	})
}

// GetRanking 获取排行
// @Summary 根据userid获取用户的score的排行榜
// @Tags Auth
// @Accept json request.GetRanking
// @Produce json
// @Success 200 {object} response.Success "登录成功"
// @Failure 400 {object} response.Err "请求参数错误"
// @Router /api/v1/user/get/rank [get]
func (c *userController) GetRanking(ctx *gin.Context) {
	var req request.GetUserInfo
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err{
			Err: err,
		})
		return
	}
	rankings, err := c.userService.GetLeaderboard(ctx, req.UserID)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Err{
			Err: fmt.Errorf("GetLeaderboard: %w", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, response.Success{Data: rankings, Msg: "success"})
}
