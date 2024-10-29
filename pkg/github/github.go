package github

import (
	"context"
	"github.com/google/go-github/v50/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"sync"
)

type GitHubAPI interface {
	GetLoginUrl() string
	SetClient(userID int64, client *github.Client)
	GetClientFromMap(userID int64) (*github.Client, bool)
	GetClientByCode(code string) (*github.Client, error)
	GetUserInfo(ctx context.Context, client *github.Client) (*github.User, error)
}

// gitHubAPI 结构体
type gitHubAPI struct {
	clients sync.Map      // 使用 sync.Map 实现并发安全
	cfg     *gitHubConfig // 引用的地址完全相同节约了内存空间
}

// 使用统一的cfg管理方案
type gitHubConfig struct {
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
}

func NewGitHubAPI() GitHubAPI {
	var cfg gitHubConfig
	err := viper.UnmarshalKey("github", &cfg)
	if err != nil {
		return nil
	}
	//每次尝试去获取一个新的githubAPI的时候就直接引用这个配置文件的地址
	return &gitHubAPI{cfg: &cfg}
}

// SetClient 设置用户的 GitHub 客户端
func (g *gitHubAPI) SetClient(userID int64, client *github.Client) {
	g.clients.Store(userID, client) // 使用 Store 方法
}

// GetClient 获取用户的 GitHub 客户端
func (g *gitHubAPI) GetClientFromMap(userID int64) (*github.Client, bool) {
	client, exists := g.clients.Load(userID) // 使用 Load 方法
	if exists {
		return client.(*github.Client), true // 类型断言
	}
	return nil, false
}

func (g *gitHubAPI) GetLoginUrl() string {
	redirectURL := "https://github.com/login/oauth/authorize?client_id=" + g.cfg.ClientID + "&scope=user"
	return redirectURL
}

func (g *gitHubAPI) GetClientByCode(code string) (*github.Client, error) {
	// 获取 access token
	token, err := g.getAccessToken(code)
	if err != nil {
		return nil, err
	}

	// 使用 access token 创建 GitHub 客户端
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return client, nil
}

func (g *gitHubAPI) GetUserInfo(ctx context.Context, client *github.Client) (*github.User, error) {
	userInfo, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (g *gitHubAPI) GetReposDetailList(repoUrls []string) (string, error) {
	return "", nil
}

func (g *gitHubAPI) getAccessToken(code string) (string, error) {
	// 创建 OAuth2 端点
	oauth2Endpoint := oauth2.Endpoint{
		TokenURL: "https://github.com/login/oauth/access_token",
	}

	// 创建 OAuth2 客户端
	ctx := context.Background()
	cf := oauth2.Config{
		ClientID:     g.cfg.ClientID,
		ClientSecret: g.cfg.ClientSecret,
		Endpoint:     oauth2Endpoint,
	}

	// 获取访问令牌
	token, err := cf.Exchange(ctx, code)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
