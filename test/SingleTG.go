package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
)

// 定义一个包级别的变量来保存机器人实例
var botInstance *tgbotapi.BotAPI
var once sync.Once

// GetBotInstance 获取机器人实例的函数
func GetBotInstance(token string) *tgbotapi.BotAPI {
	once.Do(func() {
		var err error
		botInstance, err = tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Panic(err)
		}
	})
	return botInstance
}

//func main() {
//	// 获取单例的机器人实例
//	//bot := GetBotInstance("7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4")
//	bot := GetBotInstance(global.GVA_CONFIG.Telegram.BotToken)
//	log.Printf("token值" + global.GVA_CONFIG.Telegram.BotToken)
//	bot.Debug = false
//	log.Printf("Authorized on account %s", bot.Self.UserName)
//
//	u := tgbotapi.NewUpdate(0)
//	u.Timeout = 60
//
//	updates := bot.GetUpdatesChan(u)
//
//	for update := range updates {
//		if update.Message == nil { // 忽略任何非消息更新
//			continue
//		}
//
//		// 检查消息是否提到机器人
//		if update.Message.IsCommand() {
//			//if strings.HasPrefix(msgText, "@CG33333_bot") {}
//			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
//			switch update.Message.Command() {
//			case "start":
//				msg.Text = "Hello! I am your friendly Telegram bot."
//			case "help":
//				msg.Text = "You can control me by sending these commands:\n/start - to start the bot\n/help - to get this help message"
//			default:
//				msg.Text = "I don't know that command"
//			}
//			bot.Send(msg)
//		}
//	}
//}
