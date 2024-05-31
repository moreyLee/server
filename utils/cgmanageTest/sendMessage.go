package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func SendMessage() {
	botToken := "7005107845:AAEWU9OmtzLa6YHAHROAG3wODYrUh8opFbw"

	// 初始化机器人
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	// 启用调试模式
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// 创建一个新的消息
	chatID := int64(-4275796428) // 替换为目标聊天 ID（负数表示群组）
	messageText := "开始"
	// 发送消息
	msg := tgbotapi.NewMessage(chatID, messageText)

	// 发送消息
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}

	// 设置更新配置
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// 获取更新通道
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // 忽略任何非消息更新
			continue
		}
		// 打印收到的消息
		log.Printf("打印收到的消息内容:[%s] %s", update.Message.From.UserName, update.Message.Text)
		// 检查消息是否提到了机器人
		if strings.Contains(update.Message.Text, "@"+bot.Self.UserName) {
			//
			// 创建一个新的消息  回复发送的消息
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			// 发送回复消息
			_, err = bot.Send(msg)
			if err != nil {
				log.Panic(err)
			}
		}
		fmt.Println("用户输入的命令回复的消息:", update.Message.Text)
	}
}

func main() {
	SendMessage()
}
