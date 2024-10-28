package github

import (
	"context"
	"github.com/google/go-github/v50/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type gitHubAPI struct {
	client *github.Client //需要额外进行配置
	cfg    *gitHubConfig  //引用的地址完全相同节约了内存空间
}

type GitHubAPI interface {
	GetLoginUrl() string
	SetClient(code string) error
	GetUserInfo(ctx context.Context) (*github.User, error)
}

// 使用统一的cfg管理方案
type gitHubConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

var cfg gitHubConfig

func InitGitHubConfig() error {
	err := viper.UnmarshalKey("github", &cfg)
	if err != nil {
		return nil
	}
	return nil
}

func NewGitHubAPI() GitHubAPI {
	//每次尝试去获取一个新的githubAPI的时候就直接引用这个配置文件的地址
	return &gitHubAPI{cfg: &cfg}
}

func (g *gitHubAPI) GetLoginUrl() string {
	redirectURL := "https://github.com/login/oauth/authorize?client_id=" + g.cfg.ClientID + "&scope=user"
	return redirectURL
}

func (g *gitHubAPI) SetClient(code string) error {
	// 获取 access token
	token, err := g.getAccessToken(code)
	if err != nil {
		return err
	}

	// 使用 access token 创建 GitHub 客户端
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	g.client = client
	return nil
}

func (g *gitHubAPI) GetUserInfo(ctx context.Context) (*github.User, error) {
	// 获取用户信息
	user, _, err := g.client.Users.Get(ctx, "")
	if err != nil {
		return &github.User{}, err
	}
	return user, nil
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
