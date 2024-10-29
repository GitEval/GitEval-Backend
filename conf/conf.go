package conf

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewVipperSetting,
	NewAppConf,
	NewGitHubConfig,
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
