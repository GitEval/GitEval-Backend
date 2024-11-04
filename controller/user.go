package controller

import (
	"context"
	"fmt"
	"github.com/GitEval/GitEval-Backend/api/request"
	"github.com/GitEval/GitEval-Backend/api/response"
	"github.com/GitEval/GitEval-Backend/errs"
	"github.com/GitEval/GitEval-Backend/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserServiceProxy interface {
	GetUserById(ctx context.Context, id int64) (model.User, error)
	GetLeaderboard(ctx context.Context, userId int64) ([]model.Leaderboard, error)
	GetDomains(ctx context.Context, userId int64) []string
	GetEvaluation(ctx context.Context, userId int64) (string, error)
}
type UserController struct {
	userService UserServiceProxy
}

func NewUserController(userService UserServiceProxy) *UserController {
	return &UserController{userService: userService}
}

// GetUser 获取用户
// @Summary 从userid获取用户
// @Tags User
// @Param user_id query string true "用户ID,暂时没写jwt和cookie之类的,所以直接传user_id"
// @Produce json
// @Success 200 {object} response.Success{Data=response.User} "登录成功"
// @Failure 400 {object} response.Err "请求参数错误"
// @Router /api/v1/user/get/info [get]
func (c *UserController) GetUser(ctx *gin.Context) {
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

	domain := c.userService.GetDomains(ctx, req.UserID)
	ctx.JSON(http.StatusOK, response.Success{
		Data: response.User{
			U:      user,
			Domain: domain,
		},
		Msg: "success",
	})
}

// GetRanking 获取排行
// @Summary 根据userid获取用户的score的排行榜
// @Tags User
// @Param user_id query string true "用户ID,暂时没写jwt和cookie之类的,所以直接传user_id"
// @Produce json
// @Success 200 {object} response.Success{Data=response.Ranking} "登录成功"
// @Failure 400 {object} response.Err "请求参数错误"
// @Router /api/v1/user/get/rank [get]
func (c *UserController) GetRanking(ctx *gin.Context) {
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
	ctx.JSON(http.StatusOK, response.Success{Data: response.Ranking{Leaderboard: rankings}, Msg: "success"})
}

// GetEvaluation 获取用户评价
// @Summary 根据userid获取用户评价
// @Tags User
// @Param user_id query string true "用户ID,暂时没写jwt和cookie之类的,所以直接传user_id"
// @Produce json
// @Success 200 {object} response.Success{Data=response.Evaluation} "登录成功"
// @Failure 400 {object} response.Err "请求参数错误"
// @Router /api/v1/user/get/evaluation [get]
func (c *UserController) GetEvaluation(ctx *gin.Context) {
	var req request.GetEvaluation
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err{
			Err: err,
		})
		return
	}

	evaluation, err := c.userService.GetEvaluation(ctx, req.UserID)
	switch err {
	case errs.LoginFailErr:
		//返回一个重定向的状态码,让前端做重定向,因为后端得不到实际的ip,我暂时只对这里进行了处理,看看cc有没有更好的想法
		ctx.JSON(http.StatusFound, response.Err{
			Err: fmt.Errorf("GetEvaluation: %w", err),
		})
		return
	case nil:
		ctx.JSON(http.StatusOK, response.Success{Data: response.Evaluation{Evaluation: evaluation}, Msg: "success"})
	default:
		ctx.JSON(http.StatusOK, response.Err{
			Err: fmt.Errorf("GetEvaluation: %w", err),
		})
		return
	}
}
