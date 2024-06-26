package system

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/task"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"strings"
)

var jkBuild system.JenkinsBuild

func SendMessage(bot *tgbotapi.BotAPI, message tgbotapi.Message) {
	reply := tgbotapi.NewMessage(message.Chat.ID, "构建任务已触发:  "+jkBuild.ViewName+" "+jkBuild.JobName+" 正在构建中...请稍等")
	reply.ReplyToMessageID = message.MessageID
	_, _ = bot.Send(reply)
	return
}
func BuildJobsWithText(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	msg := webhook.Message.Text
	botUsername := bot.Self.UserName
	if strings.Contains(msg, "@"+botUsername) {
		params := strings.Fields(msg)
		if len(params) == 4 {
			log.Printf("参数%s", params)
			jkBuild.ViewName = params[1]
			jkBuild.JobName = params[2]
			task.JenkinsBuildJobWithView(jkBuild.ViewName, jkBuild.JobName)
			replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, "构建任务已触发:  "+jkBuild.ViewName+" "+jkBuild.JobName+" 正在构建中...请稍等")
			replyText.ReplyToMessageID = webhook.Message.MessageID
			_, _ = bot.Send(replyText)
		} else {
			replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, "详细用法"+msg+"任务名:"+jkBuild.JobName)
			replyText.ReplyToMessageID = webhook.Message.MessageID
			_, _ = bot.Send(replyText)
		}
	}
	return
}

//	func handleUpdate(bot *tgbotapi.BotAPI, update *system.WebhookRequest) {
//		if update.Message != nil {
//			if update.Message.IsCommand() {
//				handleCommand(bot, update.Message)
//			} else if update.Message.Text != "" {
//				handleText(bot, update.Message)
//			}
//		}
//	}
func BuildJobsCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	command := message.Command()
	args := message.CommandArguments() // 获取全部参数
	argsNum := strings.Fields(args)    // 已空格为定界符 将参数进行分割
	switch command {
	case "build":
		if len(argsNum) == 3 {
			log.Printf("视图名称: %v ", argsNum[0])
			log.Printf("JobName名称: %v", argsNum[1])
			jkBuild.ViewName = argsNum[0]
			jkBuild.JobName = argsNum[1]
			task.JenkinsBuildJobWithView(jkBuild.ViewName, jkBuild.JobName)
			reply := tgbotapi.NewMessage(message.Chat.ID, "构建任务已触发:  "+jkBuild.ViewName+" "+jkBuild.JobName+" 正在构建中...请稍等")
			reply.ReplyToMessageID = message.MessageID
			_, _ = bot.Send(reply)
		}
	case "help":
		reply := tgbotapi.NewMessage(message.Chat.ID, "Jenkins构建用例:"+"  "+"字母大写\n"+
			"/build 0898国际 后台API @CG33333_bot\n"+
			"/build 0898国际 前台API @CG33333_bot\n"+
			"/build 0898国际 H5 @CG33333_bot\n"+
			"/build 0898国际 后台H5 @CG33333_bot\n"+
			"/build 0898国际 定时任务 @CG33333_bot")
		reply.ReplyToMessageID = message.MessageID
		_, _ = bot.Send(reply)
	default:
		reply := tgbotapi.NewMessage(message.Chat.ID, "Jenkins构建用例:"+"  "+"字母大写\n"+
			"/build 0898国际 后台API @CG33333_bot\n"+
			"/build 0898国际 前台API @CG33333_bot\n"+
			"/build 0898国际 H5 @CG33333_bot\n"+
			"/build 0898国际 后台H5 @CG33333_bot\n"+
			"/build 0898国际 定时任务 @CG33333_bot")
		reply.ReplyToMessageID = message.MessageID
		_, _ = bot.Send(reply)

	}
}

func (b *BaseApi) TelegramWebhook(c *gin.Context) {
	// 初始化 bot 实例
	bot, err := tgbotapi.NewBotAPI(global.GVA_CONFIG.Telegram.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	// 开启debug 模式
	bot.Debug = true
	// 创建一个webhook 网络钩子
	webhook, err := tgbotapi.NewWebhook(global.GVA_CONFIG.Telegram.WebhookUrl)
	// 对webhook实例发起请求
	_, err = bot.Request(webhook)
	// 获取webhookInfo 返回的请求状态
	info, err := bot.GetWebhookInfo()
	// 如果请求失败，打印错误信息
	if info.LastErrorDate != 0 {
		log.Printf("webhook 最新的报错消息: %s\n", info.LastErrorMessage)
	}
	log.Printf("IP地址: %s\n", info.IPAddress)
	fmt.Printf("WenHook地址： %v\n", info.URL)
	fmt.Printf("最新的错误信息： %v\n", info.LastErrorMessage)
	fmt.Printf("最大连接数：%v\n", info.MaxConnections)
	fmt.Println("-----------------------------------------------------")
	log.Printf("验证成功\n")
	fmt.Println("-----------------------------------------------------")
	log.Printf("User ID:%v\n", bot.Self.UserName)
	fmt.Println("-----------------------------------------------------")
	log.Printf("PEER ID:%v\n", bot.Self.ID)
	fmt.Println("-----------------------------------------------------")
	log.Printf("bot Name:%s %s\n", bot.Self.FirstName, bot.Self.LastName)
	fmt.Println("-----------------------------------------------------")
	var update system.WebhookRequest // telegram消息响应结构体
	//var message *tgbotapi.Message    // telegram update结构体
	//message := request.Message
	// jenkins构建结构体()
	// 解析请求体
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	BuildJobsWithText(bot, update)
	//BuildJobsCommand(bot, message)
}
