package main

import (
	"bytes"
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/task"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
	"strings"
)

type WebhookRequest struct {
	UpdateID int              `json:"update_id"`
	Message  tgbotapi.Message `json:"message"`
}

var (
	ChatID         int64 = -4275796428
	BotToken             = "7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4"
	telegramAPIURL       = "https://api.telegram.org/bot" + BotToken + "/sendMessage"
	bot, _               = tgbotapi.NewBotAPI(BotToken)
	hookUrl              = "https://e692-91-75-118-214.ngrok-free.app/telegram-webhook"
)

func MessageCommandStart(message tgbotapi.Message, bot *tgbotapi.BotAPI) {
	reply := tgbotapi.NewMessage(message.Chat.ID, "构建任务已触发:正在构建中，请稍等")
	reply.ReplyToMessageID = message.MessageID
	bot.Send(reply)
	return
}
func MessageCommandHelp(message tgbotapi.Message, bot *tgbotapi.BotAPI) {
	args := message.CommandArguments()
	reply := tgbotapi.NewMessage(message.Chat.ID, "帮助，来触发构建用例: /jenkins 0898国际 后台API @CG33333_bot"+"\n参数"+args)
	reply.ReplyToMessageID = message.MessageID
	bot.Send(reply)
}

func tgWebhook(bot *tgbotapi.BotAPI) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request WebhookRequest

		// 解析请求体
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}
		//var command string
		//message := request.Message
		//if strings.Contains(message.Text, "@") {
		//	index := strings.Index(message.Text, "@")
		//	if index != -1 && index < len(message.Text)-1 {
		//		command = strings.Split(message.Text, "@")[0]
		//	}
		//} else {
		//	command = message.Text
		//}
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
				log.Printf("视图名称%v", argsNum[0])
				log.Printf("JobName名称%v", argsNum[1])
				viewName := argsNum[0]
				JobName := argsNum[1]
				task.JenkinsBuildJobWithView(viewName, JobName)
			}
			MessageCommandStart(message, bot)
		case "help":
			MessageCommandHelp(message, bot)
		case "构建":
			MessageCommandStart(message, bot)
		case "帮助":
			MessageCommandHelp(message, bot)
		default:
			log.Printf("命令:" + command)
			reply := tgbotapi.NewMessage(message.Chat.ID, "默认"+
				"\n用例: /build 0898国际 后台API @CG33333_bot")
			reply.ReplyToMessageID = message.MessageID
			bot.Send(reply)
		}

		// 返回响应
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func webhookHandler(c *gin.Context) {
	var update system.Update
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading body"})
		return
	}

	if err := json.Unmarshal(body, &update); err != nil {
		log.Println("Error unmarshalling JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error unmarshalling JSON"})
		return
	}

	log.Printf("接收消息  from %s: %s", update.Message.From.Username, update.Message.Text)

	// 处理接收到的消息
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	sendMessage(ChatID, update.Message.Text)
}

//

// 发送消息到 Telegram
func sendMessage(chatID int64, text string) {
	message := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	messageBody, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}

	resp, err := http.Post(telegramAPIURL, "application/json", bytes.NewBuffer(messageBody))
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-OK HTTP status: %s", resp.Status)
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Response body: %s", body)
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

	// 设置 Webhook 回调路由
	r.POST("/telegram-webhook", tgWebhook(bot))

	// 启动 Gin 服务器，并监听在 127.0.0.1:5299 端口
	if err := r.Run("127.0.0.1:5299"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
