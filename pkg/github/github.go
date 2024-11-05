package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/GitEval/GitEval-Backend/conf"
	"github.com/GitEval/GitEval-Backend/model"
	"github.com/GitEval/GitEval-Backend/pkg/github/expireMap"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	"log"
	"strings"
	"time"
)

const (
	ExpireTime = time.Hour * 24 * 7
)

// GitHubAPI 结构体
// 将其当作处理所有有关github账号的中枢,因为它有map
type GitHubAPI struct {
	clients expireMap.ExpireMap // 使用 sync.Map 实现并发安全
	cfg     *conf.GitHubConfig  // 引用的地址完全相同节约了内存空间
}

func NewGitHubAPI(c *conf.GitHubConfig, clients expireMap.ExpireMap) *GitHubAPI {
	return &GitHubAPI{
		cfg:     c,
		clients: clients,
	}
}

// SetClient 设置用户的 GitHub 客户端
func (g *GitHubAPI) SetClient(userID int64, client *github.Client) {
	g.clients.Store(userID, client, ExpireTime) // 使用 Store 方法
}

// GetClientFromMap GetClient 获取用户的 GitHub 客户端
func (g *GitHubAPI) GetClientFromMap(userID int64) (*github.Client, bool) {
	client, exists := g.clients.Load(userID) // 使用 Load 方法
	if exists {
		return client.(*github.Client), true // 类型断言
	}
	return nil, false
}

func (g *GitHubAPI) GetLoginUrl() string {
	redirectURL := "https://github.com/login/oauth/authorize?client_id=" + g.cfg.ClientID + "&scope=user"
	return redirectURL
}

func (g *GitHubAPI) GetClientByCode(code string) (*github.Client, error) {
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

func (g *GitHubAPI) GetUserInfo(ctx context.Context, client *github.Client, username string) (*github.User, error) {
	userInfo, _, err := client.Users.Get(ctx, username)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (g *GitHubAPI) GetFollowing(ctx context.Context, id int64) []model.User {
	val, exist := g.clients.Load(id)
	if !exist {
		log.Println("get github client failed")
		return []model.User{}
	}
	client := val.(*github.Client)
	users, _, err := client.Users.ListFollowing(ctx, "", nil)
	if err != nil {
		log.Println("get github following user failed:", err)
		return []model.User{}
	}

	// 获取详细用户信息
	var detailedUsers []model.User
	for _, user := range users {
		detailedUser, _, err := client.Users.Get(ctx, *user.Login)
		if err != nil {
			log.Println("get user details failed:", err)
			continue
		}
		detailedUsers = append(detailedUsers, model.TransformUser(detailedUser))
	}

	return detailedUsers
}

func (g *GitHubAPI) GetFollowers(ctx context.Context, id int64) []model.User {
	val, exist := g.clients.Load(id)
	if !exist {
		log.Println("get github client failed")
		return []model.User{}
	}
	client := val.(*github.Client)
	users, _, err := client.Users.ListFollowers(ctx, "", nil)
	if err != nil {
		log.Println("get github followers user failed:", err)
		return []model.User{}
	}

	// 获取详细用户信息
	var detailedUsers []model.User
	for _, user := range users {
		detailedUser, _, err := client.Users.Get(ctx, *user.Login)
		if err != nil {
			log.Println("get user details failed:", err)
			continue
		}
		detailedUsers = append(detailedUsers, model.TransformUser(detailedUser))
	}

	return detailedUsers
}

func (g *GitHubAPI) CalculateScore(ctx context.Context, id int64, name string) float64 {
	var client = &github.Client{}
	val, exist := g.clients.Load(id)
	if !exist {
		// 创建一个 GitHub 客户端（无需认证）
		client = github.NewClient(nil) // nil 表示没有使用任何认证
	} else {
		client = val.(*github.Client)
	}
	repos, _, err := client.Repositories.List(ctx, name, nil)
	if err != nil {
		log.Printf("Error getting repositories: %v\n", err)
		return 0
	}
	// 计算评分
	score := calculateScore(repos)
	return score
}

// GetReposDetailList 根据仓库链接获取仓库的详细信息列表
func (g *GitHubAPI) GetRepoDetail(ctx context.Context, repoUrl string, client *github.Client) (*github.Repository, error) {

	// 提取用户名和仓库名
	owner, repo, err := g.parseRepoURL(repoUrl)
	if err != nil {
		return &github.Repository{}, fmt.Errorf("invalid repository URL %s: %v", repoUrl, err)
	}

	// 获取仓库的详细信息
	repository, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return &github.Repository{}, fmt.Errorf("failed to get repository details for %s: %v", repoUrl, err)
	}

	return repository, nil
}

// GetAllRepositories 获取用户的所有仓库信息
// 接受用户的昵称和userID,返回所有仓库信息
func (g *GitHubAPI) GetAllRepositories(ctx context.Context, loginName string, userId int64) []*model.Repo {
	v, exist := g.clients.Load(userId)
	if !exist {
		log.Println("get github client failed")
		return nil
	}
	client := v.(*github.Client)
	repos, _, err := client.Repositories.List(ctx, loginName, &github.RepositoryListOptions{
		Sort:        "created",                       // 按创建时间排序
		Direction:   "desc",                          // 降序排列，从新到旧
		ListOptions: github.ListOptions{PerPage: 20}, // 每页最多获取20个
	})
	if err != nil {
		log.Printf("Error getting repositories: %v\n", err)
		return nil
	}
	var resp []*model.Repo
	for _, repo := range repos {
		//尝试获取每个仓库的Readme
		me, err := g.GetReadMe(ctx, repo.GetURL(), client)
		if err != nil {
			log.Println("get github readme failed:", err)
		}
		resp = append(resp, &model.Repo{
			Name:     repo.Name,
			Readme:   &me,
			Language: repo.Language,
		})
	}
	return resp
}

func (g *GitHubAPI) GetReadMe(ctx context.Context, repoUrl string, client *github.Client) (readme string, err error) {
	// 提取用户名和仓库名
	owner, repo, err := g.parseRepoURL(repoUrl)
	if err != nil {
		return "", fmt.Errorf("invalid repository URL %s: %v", repoUrl, err)
	}
	//获取readme
	gitReadme, _, err := client.Repositories.GetReadme(ctx, owner, repo, nil)
	if err != nil || gitReadme == nil || gitReadme.Content == nil {
		return "", err
	}

	// 对 base64 内容解码
	content, err := base64.StdEncoding.DecodeString(*gitReadme.Content)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (g *GitHubAPI) GetOrganizations(ctx context.Context, userID int64) ([]*github.Organization, error) {
	val, exist := g.clients.Load(userID)
	if !exist {
		log.Println("get github client failed")
		return nil, fmt.Errorf("github client not found for user ID %d", userID)
	}
	client := val.(*github.Client)

	// 获取组织列表
	orgs, _, err := client.Organizations.List(ctx, "", nil)
	if err != nil {
		log.Println("Error getting organizations:", err)
		return nil, err
	}

	return orgs, nil
}

func (g *GitHubAPI) GetAllUserEvents(ctx context.Context, username string, client *github.Client) ([]model.UserEvent, error) {
	allEvents := make([]*github.Event, 0)

	// 分页设置
	opt := &github.ListOptions{PerPage: 100}

	// 循环获取所有用户事件
	for {
		// 获取用户事件
		events, resp, err := client.Activity.ListEventsPerformedByUser(ctx, username, false, opt)
		if err != nil {
			return nil, err // 返回nil而不是UserEvent{}，因为我们要返回切片
		}

		allEvents = append(allEvents, events...)

		// 如果没有更多页面，则退出循环
		if resp.NextPage == 0 {
			break
		} else if len(allEvents) >= 2000 {
			//如果事件过多的话就做限流,这个地方还得再明确下
			break
		}

		// 更新分页选项以请求下一页
		opt.Page = resp.NextPage
	}

	// 使用一个映射来分类不同的UserEvent
	userEventsMap := make(map[string]*model.UserEvent)
	info, err := g.getUserAllRepoInfo(ctx, client, "")
	if err != nil {
		return nil, err
	}

	for _, event := range allEvents {
		repoName := event.Repo.GetName()

		// 如果该repo的UserEvent还未创建，则初始化
		if _, exists := userEventsMap[repoName]; !exists {
			//尝试初始化
			userEventsMap[repoName] = &model.UserEvent{Repo: model.RepoInfo{Name: &repoName}}
			//如果存在于用户的仓库中则直接完全初始化这个仓库
			if _, ok := info[repoName]; ok {
				userEventsMap[repoName].Repo = *info[repoName]
			}
		}

		userEvent := userEventsMap[repoName] // 获取当前repo的UserEvent实例
		switch event.GetType() {
		case "PushEvent":
			// 更新提交计数
			userEvent.PushCount++

		case "IssuesEvent":
			// 记录 issue 信息
			userEvent.IssuesCount++

		case "PullRequestEvent":
			// 记录 Pull Request 信息
			userEvent.PullRequestCount++

		}
	}

	// 将userEventsMap转换为切片
	userEventsSlice := make([]model.UserEvent, 0, len(userEventsMap))
	for _, userEvent := range userEventsMap {
		userEventsSlice = append(userEventsSlice, *userEvent) // 将指针解引用
	}

	return userEventsSlice, nil
}

func (g *GitHubAPI) getUserAllRepoInfo(ctx context.Context, client *github.Client, username string) (map[string]*model.RepoInfo, error) {
	// 创建一个 map 用于存储仓库名称和对应的仓库信息
	repoMap := make(map[string]*model.RepoInfo)
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		// 获取用户的仓库信息
		repos, resp, err := client.Repositories.List(ctx, username, opt)

		if err != nil {
			return nil, err
		}

		// 将获取到的仓库信息存入 map
		for _, repo := range repos {
			if repo.Name != nil {
				created_at := repo.CreatedAt.Time.String()
				name := repo.GetFullName()
				repoMap[name] = &model.RepoInfo{
					Name:             &name,
					Description:      repo.Description,
					StargazersCount:  repo.StargazersCount,
					ForksCount:       repo.ForksCount,
					CreatedAt:        &created_at,
					SubscribersCount: repo.SubscribersCount,
				}
			}

		}

		// 如果没有更多页面，则退出循环,如果map的数量过于庞大自动退出
		if resp.NextPage == 0 || len(repoMap) == 1000 {
			break
		}

		// 更新分页选项以请求下一页
		opt.Page = resp.NextPage
	}

	return repoMap, nil
}

// parseRepoURL 从仓库链接中解析出用户名和仓库名
func (g *GitHubAPI) parseRepoURL(url string) (owner, repo string, err error) {
	parts := strings.Split(strings.TrimPrefix(url, "https://api.github.com/repos/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub repository URL")
	}
	return parts[0], parts[1], nil
}

// 具体的计算逻辑
func calculateScore(repos []*github.Repository) float64 {
	var totalScore float64
	for _, repo := range repos {
		if repo.StargazersCount != nil && repo.ForksCount != nil && repo.Size != nil {
			// 评分公式示例
			score := float64(*repo.StargazersCount)*0.5 + float64(*repo.ForksCount)*0.3 + float64(*repo.Size/1024)*0.2
			totalScore += score
		}
	}
	return totalScore
}

func (g *GitHubAPI) getAccessToken(code string) (string, error) {
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

	// TODO 获取访问令牌,这个地方经常出现请求超时,需要研究下原因
	token, err := cf.Exchange(ctx, code)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
