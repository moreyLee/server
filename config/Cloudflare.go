package config

type Cloudflare struct {
	CfApiLogin string `mapstructure:"cf-api-login" json:"cf-api-login" yaml:"cf-api-login"`
	GlobalKey  string `mapstructure:"" json:"global-key" yaml:"global-key"`
	ApiUrl     string `mapstructure:"" json:"api-url" yaml:"api-url"`
}
