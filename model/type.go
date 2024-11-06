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
	Name     string `json:"name"`     // 仓库名称
	Readme   string `json:"readme"`   // README 内容
	Language string `json:"language"` // 使用最多的编程语言
}

type UserEvent struct {
	Repo             RepoInfo `json:"repo"`
	PushCount        int      `json:"push_count"`
	IssuesCount      int      `json:"issues_count"`
	PullRequestCount int      `json:"pull_request_count"`
}

//弃用了,涉及的东西太多,评价体系太复杂,效果还一般
//type IssuesEventPayload struct {
//	Action string `json:"action"`
//	Issue  struct {
//		Title string `json:"title"`
//		Body  string `json:"body"`
//	} `json:"issue"`
//}
//
//type PullRequestEventPayload struct {
//	Action      string `json:"action"`
//	PullRequest struct {
//		Title string `json:"title"`
//		Body  string `json:"body"`
//	} `json:"pull_request"`
//}
