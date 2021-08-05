package initialize

import (
	"crontab/config"
	"crontab/global"

	"github.com/spf13/viper"
)

func InitConfig() {
	v := viper.New()
	v.SetConfigFile("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	serverConfig := config.ServerConfig{}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	global.Settings = serverConfig
}
