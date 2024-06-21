package system

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/task"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"strings"
)

var jkBuild system.JenkinsBuild

func SendMessageCommand(message tgbotapi.Message, bot *tgbotapi.BotAPI) {
	reply := tgbotapi.NewMessage(message.Chat.ID, "构建任务已触发:  "+jkBuild.ViewName+" "+jkBuild.JobName+" 正在构建中...请稍等")
	reply.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(reply); err != nil {
		return
	}
	return
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
	var request system.WebhookRequest // telegram消息响应结构体
	// jenkins构建结构体()
	// 解析请求体
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	// 处理请求   基于命令处理请求 不支持中文指令
	message := request.Message
	if !message.IsCommand() {
		response.OkWithDetailed("success", "不是一个命令", c)
		return
	}
	command := message.Command()       // 获取命令
	args := message.CommandArguments() // 获取全部参数
	argsNum := strings.Fields(args)    // 已空格为定界符 将参数进行分割
	// 检查命令是否包含bot的用户名（在群组中会有这种情况）
	if strings.Contains(command, "@") {
		command = strings.Split(command, "@")[0]
	}
	switch command {
	case "build":
		log.Printf("接收到的命令: %s", command)
		log.Printf("接收到的全部参数: %s", args)
		if len(argsNum) == 3 {
			log.Printf("视图名称: %v ", argsNum[0])
			log.Printf("JobName名称: %v", argsNum[1])
			jkBuild.ViewName = argsNum[0]
			jkBuild.JobName = argsNum[1]
			task.JenkinsBuildJobWithView(jkBuild.ViewName, jkBuild.JobName)
			SendMessageCommand(message, bot)
		}
	case "help":
		reply := tgbotapi.NewMessage(message.Chat.ID, "Jenkins构建用例:"+"  "+"字母大写\n"+
			"/build 0898国际 后台API @CG33333_bot\n"+
			"/build 0898国际 前台API @CG33333_bot\n"+
			"/build 0898国际 H5 @CG33333_bot\n"+
			"/build 0898国际 后台H5 @CG33333_bot\n"+
			"/build 0898国际 定时任务 @CG33333_bot")
		reply.ReplyToMessageID = message.MessageID
		if _, err := bot.Send(reply); err != nil {
			return
		}
	default:
		reply := tgbotapi.NewMessage(message.Chat.ID, "Jenkins构建用例:"+"  "+"字母大写\n"+
			"/build 0898国际 后台API @CG33333_bot\n"+
			"/build 0898国际 前台API @CG33333_bot\n"+
			"/build 0898国际 H5 @CG33333_bot\n"+
			"/build 0898国际 后台H5 @CG33333_bot\n"+
			"/build 0898国际 定时任务 @CG33333_bot")
		reply.ReplyToMessageID = message.MessageID
		if _, err := bot.Send(reply); err != nil {
			return
		}
	}

}
