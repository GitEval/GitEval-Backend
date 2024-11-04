package conf

import (
	"github.com/google/wire"
	"golang.org/x/oauth2/jwt"
)

var ProviderSet = wire.NewSet(
	NewVipperSetting,
	NewAppConf,
	NewGitHubConfig,
	NewDataConfig,
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
type DataConfig struct {
	Addr string `yaml:"addr"`
}

// 配置结构体
type LLMConfig struct {
	Addr string `yaml:"addr"`
}

type JWTConfig struct {
	SecretKey string `yaml:"secret-key"` //秘钥
	Timeout   int    `yaml:"timeout"`    //过期时间
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
func NewDataConfig(s *VipperSetting) *DataConfig {
	var dataConfig = &DataConfig{}
	s.ReadSection("data", dataConfig)
	return dataConfig
}

func NewLLMConfig(s *VipperSetting) *LLMConfig {
	var llmConfig = &LLMConfig{}
	s.ReadSection("llm", llmConfig)
	return llmConfig
}
func NewJWTConfig(s *VipperSetting) *jwt.Config {
	var jwtConf = &JWTConfig{}
	s.ReadSection("jwt", jwtConf)
	return jwtConf
}
