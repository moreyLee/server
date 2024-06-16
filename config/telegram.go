package config

type Telegram struct {
	URL         string `mapstructure:"url" yaml:"url"`
	BotToken    string `mapstructure:"bot-token" yaml:"bot-token"`
	ChatID      int64  `mapstructure:"chat-id" yaml:"chat-id"`
	WebhookUrl  string `mapstructure:"webhook-url" yaml:"webhook-url"`
	WebhookPort string `mapstructure:"webhook-port" yaml:"webhook-port"`
}
