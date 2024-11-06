package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// 替换为你的 Bot Token 和创建者 ID
const (
	Bot_token = "7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4"
	CreatorID = 6479992479 // 替换成你的 User ID
)

func main() {
	// 初始化 bot
	bot, err := tgbotapi.NewBotAPI(Bot_token)
	if err != nil {
		log.Fatalf("无法连接到 Telegram Bot API: %v", err)
	}

	// 打印当前 bot 信息
	bot.Debug = true
	log.Printf("已授权的 Bot: %s", bot.Self.UserName)

	// 创建一个更新通道
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// 监听更新
	for update := range updates {

		if update.Message != nil {
			// 检查是否是机器人被添加到群组的消息
			if update.Message.NewChatMembers != nil && update.Message.Chat.IsGroup() {
				for _, user := range update.Message.NewChatMembers {
					// 检查新成员是否是机器人本身
					if user.ID == bot.Self.ID {
						// 获取发起添加请求的用户 ID
						inviterID := update.Message.From.ID

						// 检查是否为创建者
						if inviterID != CreatorID {
							// 非创建者，离开群组
							leave := tgbotapi.LeaveChatConfig{
								ChatID: update.Message.Chat.ID,
							}
							_, err := bot.Request(leave)
							if err != nil {
								log.Printf("无法离开群组: %v", err)
							} else {
								fmt.Println("非创建者添加，已离开群组。")
							}
						}
					}
				}
			}
		}
	}
}
