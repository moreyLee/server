package main

import (
	"bytes"
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
)

var (
	ChatID         int64 = -4275796428
	BotToken             = "7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4"
	telegramAPIURL       = "https://api.telegram.org/bot" + BotToken + "/sendMessage"
	bot, _               = tgbotapi.NewBotAPI(BotToken)
)

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
	// 检查消息是否以 "@机器人名称" 开头
	//botName := "@CG33333_bot" // 替换为你的机器人名称
	//if strings.HasPrefix(update.Message.Text, "@"+bot.Self.UserName) {
	//	// 提取命令及参数
	//	parts := strings.Fields(update.Message.Text[len(botName):])
	//	if len(parts) > 0 {
	//		command := parts[0]
	//		handleCommand(update.Message.Chat.ID, update.Message.From.Username, command)
	//	} else {
	//		log.Println("No command found after bot name.")
	//	}
	//}

	// 解析收到的消息
	// 响应消息
	responseText := "Hello, " + update.Message.From.FirstName + "! You said: " + update.Message.Text
	switch update.Message.Text {

	}
	sendMessage(update.Message.Chat.ID, responseText)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// 处理命令
func handleCommand(chatID int64, username, command string) {
	var responseText string

	// 根据命令执行相应的操作
	switch command {
	case "/start":
		responseText = "Hello, @" + username + "! Welcome! How can I assist you today?"
	case "/help":
		responseText = "Hello, @" + username + "! Available commands:\n/start - Start the bot\n/help - Show this help message"
	default:
		responseText = "Hello, @" + username + "! Unknown command. Type /help to see available commands."
	}

	// 发送回复消息到 Telegram
	sendMessage(chatID, responseText)
}

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

	// 创建一个 Gin 路由器
	r := gin.Default()

	// 设置 Webhook 回调路由
	r.POST("/telegram-webhook", webhookHandler)

	// 启动 Gin 服务器，并监听在 127.0.0.1:5299 端口
	if err := r.Run("127.0.0.1:5299"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
