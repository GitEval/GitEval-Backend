package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

type VipperSetting struct {
	*viper.Viper
}

func (s *VipperSetting) ReadSection(k string, v interface{}) error {
	err := s.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}

func NewVipperSetting(ConfFilePath string) *VipperSetting {
	vp := viper.New()
	vp.SetConfigFile(ConfFilePath) // 指定配置文件路径
	err := vp.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Read confg err:%v", err))
	}
	return &VipperSetting{
		Viper: vp,
	}
}
