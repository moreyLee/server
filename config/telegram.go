package config

type Telegram struct {
	URL            string   `mapstructure:"url" yaml:"url"`
	BotName        string   `mapstructure:"bot-name" yaml:"bot-name"`
	BotToken       string   `mapstructure:"bot-token" yaml:"bot-token"`
	ChatID         int64    `mapstructure:"chat-id" yaml:"chat-id"`
	WebhookUrl     string   `mapstructure:"webhook-url" yaml:"webhook-url"`
	WebhookPort    string   `mapstructure:"webhook-port" yaml:"webhook-port"`
	AuthorizedUser []string `mapstructure:"authorized_users" yaml:"authorized_users"`
}
