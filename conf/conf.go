package conf

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewVipperSetting,
	NewAppConf,
	NewGitHubConfig,
	NewLLMConfig,
)

type AppConf struct {
	Addr string `json:"addr"`
	//其他配置也可以加到这个里面
}

// GitHubConfig 使用统一的cfg管理方案
type GitHubConfig struct {
	ClientID     string `yaml:"client-id"`
	ClientSecret string `yaml:"client-secret"`
}

// 配置结构体
type LLMConfig struct {
	Addr string `yaml:"addr"`
}

func NewAppConf(s *VipperSetting) *AppConf {
	var appconf = &AppConf{}
	s.ReadSection("app", appconf)
	return appconf
}

func NewGitHubConfig(s *VipperSetting) *GitHubConfig {
	var GitHubConf = &GitHubConfig{}
	s.ReadSection("github", GitHubConf)
	return GitHubConf
}

// NewLLMClient 创建新的 LLMClient 实例
func NewLLMConfig(s *VipperSetting) *LLMConfig {
	var LLMConf = &LLMConfig{}
	s.ReadSection("github", LLMConf)
	return LLMConf
}
