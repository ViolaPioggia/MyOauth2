package global

import (
	"MyOauth2/config/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	Config  *config.Config
	MysqlDB *gorm.DB
	Rdb     *redis.Client
)
