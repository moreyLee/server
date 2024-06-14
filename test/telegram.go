package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

var (
	botToken = "7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4"
	// 测试群组ID
	groupID = int64(-4275796428)
)

func Send() {

	// 初始化机器人
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	// 启用调试模式
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// 创建一个新的消息
	chatID := groupID // 替换为目标聊天 ID（负数表示群组）
	messageText := "群组消息"
	// 发送消息
	msg := tgbotapi.NewMessage(chatID, messageText)

	// 发送消息
	_, err = bot.Send(msg)

	log.Printf("Message sent to chat ID %d", chatID)

	// 设置更新配置
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// 获取更新通道
	updates := bot.GetUpdatesChan(u)

	// 处理更新
	for update := range updates {
		if update.Message == nil { // 忽略任何非消息更新
			continue
		}
		// 打印收到的消息
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// 检查消息是否提到了机器人
		if update.Message.IsCommand() || strings.Contains(update.Message.Text, "@"+bot.Self.UserName) {
			//if strings.HasPrefix(msgText, "@CG33333_bot") {
			// 回复消息
			responseText := "你提到我了吗？我在这里！"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
			msg.ReplyToMessageID = update.Message.MessageID

			// 发送回复消息
			bot.Send(msg)

		}
	}
}

func main() {
	//Send()
	//GetchatID()
}
