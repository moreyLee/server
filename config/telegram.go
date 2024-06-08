package config

type Telegram struct {
	BotToken string `mapstructure:"bot-token" yaml:"bot-token"`
	ChatID   int64  `mapstructure:"chat-id" yaml:"chat-id"`
}
