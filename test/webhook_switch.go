package main

import (
	"github.com/flipped-aurora/gin-vue-admin/server/task"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"strings"
)

type WebhookRequest struct {
	UpdateID int              `json:"update_id"`
	Message  tgbotapi.Message `json:"message"`
}
type JenkinsBuild struct {
	ViewName string `json:"view_name"`
	JobName  string `json:"job_name"`
}

var (
	BotToken = "7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4"
	hookUrl  = "https://8ca3-91-75-118-214.ngrok-free.app/jenkins/telegram-webhook"
)

func (j *JenkinsBuild) SendMessageCommand(message tgbotapi.Message, bot *tgbotapi.BotAPI) {
	reply := tgbotapi.NewMessage(message.Chat.ID, "构建任务已触发:  "+j.ViewName+j.JobName+"正在构建中...请稍等")
	reply.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(reply); err != nil {
		return
	}
	return
}
func MessageCommandHelp(message tgbotapi.Message, bot *tgbotapi.BotAPI) {
	args := message.CommandArguments()
	reply := tgbotapi.NewMessage(message.Chat.ID, "触发构建用例: /jenkins 0898国际 后台API @CG33333_bot"+"\n参数"+args)
	reply.ReplyToMessageID = message.MessageID
	if _, err := bot.Send(reply); err != nil {
		return
	}
}

func (j *JenkinsBuild) tgWebhook(bot *tgbotapi.BotAPI) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request WebhookRequest
		// 解析请求体
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		// 处理请求   基于命令处理请求 不支持中文指令
		message := request.Message
		if !message.IsCommand() {
			c.JSON(http.StatusOK, gin.H{"status": "not a command"})
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
				j.ViewName = argsNum[0]
				j.JobName = argsNum[1]
				task.JenkinsBuildJobWithView(j.ViewName, j.JobName)
				j.SendMessageCommand(message, bot)
			}
		case "help":
			MessageCommandHelp(message, bot)
		case "构建":
			j.SendMessageCommand(message, bot)
		case "帮助":
			MessageCommandHelp(message, bot)
		default:
			log.Printf("命令:" + command)
			reply := tgbotapi.NewMessage(message.Chat.ID, "默认"+
				"\n用例: /build 0898国际 后台API @CG33333_bot")
			reply.ReplyToMessageID = message.MessageID
			if _, err := bot.Send(reply); err != nil {
				return
			}
		}

		// 返回响应
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func main() {
	// 初始化 bot 实例
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatal(err)
	}
	// 开启debug 模式
	bot.Debug = true
	// 创建一个webhook 网络钩子
	webhook, err := tgbotapi.NewWebhook(hookUrl)
	// 对webhook实例发起请求
	_, err = bot.Request(webhook)
	// 获取webhookInfo 返回的请求状态
	info, err := bot.GetWebhookInfo()
	// 如果请求失败，打印错误信息
	if info.LastErrorDate != 0 {
		log.Printf("Telegram webhook 初始化报错: %s\n", info.LastErrorMessage)

	}
	// 创建一个 Gin 路由器
	r := gin.Default()
	jb := &JenkinsBuild{}
	// 设置 Webhook 回调路由
	r.POST("/jenkins/telegram-webhook", jb.tgWebhook(bot))

	// 启动 Gin 服务器，并监听在 127.0.0.1:5299 端口
	if err := r.Run("127.0.0.1:8888"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// 用法 /build AK国际 后台API @机器人
