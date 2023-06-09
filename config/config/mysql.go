package config

import (
	"fmt"
	"time"
)

type Mysql struct {
	Addr            string `mapstructure:"addr" yaml:"addr"`
	Port            string `mapstructure:"port" yaml:"port"`
	Db              string `mapstructure:"db" yaml:"db"`
	Username        string `mapstructure:"username" yaml:"username"`
	Password        string `mapstructure:"password" yaml:"password"`
	Charset         string `mapstructure:"charset" yaml:"charset"`
	ConnMaxIdleTime string `mapstructure:"connMaxIdleTime" yaml:"connMaxIdleTime"`
	ConnMaxLifeTime string `mapstructure:"connMaxLifeTime" yaml:"connMaxLifeTime"`
	Place           string `mapstructure:"place" yaml:"place"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenCoons" yaml:"maxOpenConns"`
}

func (m *Mysql) GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s",
		m.Username,
		m.Password,
		m.Addr,
		m.Port,
		m.Db,
		m.Charset,
		m.Place,
	)
}

func (m *Mysql) GetConnMaxIDleTime() time.Duration {
	t, _ := time.ParseDuration(m.ConnMaxIdleTime)
	return t
}

func (m *Mysql) GetconnMaxLifeTime() time.Duration {
	t, _ := time.ParseDuration(m.ConnMaxLifeTime)
	return t
}
