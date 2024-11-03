package service

import (
	"context"
	"github.com/GitEval/GitEval-Backend/model"
	"github.com/GitEval/GitEval-Backend/pkg/github"
	"github.com/GitEval/GitEval-Backend/pkg/llm"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx context.Context) (url string, err error)
	CallBack(ctx context.Context, code string) (userId int64, err error)
}

type UserServiceProxy interface {
	InitUser(ctx context.Context, u model.User) (err error)
	GetUserById(ctx context.Context, id int64) (model.User, error)
	SaveUser(ctx context.Context, user model.User) error
}

type authService struct {
	githubAPI github.GitHubAPI
	u         UserServiceProxy
	llmClient llm.LLMClient
}

func NewAuthService(u UserServiceProxy, api github.GitHubAPI, llmClient llm.LLMClient) AuthService {
	return &authService{
		u: u,
		//因为让其成为中枢，必然要依赖注入到这个authService
		githubAPI: api,
		llmClient: llmClient,
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
	switch err {
	case gorm.ErrRecordNotFound:
		user = model.TransformUser(userInfo)
		// 创建用户
		err = s.u.InitUser(ctx, user)
		if err != nil {
			return 0, err
		}

		//存储用户
		s.githubAPI.SetClient(user.ID, client)
		go func() {
			ctx = context.Background()

			var followerAreas []string
			followers := s.githubAPI.GetFollowers(ctx, user.ID)
			for _, follower := range followers {
				followerAreas = append(followerAreas, *follower.Location)
			}

			var followingAreas []string
			followings := s.githubAPI.GetFollowers(ctx, user.ID)
			for _, following := range followings {
				followingAreas = append(followingAreas, *following.Location)
			}

			area, err := s.llmClient.GetArea(ctx, llm.GetAreaRequest{
				Bio:            *user.Bio,
				Company:        *user.Company,
				Location:       *user.Location,
				FollowerAreas:  followerAreas,
				FollowingAreas: followingAreas,
			})
			if err != nil {
				return
			}

			if area.Confidence >= 0.5 {
				user.Nationality = &area.Area
			} else {
				*user.Nationality = "N/A"
			}

			//先尝试保存一次
			err = s.u.SaveUser(ctx, user)
			if err != nil {
				return
			}

			//获取仓库
			details, err := s.githubAPI.GetAllRepos(ctx, userInfo, client)
			if err != nil {
				return
			}
			var repos []llm.Repo
			for _, detail := range details {
				repos = append(repos, llm.Repo{
					Name:     detail.Name,
					Readme:   detail.Readme,
					Language: detail.Language,
					Commit:   detail.Commit,
					Star:     detail.Star,
					Fork:     detail.Fork,
				})
			}

			domain, err := s.llmClient.GetDomain(ctx, llm.GetDomainRequest{
				Repos: repos,
				Bio:   *user.Bio,
			})
			if err != nil {
				return
			}

			user.Domain = domain.Domain
			//尝试保存一次
			err = s.u.SaveUser(ctx, user)
			if err != nil {
				return
			}

		}()
	case nil:
		//尝试获取用户的技术领域
		if user.Nationality == nil {
			go func() {
				ctx = context.Background()

				var followerAreas []string
				followers := s.githubAPI.GetFollowers(ctx, user.ID)
				for _, follower := range followers {
					followerAreas = append(followerAreas, *follower.Location)
				}

				var followingAreas []string
				followings := s.githubAPI.GetFollowers(ctx, user.ID)
				for _, following := range followings {
					followingAreas = append(followingAreas, *following.Location)
				}

				area, err := s.llmClient.GetArea(ctx, llm.GetAreaRequest{
					Bio:            *user.Bio,
					Company:        *user.Company,
					Location:       *user.Location,
					FollowerAreas:  followerAreas,
					FollowingAreas: followingAreas,
				})
				if err != nil {
					return
				}

				if area.Confidence >= 0.5 {
					user.Nationality = &area.Area
				} else {
					*user.Nationality = "N/A"
				}

				//先尝试保存一次
				err = s.u.SaveUser(ctx, user)
				if err != nil {
					return
				}
			}()
		}
		//尝试获取用户的技术领域
		if user.Domain == nil {
			go func() {
				ctx = context.Background()
				//获取仓库
				details, err := s.githubAPI.GetAllRepos(ctx, userInfo, client)
				if err != nil {
					return
				}
				var repos []llm.Repo
				for _, detail := range details {
					repos = append(repos, llm.Repo{
						Name:     detail.Name,
						Readme:   detail.Readme,
						Language: detail.Language,
						Commit:   detail.Commit,
						Star:     detail.Star,
						Fork:     detail.Fork,
					})
				}

				domain, err := s.llmClient.GetDomain(ctx, llm.GetDomainRequest{
					Repos: repos,
					Bio:   *user.Bio,
				})
				if err != nil {
					return
				}

				user.Domain = domain.Domain
				//尝试保存一次
				err = s.u.SaveUser(ctx, user)
				if err != nil {
					return
				}
			}()
		}

	default:
		return 0, err

	}

	return user.ID, nil
}
