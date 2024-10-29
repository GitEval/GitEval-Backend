package service

import (
	"context"
	"github.com/GitEval/GitEval-Backend/pkg/github"
)

type AuthService interface {
	Login(ctx context.Context) (url string, err error)
}

type authService struct {
}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Login(ctx context.Context) (url string, err error) {
	githubAPI := github.NewGitHubAPI()
	url = githubAPI.GetLoginUrl()
	return url, nil
}

func (s *authService) CallBack(ctx context.Context) error {
	return nil
}
