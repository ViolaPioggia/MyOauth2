package global

import (
	"MyOauth2/config/config"
	"gorm.io/gorm"
)

var (
	Config  *config.Config
	MysqlDB *gorm.DB
	//Rdb     *redis.Client
)
