package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

// LLMClient 接口定义
type LLMClient interface {
	GetDomain(ctx context.Context, req GetDomainRequest) (GetDomainResponse, error)
	GetEvaluation(ctx context.Context, req GetEvaluationRequest) (GetEvaluationResponse, error)
}

// lLMClient 结构体
type lLMClient struct {
	cfg *llmConfig
}

// 配置结构体
type llmConfig struct {
	Addr string `yaml:"addr"`
}

// NewLLMClient 创建新的 LLMClient 实例
func NewLLMClient() LLMClient {
	var cfg llmConfig
	err := viper.UnmarshalKey("llm", &cfg)
	if err != nil {
		return nil
	}
	return &lLMClient{cfg: &cfg}
}

// GetEvaluationResponse 结构体与 Python 的 EvaluationResponse 对应
type GetEvaluationResponse struct {
	Evaluation string `json:"evaluation"` // 响应消息内容
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
	Name     string             `json:"name"`     // 仓库名称
	Readme   string             `json:"readme"`   // README 内容
	Language map[string]float64 `json:"language"` // 编程语言及其使用比例
	Commit   int                `json:"commit"`   // 用户的 commit 次数
	Add      int                `json:"add"`      // 增加的 commit 数量
	Delete   int                `json:"delete"`   // 减少的 commit 数量
	Star     int                `json:"star"`     // 被 star 的数量
	Fork     int                `json:"fork"`     // 被 fork 的数量
}

// DomainRequest 结构体与 Python 的 DomainRequest 对应
type GetDomainRequest struct {
	Repos         []Repo   `json:"repos"`         // 仓库列表
	Bio           string   `json:"bio"`           // 个人简介
	Organizations []string `json:"organizations"` // 所属组织
}

// DomainResponse 结构体与 Python 的 DomainResponse 对应
type GetDomainResponse struct {
	Domain []string `json:"domain"` // 响应消息内容
}

// EvaluationRequest 结构体与 Python 的 EvaluationRequest 对应
type GetEvaluationRequest struct {
	Bio               string   `json:"bio"`                 // 个人简介
	Followers         int      `json:"followers"`           // 粉丝
	Following         int      `json:"following"`           // 关注
	TotalPrivateRepos int      `json:"total_private_repos"` // 私人仓库数量
	TotalPublicRepos  int      `json:"total_public_repos"`  // 公开仓库数量
	Repos             []Repo   `json:"repos"`               // 仓库列表
	Domains           []string `json:"domains"`             // 技术领域
	CreatedAt         string   `json:"created_at"`          // 建号时间
	Organizations     []string `json:"organizations"`       // 所属组织
	DiskUsage         float64  `json:"disk_usage"`          // 硬盘使用量单位是KB
}
