package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/GitEval/GitEval-Backend/conf"
	"net/http"
)

// LLMClient 结构体
type LLMClient struct {
	cfg *conf.LLMConfig
}

// GetEvaluationResponse 结构体与 Python 的 EvaluationResponse 对应
type GetEvaluationResponse struct {
	Evaluation string `json:"evaluation"` // 响应消息内容
}

func NewLLMClient(cfg *conf.LLMConfig) *LLMClient {
	return &LLMClient{cfg: cfg}
}

// GetDomain 发送请求以获取领域
func (c *LLMClient) GetDomain(ctx context.Context, req GetDomainRequest) (GetDomainResponse, error) {
	url := fmt.Sprintf("%s/getDomain", c.cfg.Addr)
	resp, err := c.sendPostRequest(ctx, url, req)
	if err != nil {
		return GetDomainResponse{}, err
	}
	defer resp.Body.Close()
	// 处理响应
	var response GetDomainResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return GetDomainResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}

// GetEvaluation 发送请求以获取评估
func (c *LLMClient) GetEvaluation(ctx context.Context, req GetEvaluationRequest) (GetEvaluationResponse, error) {
	url := fmt.Sprintf("%s/getEvaluation", c.cfg.Addr)

	resp, err := c.sendPostRequest(ctx, url, req)
	if err != nil {
		return GetEvaluationResponse{}, err
	}
	defer resp.Body.Close()

	// 处理响应
	var response GetEvaluationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return GetEvaluationResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}

// GetEvaluation 发送请求以获取评估
func (c *LLMClient) GetArea(ctx context.Context, req GetAreaRequest) (GetAreaResponse, error) {
	url := fmt.Sprintf("%s/getArea", c.cfg.Addr)

	resp, err := c.sendPostRequest(ctx, url, req)
	if err != nil {
		return GetAreaResponse{}, err
	}
	defer resp.Body.Close()
	// 处理响应
	var response GetAreaResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return GetAreaResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}

// sendPostRequest 通用的 POST 请求发送函数
func (c *LLMClient) sendPostRequest(ctx context.Context, url string, req interface{}) (*http.Response, error) {
	// 将请求体序列化为 JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建新的 HTTP 请求
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}

// Repo 结构体与 Python 的 Repo 对应
type Repo struct {
	Name         *string `json:"name"`     // 仓库名称
	MainLanguage *string `json:"language"` // 主要的编程语言
	Readme       *string `json:"readme"`   //这个仓库的readme
}

// DomainRequest 结构体与 Python 的 DomainRequest 对应
type GetDomainRequest struct {
	Repos []Repo `json:"repos"` // 仓库列表
	Bio   string `json:"bio"`   // 个人简介
}

// DomainResponse 结构体与 Python 的 DomainResponse 对应
type GetDomainResponse struct {
	Domains []Domain `json:"domains"`
}
type Domain struct {
	Domain     string  `json:"domain"`
	Confidence float64 `json:"confidence"`
}
type RepoInfo struct {
	Name             *string `json:"name,omitempty"`
	Description      *string `json:"description,omitempty"`
	StargazersCount  *int    `json:"stargazers_count,omitempty"`
	ForksCount       *int    `json:"forks_count,omitempty"`
	CreatedAt        *string `json:"created_at,omitempty"`
	SubscribersCount *int    `json:"subscribers_count,omitempty"`
}

type UserEvent struct {
	Repo             *RepoInfo `json:"repo,omitempty"`
	CommitCount      int       `json:"commit_count"`
	IssuesCount      int       `json:"issues_count"`
	PullRequestCount int       `json:"pull_request_count"`
}

type GetEvaluationRequest struct {
	Bio               *string     `json:"bio,omitempty"`
	Followers         int         `json:"followers"`
	Following         int         `json:"following"`
	TotalPrivateRepos int         `json:"total_private_repos"`
	TotalPublicRepos  int         `json:"total_public_repos"`
	UserEvents        []UserEvent `json:"user_events"`
	Domains           *[]string   `json:"domains"`
}

// AreaRequest 结构体与 Python 的 AreaRequest 对应
type GetAreaRequest struct {
	Bio            *string  `json:"bio"`             // 个人简介
	Company        *string  `json:"company"`         // 公司
	Location       *string  `json:"location"`        // 地点
	FollowerAreas  []string `json:"follower_areas"`  // 粉丝地区
	FollowingAreas []string `json:"following_areas"` // 关注者地区
}

// AreaResponse 结构体与 Python 的 AreaResponse 对应
type GetAreaResponse struct {
	Area       string  `json:"area"`       // 区域
	Confidence float64 `json:"Confidence"` // 置信度
}
