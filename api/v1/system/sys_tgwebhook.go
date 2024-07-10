package system

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/task"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"strings"
)

var jkBuild system.JenkinsBuild

func BuildJobsWithText(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	msg := webhook.Message.Text
	botUsername := bot.Self.UserName
	if strings.Contains(msg, "@"+botUsername) {
		// 获取外部参数
		params := strings.Fields(msg)
		log.Printf("参数: %v", params)
		if len(params) > 3 {
			log.Printf("参数%s", params)
			jkBuild.ViewName = params[1]
			jkBuild.JobName = params[2]
			if params[3] == "更新" {
				task.JenkinsBuildJobWithView(jkBuild.ViewName, jkBuild.JobName)
				replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, "构建任务已触发:  "+jkBuild.ViewName+" "+jkBuild.JobName+" 正在构建中...请稍等")
				replyText.ReplyToMessageID = webhook.Message.MessageID
				_, _ = bot.Send(replyText)
			}
		}
	} else {
		log.Printf("未找到@%s", botUsername)
		//text := tgbotapi.NewMessage(webhook.Message.Chat.ID, "请提供足够的参数触发构建任务 \n"+"构建用例:\n"+
		//	"@"+botUsername+" 0898国际 后台API  \n"+
		//	"@"+botUsername+" 0898国际 前台API \n"+
		//	"@"+botUsername+" 0898国际 H5 \n"+
		//	"@"+botUsername+" 0898国际 后台H5 \n"+
		//	"@"+botUsername+" 0898国际 定时任务 ")
		//text.ReplyToMessageID = webhook.Message.MessageID
		//_, _ = bot.Send(text)
	}
	return
}

func GetProjectParams(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	msg := webhook.Message.Text
	botUsername := bot.Self.UserName
	if strings.Contains(msg, "@"+botUsername) {
		// 获取外部参数
		params := strings.Fields(msg)
		log.Printf("参数: %v", params)
		if len(params) > 3 {
			log.Printf("参数%s", params)
			jkBuild.ViewName = params[1]
			jkBuild.JobName = params[2]
			if params[3] == "获取分支" {
				jobConfig := task.GetBranch(jkBuild.ViewName, jkBuild.JobName)
				if len(jobConfig.SCM.UserRemoteConfigs.URLs) > 0 && len(jobConfig.SCM.Branches) > 0 {
					replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, jkBuild.ViewName+"  "+jkBuild.JobName+"的git分支:  "+jobConfig.SCM.Branches[0])
					replyText.ReplyToMessageID = webhook.Message.MessageID
					_, _ = bot.Send(replyText)
				}
			}
		}
	}

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
	log.Printf("userID:%v\n", bot.Self.UserName)
	fmt.Println("-----------------------------------------------------")
	var update system.WebhookRequest // telegram消息响应结构体
	// 解析请求体
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	BuildJobsWithText(bot, update)
	GetProjectParams(bot, update)
}
