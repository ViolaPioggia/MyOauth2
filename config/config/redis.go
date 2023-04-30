package config

type Redis struct {
	Addr     string `mapstructure:"addr" yaml:"addr"`
	Port     string `mapstructure:"port" yaml:"port"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	Db       int    `mapstructure:"db" yaml:"db"`
	PoolSize int    `mapstructure:"poolSize" yaml:"poolSize"`
}
