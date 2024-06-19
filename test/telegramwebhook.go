package main

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
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

// 接收消息
var (
	URL  = "https://api.telegram.org/bot"
	port = "5299"
	//webhookURL = "https://devops.3333d.vip/telegram-webhook"
	webhookUrl = "https://e692-91-75-118-214.ngrok-free.app/telegram-webhook"
	token      = "7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	message := system.ReceiveMessage{}

	chatID := 0
	msgText := ""

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		fmt.Println(err)
	}
	// if private or group
	if message.Message.Chat.ID != 0 {
		fmt.Println("群组ID与消息文本", message.Message.Chat.ID, message.Message.Text)
		// 获取到消息的chatID 消息内容
		chatID = message.Message.Chat.ID
		msgText = message.Message.Text
	} else {
		// if channel
		fmt.Println("频道id和频道文本", message.ChannelPost.Chat.ID, message.ChannelPost.Text)
		//chatID = message.ChannelPost.Chat.ID
		//msgText = message.ChannelPost.Text
	}

	respMsg := fmt.Sprintf("%s%s/sendMessage?chat_id=%d&text=Received: %s", URL, token, chatID, msgText)

	// send echo resp
	_, err = http.Get(respMsg)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 初始化 bot 实例
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	// 开启debug 模式
	bot.Debug = true
	// 创建一个webhook 网络钩子
	webhook, err := tgbotapi.NewWebhook(webhookUrl)
	// 对webhook实例发起请求
	_, err = bot.Request(webhook)
	// 获取webhookInfo 返回的请求状态
	info, err := bot.GetWebhookInfo()
	// 如果请求失败，打印错误信息
	if info.LastErrorDate != 0 {
		log.Printf("Telegram webhook 初始化报错: %s\n", info.LastErrorMessage)

	}
	// 遍历收到的消息
	update := bot.ListenForWebhook("/")

	http.HandleFunc("/telegram-webhook", Webhook)
	log.Printf("Starting server on port %s", port)
	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	for v := range update {
		if v.Message != nil {
			continue
		}
		msg := tgbotapi.NewMessage(v.Message.Chat.ID, "1")
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}
