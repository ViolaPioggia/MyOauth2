package config

import (
	"MyOauth2/config/global"
	"fmt"
	"github.com/spf13/viper"
)

func ViperSetup() {
	DbViper := viper.New()
	Oauth2Viper := viper.New()
	DbViper.SetConfigType("yaml")
	DbViper.SetConfigName("db")
	Oauth2Viper.SetConfigType("yaml")
	Oauth2Viper.SetConfigName("oauth2")
	DbViper.AddConfigPath("./config/")
	Oauth2Viper.AddConfigPath("./config/")
	err := DbViper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("get db file failed ,err:%v", err))
	}
	err = Oauth2Viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("get oauth file failed ,err:%v", err))
	}
	if err = Oauth2Viper.Unmarshal(&global.Config); err != nil {
		//将配置文件反序列化到config结构体
		panic(fmt.Errorf("get Oauth2 Unmarshal failed.err:%v", err))
	}
	if err = DbViper.Unmarshal(&global.Config); err != nil {
		//将配置文件反序列化到config结构体
		panic(fmt.Errorf("get Db Unmarshal failed.err:%v", err))
	}
}
