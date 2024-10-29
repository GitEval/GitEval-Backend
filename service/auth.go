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

type authService struct {
	userDAO   model.UserDAO
	githubAPI github.GitHubAPI
}

func NewAuthService(userDAO model.UserDAO, api github.GitHubAPI) AuthService {
	return &authService{
		userDAO: userDAO,
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

	userInfo, err := s.githubAPI.GetUserInfo(ctx, client)
	if err != nil {
		return 0, err
	}

	// 根据用户 ID 查找用户
	user, err := s.userDAO.GetUserByID(userInfo.GetID())
	if err != nil {
		return 0, err
	}

	// 如果用户不存在，创建新用户,如果存在
	if user == nil {
		user = &model.User{
			ID:                userInfo.GetID(),
			AvatarURL:         userInfo.GetAvatarURL(),
			Name:              userInfo.Name,
			Company:           userInfo.Company,
			Blog:              userInfo.Blog,
			Location:          userInfo.Location,
			Email:             userInfo.GetEmail(),
			Bio:               userInfo.Bio,
			PublicRepos:       userInfo.GetPublicRepos(),
			Followers:         userInfo.GetFollowers(),
			Following:         userInfo.GetFollowing(),
			TotalPrivateRepos: userInfo.GetTotalPrivateRepos(),
			Collaborators:     userInfo.GetCollaborators(),
		}
		// 创建用户
		err = s.userDAO.CreateUser(user)
		if err != nil {
			return 0, err
		}
		//存储用户到
		s.githubAPI.SetClient(user.ID, client)
	}

	return user.ID, nil
}
