package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/GitEval/GitEval-Backend/conf"
	"net/http"
)

// lLMClient 结构体
type lLMClient struct {
	cfg *conf.LLMConfig
}

// GetEvaluationResponse 结构体与 Python 的 EvaluationResponse 对应
type GetEvaluationResponse struct {
	Evaluation string `json:"evaluation"` // 响应消息内容
}

func NewLLMClient(cfg *conf.LLMConfig) *lLMClient {
	return &lLMClient{cfg: cfg}
}

// GetDomain 发送请求以获取领域
func (c *lLMClient) GetDomain(ctx context.Context, req GetDomainRequest) (GetDomainResponse, error) {
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
func (c *lLMClient) GetEvaluation(ctx context.Context, req GetEvaluationRequest) (GetEvaluationResponse, error) {
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
func (c *lLMClient) GetArea(ctx context.Context, req GetAreaRequest) (GetAreaResponse, error) {
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
func (c *lLMClient) sendPostRequest(ctx context.Context, url string, req interface{}) (*http.Response, error) {
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
	Name         string `json:"name"`     // 仓库名称
	MainLanguage string `json:"language"` // 主要的编程语言
}

// DomainRequest 结构体与 Python 的 DomainRequest 对应
type GetDomainRequest struct {
	Repos []Repo `json:"repos"` // 仓库列表
	Bio   string `json:"bio"`   // 个人简介
}

// DomainResponse 结构体与 Python 的 DomainResponse 对应
type GetDomainResponse struct {
	Domain []string `json:"domain"` // 响应消息内容
}

// EvaluationRequest 结构体与 Python 的 EvaluationRequest 对应
type GetEvaluationRequest struct {
	Bio               string      `json:"bio"`                 // 个人简介
	Followers         int         `json:"followers"`           // 粉丝
	Following         int         `json:"following"`           // 关注
	TotalPrivateRepos int         `json:"total_private_repos"` // 私人仓库数量
	TotalPublicRepos  int         `json:"total_public_repos"`  // 公开仓库数量
	UserEvents        []UserEvent `json:"user_events"`         // 用户事件
	Domains           []string    `json:"domains"`             // 技术领域
	CreatedAt         string      `json:"created_at"`          // 建号时间
	Organizations     []string    `json:"organizations"`       // 所属组织
	DiskUsage         float64     `json:"disk_usage"`          // 硬盘使用量单位是KB
}

// UserEvent 结构体与 Python 的 UserEvent 对应
type UserEvent struct {
	Repo             *RepoInfo `json:"repo"`               // 仓库信息
	CommitCount      int       `json:"commit_count"`       // 提交计数
	IssuesCount      int       `json:"issues_count"`       // Issue 计数
	PullRequestCount int       `json:"pull_request_count"` // Pull Request 计数
}

type RepoInfo struct {
	Description      *string `json:"description,omitempty"`       // 仓库描述
	StargazersCount  *int    `json:"stargazers_count,omitempty"`  // Star 数量
	ForksCount       *int    `json:"forks_count,omitempty"`       // Fork 数量
	CreatedAt        *string `json:"created_at,omitempty"`        // 创建时间
	SubscribersCount *int    `json:"subscribers_count,omitempty"` // 订阅者数量
}

// AreaRequest 结构体与 Python 的 AreaRequest 对应
type GetAreaRequest struct {
	Bio            string   `json:"bio"`             // 个人简介
	Company        string   `json:"company"`         // 公司
	Location       string   `json:"location"`        // 地点
	FollowerAreas  []string `json:"follower_areas"`  // 粉丝地区
	FollowingAreas []string `json:"following_areas"` // 关注者地区
}

// AreaResponse 结构体与 Python 的 AreaResponse 对应
type GetAreaResponse struct {
	Area       string  `json:"area"`       // 区域
	Confidence float64 `json:"Confidence"` // 信心值
}
