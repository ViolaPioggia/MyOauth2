package config

type Config struct {
	Mysql *Mysql `mapstructure:"mysql" yaml:"mysql"`
	//Redis   *Redis   `mapstructure:"redis" yaml:"redis"`
	Oauth2  *Oauth2  `mapstructure:"oauth2" yaml:"oauth2"`
	Session *Session `mapstructure:"session" yaml:"session"`
}
