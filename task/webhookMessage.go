package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
)

var bot *tgbotapi.BotAPI

func WebhookMessage(r *gin.Context) {
	initialize.InitBot()
	message := common.ReceiveMessage{}
	Body := new(bytes.Buffer)
	json.NewDecoder(Body).Decode(&message)
	updates := bot.ListenForWebhook("/")
	respMsg := fmt.Sprintf("%s%s/sendMessage?chat_id=%d&text=Received: %s",
		global.GVA_CONFIG.Telegram.URL,
		global.GVA_CONFIG.Telegram.BotToken,
		global.GVA_CONFIG.Telegram.ChatID,
		"")
	http.Get(respMsg)
	go func() {
		err := http.ListenAndServe("127.0.0.1:5299", nil)
		if err != nil {
			return
		}

		for update := range updates {
			if update.Message == nil {
				continue
			}
			switch update.Message.Text {
			case "jenkins":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "1")
				bot.Send(msg)
			default:
				fmt.Println("默认输出")
			}
		}
	}()
}
