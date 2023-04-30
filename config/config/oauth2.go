package config

type Oauth2 struct {
	AccessTokenExp    int      `mapstructure:"accessTokenExp" yaml:"accessTokenExp"`
	RefreshTokenExp   int      `mapstructure:"refreshTokenExp" yaml:"refreshTokenExp"`
	IsGenerateRefresh bool     `mapstructure:"isGenerateRefresh" yaml:"isGenerateRefresh"`
	JWTSignedKey      string   `mapstructure:"jwtSignedKey" yaml:"jwtSignedKey"`
	Client            []Client `mapstructure:"client" yaml:"client"`
}

type Client struct {
	ID     string `mapstructure:"id" yaml:"id"`
	Secret string `mapstructure:"secret" yaml:"secret"`
	Name   string `mapstructure:"name" yaml:"name"`
	Addr   string `mapstructure:"addr" yaml:"addr"`
	Port   string `mapstructure:"port" yaml:"port"`
}
