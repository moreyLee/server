package task

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func SendMessage() {
	// 启动协程运行函数
	go func() {
		defer func() {
			// recover 函数只能在defer()函数中调用 用于恢复程序控制流
			if err := recover(); err != nil {
				log.Printf("telegram 机器人运行出错捕获异常信息:\n%v\n", err)
			}
		}()
		botToken := "7438242996:AAFwGnP8mQBmvcjiDggltiiOTMo14XeOoT4"
		// 初始化机器人
		bot, err := tgbotapi.NewBotAPI(botToken)
		if err != nil {
			log.Panic(err)
		}
		// 启用调试模式 慢sql 语句优化
		bot.Debug = false

		log.Printf("机器人名称: @%s", bot.Self.UserName)

		// 创建一个新的消息
		chatID := int64(-4275796428) // 替换为目标聊天 ID（负数表示群组）
		messageText := "开始"
		// 发送消息
		msg := tgbotapi.NewMessage(chatID, messageText)
		//
		// 发送消息
		_, err = bot.Send(msg)
		if err != nil {
			log.Panic(err)
		}

		// 设置更新配置
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		// 触发panic 异常
		//actualFunction()
		// 获取更新通道
		updates := bot.GetUpdatesChan(u)
		for update := range updates {
			if update.Message == nil { // 忽略任何非消息更新
				continue
			}
			// 打印收到的消息
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// 检查消息是否提到了机器人
			if strings.Contains(update.Message.Text, "@"+bot.Self.UserName) {

				// 检查命令
				switch update.Message.Command() {
				case "jenkins":
					err := JenkinsBuildJob()
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "JobName")
					msg.ReplyToMessageID = update.Message.MessageID
					// 发送回复消息
					_, err = bot.Send(msg)
					if err != nil {
						log.Panic(err)
					}
					//
				}
				// 回复消息
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID
				// 发送回复消息
				_, err = bot.Send(msg)
				if err != nil {
					log.Panic(err)
				}
			}
		}

		log.Printf("异常以后继续执行后面的业务逻辑")
	}()
}

// 测试中断函数 用于触发panic 异常退出
func actualFunction() {
	global.GVA_LOG.Info("开始执行触发异常函数")
	// 触发panic
	panic("触发panic错误")
}
