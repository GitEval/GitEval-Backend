package service

import (
	"context"
	"github.com/GitEval/GitEval-Backend/model"
	"github.com/GitEval/GitEval-Backend/pkg/llm"
	"github.com/google/go-github/v50/github"
)

type GitHubAPIProxy interface {
	GetLoginUrl() string
	SetClient(userID int64, client *github.Client)
	GetClientFromMap(userID int64) (*github.Client, bool)
	GetClientByCode(code string) (*github.Client, error)
	GetUserInfo(ctx context.Context, client *github.Client, username string) (*github.User, error)
}

// LLMClientProxy 接口定义
type LLMClientProxy interface {
	GetDomain(ctx context.Context, req llm.GetDomainRequest) (llm.GetDomainResponse, error)
	GetEvaluation(ctx context.Context, req llm.GetEvaluationRequest) (llm.GetEvaluationResponse, error)
}

type UserServiceProxy interface {
	InitUser(ctx context.Context, u model.User) (err error)
	GetUserById(ctx context.Context, id int64) (model.User, error)
	CreateUser(ctx context.Context, u model.User) error
}

type AuthService struct {
	githubAPI GitHubAPIProxy
	u         UserServiceProxy
	llmClient LLMClientProxy
}

func NewAuthService(u UserServiceProxy, api GitHubAPIProxy, llmClient LLMClientProxy) *AuthService {
	return &AuthService{
		u: u,
		//因为让其成为中枢，必然要依赖注入到这个authService
		githubAPI: api,
		llmClient: llmClient,
	}
}

func (s *AuthService) Login(ctx context.Context) (url string, err error) {
	url = s.githubAPI.GetLoginUrl()
	return url, nil
}

func (s *AuthService) CallBack(ctx context.Context, code string) (userId int64, err error) {
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
	// 如果用户不存在，创建新用户,如果存在
	if (user == model.User{}) {
		user = model.User{
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
		//首次创建用户
		err = s.u.CreateUser(ctx, user)
		if err != nil {
			return 0, err
		}
		//这里做异步主要是为了保证用户体验,否则等待时间过长了
		go func() {
			// 初始化用户关系网
			err = s.u.InitUser(context.Background(), user)
			if err != nil {
				return
			}
		}()
		//存储用户
		s.githubAPI.SetClient(user.ID, client)
	}

	return user.ID, nil
}
