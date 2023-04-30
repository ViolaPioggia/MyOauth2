package config

import (
	"MyOauth2/config/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlDbSetup() {
	config := global.Config.Mysql

	db, err := gorm.Open(mysql.Open(config.GetDsn()), &gorm.Config{})
	if err != nil {
		panic("initialize mysql failed")
	}

	global.MysqlDB = db
}
