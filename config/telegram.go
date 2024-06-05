package config

type telegram struct {
	BotToken string `mapstructure:"BotToken" yaml:"bot-token"`
	ChatID   string `mapstructure:"ChatID" yaml:"chat-id"`
}
