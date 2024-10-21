package system

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

// ReceiveMessage struct
type ReceiveMessage struct {
	UpdateID    int         `json:"update_id"`
	Message     Message     `json:"message"`
	ChannelPost ChannelPost `json:"channel_post"`
}

// Message struct
type Message struct {
	MessageID int        `json:"message_id"`
	From      From       `json:"from"`
	Chat      Chat       `json:"chat"`
	Date      int        `json:"date"`
	Text      string     `json:"text"`
	Entities  []Entities `json:"entities"`
}

// ChannelPost struct
type ChannelPost struct {
	MessageID int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

// SendMessage struct
type SendMessage struct {
	Ok     bool   `json:"ok"`
	Result Result `json:"result"`
}

// Result struct
type Result struct {
	MessageID int    `json:"message_id"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
}

// From struct
type From struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	UserName     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

// Chat struct
type Chat struct {
	ID                          int    `json:"id"`
	FirstName                   string `json:"first_name"`
	UserName                    string `json:"username"`
	Type                        string `json:"type"`
	Title                       string `json:"title"`
	AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
}

// Entities struct
type Entities struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

// Update struct
type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

// WebhookRequest struct 处理webhook请求的结构体 响应群组中的消息
type WebhookRequest struct {
	UpdateID int              `json:"update_id"`
	Message  tgbotapi.Message `json:"message"`
}

// JenkinsBuild struct 构建项目时的视图名和项目名
type JenkinsBuild struct {
	ViewName string `json:"view_name"` // 视图名称
	JobName  string `json:"job_name"`  // 项目名称
	TaskType string `json:"task_type"` // 任务类型  如 后台API 前台API
}

type AdminLoginToken struct {
	ID        uint   `gorm:"primary_key;auto_increment" json:"id"`
	HttpToken string `gorm:"type:text" json:"http_token"`
	CreatedAt time.Time
}

func (AdminLoginToken) TableName() string {
	return "tg_admin_token"
}
