package config

type Jenkins struct {
	Url      string `mapstructure:"url" yaml:"url"`
	User     string `mapstructure:"user" yaml:"user"`
	ApiToken string `mapstructure:"api-token" yaml:"api-token"`
}
