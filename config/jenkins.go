package config

type Jenkins struct {
	TestUrl   string `mapstructure:"testUrl" yaml:"testUrl"` //
	TestUser  string `mapstructure:"testUser" yaml:"testUser"`
	TestToken string `mapstructure:"testToken" yaml:"testToken"`
	Url       string `mapstructure:"url" yaml:"url"`
	User      string `mapstructure:"user" yaml:"user"`
	ApiToken  string `mapstructure:"api-token" yaml:"api-token"`
}
