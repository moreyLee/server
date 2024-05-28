package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type telegramMessage struct {
	ChatID                int64  `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

func SendMessageToTelegramGroup(token string, chatID int64, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	msg := telegramMessage{
		ChatID:                chatID,
		Text:                  message,
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
	}
	messageData, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(messageData))
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
