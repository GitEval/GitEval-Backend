package model

type RepoInfo struct {
	Name             string `json:"name"` // 仓库名称
	Description      string `json:"description"`
	StargazersCount  int    `json:"stargazers_count"`
	ForksCount       int    `json:"forks_count"`
	CreatedAt        string `json:"created_at"`
	SubscribersCount int    `json:"subscribers_count"`
}

type Repo struct {
	Name     string `json:"name"`
	Readme   string `json:"readme"`
	Language string `json:"language"` // 使用最多的编程语言
}

type UserEvent struct {
	Repo             RepoInfo `json:"repo"`
	PushCount        int      `json:"push_count"`
	IssuesCount      int      `json:"issues_count"`
	PullRequestCount int      `json:"pull_request_count"`
}
