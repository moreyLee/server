package system

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
)

var (
	URL  = "https://api.telegram.org/bot"
	port = "5299"
	//webhookURL = "https://devops.3333d.vip/telegram-webhook"
	webhookUrl = "https://cf4b-91-75-118-214.ngrok-free.app/telegram-webhook"
	token      = "7449933946:AAGSpUHIsi9cTgc65O9CFheOia3czrLS8l4"
)

func (b *BaseApi) TelegramWebhook(c *gin.Context) {
	// 类似与main 函数中的初始化
	// 初始化 bot 实例
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	// 开启debug 模式
	bot.Debug = true
	// 创建一个webhook 网络钩子
	webhook, err := tgbotapi.NewWebhook(webhookUrl)
	// 对webhook实例发起请求
	_, err = bot.Request(webhook)
	// 获取webhookInfo 返回的请求状态
	info, err := bot.GetWebhookInfo()
	// 如果请求失败，打印错误信息
	if info.LastErrorDate != 0 {
		log.Printf("webhook 最新的报错消息: %s\n", info.LastErrorMessage)
	}
	fmt.Printf("IP地址: %v\n", info.IPAddress)
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

	// 遍历收到的消息
	updates := bot.ListenForWebhook("/")

	// 调用webhook 函数 通过post方法
	message := system.ReceiveMessage{}
	json.NewDecoder(c.Request.Body).Decode(&message)
	chatID := 0
	msgText := ""
	if message.Message.Chat.ID != 0 {
		fmt.Println("群组ID与消息文本", message.Message.Chat.ID, message.Message.Text)
		// 获取到消息的chatID 消息内容
		chatID = message.Message.Chat.ID
		msgText = message.Message.Text
	}
	respMsg := fmt.Sprintf("%s%s/sendMessage?chat_id=%d&text=Received: %s",
		URL,
		token,
		chatID,
		msgText)
	http.Get(respMsg)
	fmt.Println("发送信息的url", respMsg)
	//go func() {
	//	err := http.ListenAndServe(":5299", nil) // 原生的http包
	//	if err != nil {
	//		return
	//	}
	//
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
	//}()
}
