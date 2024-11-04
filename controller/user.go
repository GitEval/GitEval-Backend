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
	GetNationByUserId(ctx context.Context, userId int64) (string, error)
	GetDomainByUserId(ctx context.Context, userId int64) ([]string, error)
	SearchUser(ctx context.Context, nation, domain string, page int, pageSize int) ([]model.User, error)
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
// @Router /api/v1/user/getInfo [get]
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

// GetRanking 获取用户排行(和自己的好友之间的
// @Summary 根据userid获取用户的score的排行榜
// @Tags User
// @Param user_id query string true "用户ID,暂时没写jwt和cookie之类的,所以直接传user_id"
// @Produce json
// @Success 200 {object} response.Success{Data=response.Ranking} "登录成功"
// @Failure 400 {object} response.Err "请求参数错误"
// @Router /api/v1/user/getRank [get]
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
// @Success 200 {object} response.Success{Data=response.EvaluationResp} "登录成功"
// @Failure 400 {object} response.Err "请求参数错误"
// @Router /api/v1/user/getEvaluation [get]
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
		ctx.JSON(http.StatusOK, response.Success{Data: response.EvaluationResp{Evaluation: evaluation}, Msg: "success"})
	default:
		ctx.JSON(http.StatusOK, response.Err{
			Err: fmt.Errorf("GetEvaluation: %w", err),
		})
		return
	}
}

// GetNation 获取用户所在国家
// @Summary 根据用户 ID 获取用户所在国家
// @Tags User
// @Param user_id query string true "用户ID, 暂时没写 JWT 和 cookie 之类的, 所以直接传 user_id"
// @Produce json
// @Success 200 {object} response.Success{Data=response.NationResp} "国家获取成功"
// @Failure 400 {object} response.Err "请求参数错误"
// @Failure 404 {object} response.Err "用户未找到"
// @Router /api/v1/user/getNation [get]
func (c *UserController) GetNation(ctx *gin.Context) {
	var req request.GetNation
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err{
			Err: err,
		})
		return
	}

	nation, err := c.userService.GetNationByUserId(ctx, req.UserID)
	switch err {
	case errs.LoginFailErr:
		//返回一个重定向的状态码,让前端做重定向,因为后端得不到实际的ip,我暂时只对这里进行了处理,看看cc有没有更好的想法
		ctx.JSON(http.StatusFound, response.Err{
			Err: fmt.Errorf("GetDomain: %w", err),
		})
		return
	case nil:
		ctx.JSON(http.StatusOK, response.Success{Data: response.NationResp{Nation: nation}, Msg: "success"})
	default:
		ctx.JSON(http.StatusOK, response.Err{
			Err: fmt.Errorf("GetDomain: %w", err),
		})
		return
	}
}

// GetDomain 获取用户的领域
// @Summary 根据用户 ID 获取用户的领域
// @Tags User
// @Param user_id query string true "用户ID, 暂时没写 JWT 和 cookie 之类的, 所以直接传 user_id"
// @Produce json
// @Success 200 {object} response.Success{Data=response.DomainResp} "领域获取成功"
// @Failure 400 {object} response.Err "请求参数错误"
// @Failure 404 {object} response.Err "用户未找到"
// @Router /api/v1/user/getDomain [get]
func (c *UserController) GetDomain(ctx *gin.Context) {
	var req request.GetDomain
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err{
			Err: err,
		})
		return
	}

	domain, err := c.userService.GetDomainByUserId(ctx, req.UserID)
	switch err {
	case errs.LoginFailErr:
		//返回一个重定向的状态码,让前端做重定向,因为后端得不到实际的ip,我暂时只对这里进行了处理,看看cc有没有更好的想法
		ctx.JSON(http.StatusFound, response.Err{
			Err: fmt.Errorf("GetDomain: %w", err),
		})
		return
	case nil:
		ctx.JSON(http.StatusOK, response.Success{Data: response.DomainResp{Domain: domain}, Msg: "success"})
	default:
		ctx.JSON(http.StatusOK, response.Err{
			Err: fmt.Errorf("GetDomain: %w", err),
		})
		return
	}
}

// SearchUser 根据国家和领域搜索用户
// @Summary 根据国家和领域搜索用户
// @Tags User
// @Param nation query string false "国家，选择性参数"
// @Param domain query string false "领域，选择性参数"
// @Param page query int false "分页参数表示这是第几页"
// @Param pageSize query int false "每页返回的用户数量，建议一次返回10个"
// @Produce json
// @Success 200 {object} response.Success{Data=response.SearchResp} "用户搜索成功"
// @Failure 400 {object} response.Err "请求参数错误"
// @Router /api/v1/user/search [get]
func (c *UserController) SearchUser(ctx *gin.Context) {
	var req request.SearchUser
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err{
			Err: err,
		})
		return
	}

	users, err := c.userService.SearchUser(ctx, req.Nation, req.Domain, req.Page, req.PageSize)
	switch err {
	case errs.LoginFailErr:
		//返回一个重定向的状态码,让前端做重定向,因为后端得不到实际的ip,我暂时只对这里进行了处理,看看cc有没有更好的想法
		ctx.JSON(http.StatusFound, response.Err{
			Err: fmt.Errorf("SearchUser: %w", err),
		})
		return
	case nil:
		ctx.JSON(http.StatusOK, response.Success{Data: response.SearchResp{Users: users}, Msg: "success"})
	default:
		ctx.JSON(http.StatusOK, response.Err{
			Err: fmt.Errorf("SearchUser: %w", err),
		})
		return
	}
}
