package service

import (
	"context"
	"github.com/GitEval/GitEval-Backend/model"
	"github.com/GitEval/GitEval-Backend/pkg/github"
)

type AuthService interface {
	Login(ctx context.Context) (url string, err error)
	CallBack(ctx context.Context, code string) (userId int64, err error)
}

type UserServiceProxy interface {
	InitUser(ctx context.Context, u model.User) (err error)
	GetUserById(ctx context.Context, id int64) (model.User, error)
}

type authService struct {
	githubAPI github.GitHubAPI
	u         UserServiceProxy
}

func NewAuthService(u UserServiceProxy, api github.GitHubAPI) AuthService {
	return &authService{
		u: u,
		//因为让其成为中枢，必然要依赖注入到这个authService
		githubAPI: api,
	}
}

func (s *authService) Login(ctx context.Context) (url string, err error) {
	url = s.githubAPI.GetLoginUrl()
	return url, nil
}

func (s *authService) CallBack(ctx context.Context, code string) (userId int64, err error) {
	client, err := s.githubAPI.GetClientByCode(code)
	if err != nil {
		return 0, err
	}

	userInfo, err := s.githubAPI.GetUserInfo(ctx, client, "")
	if err != nil {
		return 0, err
	}

	// 根据用户 ID 查找用户
	user, err := s.u.GetUserById(ctx, userInfo.GetID())
	if err != nil {
		return 0, err
	}

	// 如果用户不存在，创建新用户,如果存在
	if (user == model.User{}) {
		user = model.TransformUser(userInfo)
		// 创建用户
		err = s.u.InitUser(ctx, user)
		if err != nil {
			return 0, err
		}
		//存储用户
		s.githubAPI.SetClient(user.ID, client)
	}

	return user.ID, nil
}
