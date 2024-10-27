package conf

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewVipperSetting,
	NewAppConf,
)

type AppConf struct {
	Addr string `json:"addr"`
	//其他配置也可以加到这个里面
}

func NewAppConf(s *VipperSetting) *AppConf {
	var appconf = &AppConf{}
	s.ReadSection("app", appconf)
	return appconf
}
