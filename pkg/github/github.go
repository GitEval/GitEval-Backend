package github

import (
	"context"
	"github.com/GitEval/GitEval-Backend/conf"
	"github.com/GitEval/GitEval-Backend/model"
	"github.com/GitEval/GitEval-Backend/pkg/github/expireMap"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	"log"
	"time"
)

const (
	ExpireTime = time.Hour * 24 * 7
)

//cc:
//接口的声明一般由调用方来决定,比如你的authService只是用到了几个方法,就没必要直接将该接口给它
//这里我懒得改了，哈哈哈

type GitHubAPI interface {
	GetLoginUrl() string
	SetClient(userID int64, client *github.Client)
	GetClientFromMap(userID int64) (*github.Client, bool)
	GetClientByCode(code string) (*github.Client, error)
	GetUserInfo(ctx context.Context, client *github.Client, username string) (*github.User, error)
}

// gitHubAPI 结构体
// 将其当作处理所有有关github账号的中枢,因为它有map
type gitHubAPI struct {
	clients expireMap.ExpireMap // 使用 sync.Map 实现并发安全
	cfg     *conf.GitHubConfig  // 引用的地址完全相同节约了内存空间
}

func NewGitHubAPI(c *conf.GitHubConfig, clients expireMap.ExpireMap) GitHubAPI {
	return &gitHubAPI{
		cfg:     c,
		clients: clients,
	}
}

// SetClient 设置用户的 GitHub 客户端
func (g *gitHubAPI) SetClient(userID int64, client *github.Client) {
	g.clients.Store(userID, client, ExpireTime) // 使用 Store 方法
}

// GetClientFromMap GetClient 获取用户的 GitHub 客户端
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

func (g *gitHubAPI) GetUserInfo(ctx context.Context, client *github.Client, username string) (*github.User, error) {
	userInfo, _, err := client.Users.Get(ctx, username)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (g *gitHubAPI) GetFollowing(ctx context.Context, id int64) []model.User {
	val, exist := g.clients.Load(id)
	if !exist {
		log.Println("get github client failed")
		return nil
	}
	client := val.(*github.Client)
	users, _, err := client.Users.ListFollowing(ctx, "", nil)
	if err != nil {
		log.Println("get github following user failed")
		return nil
	}
	return model.TransformUsers(users)
}

func (g *gitHubAPI) GetFollowers(ctx context.Context, id int64) []model.User {
	val, exist := g.clients.Load(id)
	if !exist {
		log.Println("get github client failed")
		return nil
	}
	client := val.(*github.Client)
	users, _, err := client.Users.ListFollowers(ctx, "", nil)
	if err != nil {
		log.Println("get github following user failed")
		return nil
	}
	return model.TransformUsers(users)
}
func (g *gitHubAPI) CalculateScore(ctx context.Context, id int64, name string) float64 {
	val, exist := g.clients.Load(id)
	if !exist {
		log.Println("get github client failed")
		return 0
	}
	client := val.(*github.Client)
	// 获取用户的公开仓库
	repos, _, err := client.Repositories.List(ctx, name, nil)
	if err != nil {
		log.Printf("Error getting repositories: %v\n", err)
		return 0
	}
	// 计算评分
	score := calculateScore(repos)
	return score
}

// 具体的计算逻辑
func calculateScore(repos []*github.Repository) float64 {
	var totalScore float64
	for _, repo := range repos {
		if repo.StargazersCount != nil && repo.ForksCount != nil && repo.Size != nil {
			// 评分公式示例
			score := float64(*repo.StargazersCount)*0.5 + float64(*repo.ForksCount)*0.3 + float64(*repo.Size)*0.2
			totalScore += score
		}
	}
	return totalScore
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
