package github

type RepoInfo struct {
	Description      string `json:"Description"`
	StargazersCount  int    `json:"StargazersCount"`
	ForksCount       int    `json:"ForksCount"`
	CreatedAt        string `json:"CreatedAt"`
	SubscribersCount int    `json:"SubscribersCount"`
}

type Repo struct {
	Name     string `json:"name"`     // 仓库名称
	Readme   string `json:"readme"`   // README 内容
	Language string `json:"language"` // 使用最多的编程语言
	Commit   int    `json:"commit"`   // 用户的 commit 次数
	Star     int    `json:"star"`     // 被 star 的数量
	Fork     int    `json:"fork"`     // 被 fork 的数量
}

type UserEvent struct {
	Repo             *RepoInfo `json:"repo"`
	CommitCount      int       `json:"commit_count"`
	Commit           []string  `json:"commit"`
	IssuesCount      int       `json:"issues_count"`
	Issues           []string  `json:"issues"`
	PullRequest      []string  `json:"pull_request"`
	PullRequestCount int       `json:"pull_request_count"`
}

type IssuesEventPayload struct {
	Action string `json:"action"`
	Issue  struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"issue"`
}

type PullRequestEventPayload struct {
	Action      string `json:"action"`
	PullRequest struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"pull_request"`
}
