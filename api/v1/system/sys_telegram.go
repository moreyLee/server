package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
)

func (b *BaseApi) TelegramWebhook(c *gin.Context) {
	var update tgbotapi.BotAPI
	err := c.ShouldBindJSON(&update)
	// 处理Telegram Webhook请求
	message := common.ReceiveMessage{}
	requestBody := new(bytes.Buffer)
	chatID := 0
	msgText := ""

	err = json.NewDecoder(requestBody).Decode(&message)
	if err != nil {
		fmt.Println(err)
	}

	// if private or group
	if message.Message.Chat.ID != 0 {
		fmt.Println(message.Message.Chat.ID, message.Message.Text)
		chatID = message.Message.Chat.ID
		msgText = message.Message.Text
	} else {
		// if channel
		fmt.Println(message.ChannelPost.Chat.ID, message.ChannelPost.Text)
		chatID = message.ChannelPost.Chat.ID
		msgText = message.ChannelPost.Text
	}

	respMsg := fmt.Sprintf("%s%s/sendMessage?chat_id=%d&text=Received: %s",
		global.GVA_CONFIG.Telegram.URL,
		global.GVA_CONFIG.Telegram.BotToken,
		chatID,
		msgText)

	// send echo resp
	_, err = http.Get(respMsg)
	if err != nil {
		fmt.Println(err)
	}
}
