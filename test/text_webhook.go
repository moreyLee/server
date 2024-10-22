package main

import (
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
)

type TextMessage struct {
	UpdateID int              `json:"update_id"`
	Message  tgbotapi.Message `json:"message"`
}
type Build struct {
	ViewName string `json:"view_name"`
	JobName  string `json:"job_name"`
}
type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

const (
	Token      = "7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4"                     // Replace with your Telegram Bot Token
	WebhookURL = "https://7988-87-200-210-97.ngrok-free.app/jenkins/telegram-webhook" // Replace with your webhook URL
)

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message != nil {
		if update.Message.IsCommand() {
			handleCommand(bot, update.Message)
		} else if update.Message.Text != "" {
			handleText(bot, update.Message)
		}
	}
}

func handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	command := message.Command()
	switch command {
	case "start":
		reply := tgbotapi.NewMessage(message.Chat.ID, "Welcome! How can I assist you today?")
		bot.Send(reply)
	case "help":
		reply := tgbotapi.NewMessage(message.Chat.ID, "Here are the available commands:\n/start - Start the bot\n/help - Show this help message")
		bot.Send(reply)
	default:
		reply := tgbotapi.NewMessage(message.Chat.ID, "Unknown command. Type /help for a list of commands.")
		bot.Send(reply)
	}
}

func handleText(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	text := message.Text
	responseText := "You said: " + text
	reply := tgbotapi.NewMessage(message.Chat.ID, responseText)
	if _, err := bot.Send(reply); err != nil {
		log.Println("Failed to send reply:", err)
	}
}

func main() {

	// Initialize the Telegram bot instance
	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		log.Fatal(err)
	}

	// Enable debug mode
	bot.Debug = true

	// Create a Gin router
	r := gin.Default()

	// Handle the Telegram webhook POST request
	r.POST("/jenkins/telegram-webhook", func(c *gin.Context) {
		var update tgbotapi.Update
		if err := c.ShouldBindJSON(&update); err != nil {
			log.Println("Failed to parse update:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request data"})
			return
		}

		// Call handleUpdate function
		handleUpdate(bot, update)

		// Return a success response
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Set up the Telegram Webhook
	webhook, err := tgbotapi.NewWebhook(WebhookURL)
	// 对webhook实例发起请求
	_, err = bot.Request(webhook)
	if err := r.Run(":" + "5888"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
