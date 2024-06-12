package main

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
)

// 接收消息
var (
	//webhookURL = "https://devops.3333d.vip/telegram-webhook"
	webhookUrl = "https://7018-91-75-118-214.ngrok-free.app/telegram-webhook"
	token      = "7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4"
)

const chatID int64 = -4275796428

func handleTelegramWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	//body, err := io.ReadAll(r.Body)
	//if err != nil {
	//	log.Printf("Error reading request body: %v", err)
	//	http.Error(w, "Error reading request body", http.StatusBadRequest)
	//	return
	//}

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("Error decoding update: %v", err)
		//log.Printf("Request body: %s\n", string(body))
		http.Error(w, "Error decoding update", http.StatusBadRequest)
		return
	}

	//log.Printf("Webhook update: %+v", update)
	// 获取update_id
	//updateID := update.UpdateID
	//log.Printf("Update ID: %d", updateID)

	// 处理消息
	if update.Message != nil && update.Message.Chat != nil && update.Message.Chat.ID == -4275796428 {
		log.Printf("Webhook update: %+v", update)
		// Process the message here
	}
	//if update.Message != nil {
	//	chatID := update.Message.Chat.ID
	//	text := update.Message.Text
	//
	//	responseText := "You said: " + text
	//	msg := tgbotapi.NewMessage(chatID, responseText)
	//	if _, err := bot.Send(msg); err != nil {
	//		log.Printf("Error sending message: %v", err)
	//	}
	//}

	w.WriteHeader(http.StatusOK)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	// ...
	_, err = tgbotapi.NewWebhook(webhookUrl)

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	http.HandleFunc("/telegram-webhook", handleTelegramWebhook)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
