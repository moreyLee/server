package task

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
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

func BotJob() {
	// 启动协程运行函数
	go func() {
		defer func() {
			// recover 函数只能在defer()函数中调用 用于恢复程序控制流
			if err := recover(); err != nil {
				log.Printf("telegram 机器人运行出错,参数错误:\n%v\n", err)
			}
		}()
		// 初始化机器人
		//bot := initialize.SingleBot()
		//bot := GetBotInstance("7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4")
		bot := GetBotInstance(global.GVA_CONFIG.Telegram.BotToken)

		//bot, err := tgbotapi.NewBotAPI(global.GVA_CONFIG.Telegram.BotToken)
		//if err != nil {
		//	log.Panic(err)
		//}
		// 启用调试模式 慢sql 语句优化
		bot.Debug = false
		log.Printf("机器人名称: @%s", bot.Self.UserName)

		// 创建一个新的消息
		//chatID := global.GVA_CONFIG.Telegram.ChatID // 替换为目标聊天 ID（负数表示群组）
		//messageText := "欢迎使用CG机器人"
		//发送消息
		//msg := tgbotapi.NewMessage(chatID, messageText)
		//发送消息
		//bot.Send(msg)

		// 设置更新配置
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 120 // 设置超时时间
		// 触发panic 异常
		//actualFunction()
		// 发送启动消息

		// 获取更新通道
		updates := bot.GetUpdatesChan(u)

		// 启动消息发送协程
		//SendMessage(messages)
		for update := range updates {
			if update.Message == nil { // 忽略任何非消息更新
				continue
			}
			//time.Sleep(3 * time.Second)
			// 打印收到的消息
			//chatID := update.Message.Chat.ID
			msgText := update.Message.Text
			// 检查消息是否提到了机器人
			if update.Message.IsCommand() || strings.Contains(msgText, "@"+bot.Self.UserName) {
				// 回复消息
				//responseText := "你提到我了吗？我在这里！"
				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
				//msg.ReplyToMessageID = update.Message.MessageID
				// 发送回复消息
				//bot.Send(msg)
				args := update.Message.CommandArguments()
				// 使用strings.fields 来分割参数
				argsNum := strings.Fields(args)
				//log.Printf("args参数个数：%v", len(argsNum))
				ViewParam := update.Message.CommandArguments()
				// 检查命令
				switch update.Message.Command() {
				case "jenkins":
					// 从机器人输入四个参数
					if len(argsNum) == 3 {
						//log.Printf("入参:  " + args)
						field := strings.Fields(args)
						jobName := field[1]
						log.Printf("JobName名称: %s", jobName)
						viewName := strings.SplitN(ViewParam, " ", 2)[0]
						log.Printf("Jenkins视图名称: %s", viewName)
						JenkinsBuildJobWithView(viewName, jobName)
						//reply := tgbotapi.NewMessage(update.Message.Chat.ID, "已收到请求，"+viewName+"_"+jobName+"正在构建中，请稍等")
						//reply.ReplyToMessageID = update.Message.MessageID
						//bot.Send(reply)
					} else {
						log.Printf("请输入正确的参数,参考/help")
						//reply := tgbotapi.NewMessage(update.Message.Chat.ID, "错误: 请输入正确的参数个数,请参考 /help @CG33333_bot")
						//reply.ReplyToMessageID = update.Message.MessageID
						//bot.Send(reply)
					}
				case "help":
					//reply := tgbotapi.NewMessage(update.Message.Chat.ID, "请使用 /jenkins 项目名 [后台API/前台API/H5/后台H5/定时任务]   任选其一，触发构建"+
					//	"\n用例: /jenkins 0898国际 后台API @CG33333_bot")
					//reply.ReplyToMessageID = update.Message.MessageID
					//bot.Send(reply)
				default:
					//reply := tgbotapi.NewMessage(update.Message.Chat.ID, "请使用 /jenkins 项目名 [后台API/前台API/H5/后台H5/定时任务] 任选其一，来触发构建"+
					//	"\n用例: /jenkins 0898国际 后台API @CG33333_bot")
					//reply.ReplyToMessageID = update.Message.MessageID
					//bot.Send(reply)
				}
			}
		}
		log.Printf("异常以后继续执行后面的业务逻辑")
	}()
}

// 测试中断函数 用于触发panic 异常退出
//func actualFunction() {
//	global.GVA_LOG.Info("开始执行触发异常函数")
//	// 触发panic
//	panic("触发panic错误")
//}
