package system

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/task"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/selenium"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var jkBuild system.JenkinsBuild

// BuildJobsWithText 基于文本消息构建jenkins 任务
func BuildJobsWithText(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	msg := webhook.Message.Text
	botUsername := bot.Self.UserName
	//username := webhook.Message.From.UserName
	if strings.Contains(msg, "@"+botUsername) {
		// 获取外部参数
		params := strings.Fields(msg)
		log.Printf("参数: %v", params)
		if len(params) > 3 {
			log.Printf("参数%s", params)
			jkBuild.ViewName = params[1]
			jkBuild.JobName = params[2]
			if params[3] == "更新" {
				if webhook.Message.From.UserName == "David5886" || webhook.Message.From.UserName == "@zero666_777" || webhook.Message.From.UserName == "sweet" {
					//log.Printf("触发的用户:%s", username)
					task.JenkinsBuildJob(jkBuild.ViewName, jkBuild.JobName) // 包网
					task.JenkinsBuildJobWithView(jkBuild.ViewName, jkBuild.JobName)
					replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, "构建任务已触发:  "+jkBuild.ViewName+" "+jkBuild.JobName+" 正在构建中...30秒后获取构建状态")
					replyText.ReplyToMessageID = webhook.Message.MessageID
					_, _ = bot.Send(replyText)
					go func() { // goroutine 异步方式  不影响主逻辑 延迟操作可以在后台运行
						time.Sleep(30 * time.Second)
						//获取构建任务的状态
						status, _ := task.GetLastBuildStatus(jkBuild.ViewName, jkBuild.JobName)
						formattedTime := time.Unix(status.Timestamp/1000, (status.Timestamp%1000)*1000000).UTC().In(func() *time.Location { loc, _ := time.LoadLocation("Asia/Dubai"); return loc }()).Format("2006-01-02 15:04:05")
						result := tgbotapi.NewMessage(webhook.Message.Chat.ID, "执行结果: "+status.Result+"\n最近的构建号: "+strconv.Itoa(status.Number)+"\n构建时间: "+formattedTime)
						result.ReplyToMessageID = webhook.Message.MessageID
						_, _ = bot.Send(result)
					}()
				}
			}
		} else {
			// 参数不足 打印帮助信息
			printHelp(bot, webhook)
		}
	}
	return
}

// GetProjectParams 获取构建任务 git 分支
func GetProjectParams(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	msg := webhook.Message.Text
	botUsername := bot.Self.UserName
	if strings.Contains(msg, "@"+botUsername) {
		// 获取外部参数
		params := strings.Fields(msg)
		if len(params) > 3 {
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

// GetAdminLink 获取管理后台地址
func GetAdminLink(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	msg := webhook.Message.Text
	botUsername := bot.Self.UserName
	if strings.Contains(msg, "@"+botUsername) {
		// 获取外部参数
		params := strings.Fields(msg)
		if len(params) > 3 {
			siteName := params[1]
			envName := params[2]
			log.Printf("参数一：%s", siteName)
			log.Printf("参数2：%s", envName)
			if params[3] == "后台链接" {
				go func() {
					defer func() {
						if err := recover(); err != nil {
							code, _ := selenium.GetCaptchaCode()
							message := siteName + "解析验证码错误:   " + code + "\n有效验证码为4位数字，请@机器人重试！"
							replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, message)
							replyText.ReplyToMessageID = webhook.Message.MessageID
							_, _ = bot.Send(replyText)
						}
					}()
					siteLink := selenium.GetAdminLinkTools(siteName)
					message := siteName + "站点地址: " + siteLink
					replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, message)
					replyText.ReplyToMessageID = webhook.Message.MessageID
					_, _ = bot.Send(replyText)
				}()
			}
			if params[3] == "后台地址" {
				go func() {
					siteLink := selenium.GetLinkNoLogin(siteName)
					message := siteName + "站点地址: " + siteLink
					replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, message)
					replyText.ReplyToMessageID = webhook.Message.MessageID
					_, _ = bot.Send(replyText)
					defer func() {
						if err := recover(); err != nil {
							message := siteName + "异常"
							replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, message)
							replyText.ReplyToMessageID = webhook.Message.MessageID
							_, _ = bot.Send(replyText)
						}
					}()
				}()
			}
		}
	}
}

func GetGameFile(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {

}

// 测试环境 重启服务
func restartService(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	msg := webhook.Message.Text
	botUsername := bot.Self.UserName
	if strings.Contains(msg, "@"+botUsername) {
		params := strings.Fields(msg)
		if len(params) > 3 {
			// 测试环境
			siteName := params[1]
			EnvName := params[2]
			command := params[3]
			switch EnvName {
			case "测试":
				replyWithMessage(bot, webhook, "这是测试环境")
				task.ManageService(siteName, EnvName, command)
			}

		}
	}
}

// printHelp 打印帮助信息
func printHelp(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	helpMessage := "使用说明:\n" +
		"@机器人 <视图名> <任务名> 更新 - 指定视图名和任务名触发构建\n" +
		"示例: @CG33333_bot 28国际 后台API 更新\n" +
		"@CG33333_bot 28国际 后台API 获取分支\n" +
		"@CG33333_bot 28国际 生产 后台地址\n"
	replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, helpMessage)
	replyText.ReplyToMessageID = webhook.Message.MessageID
	_, _ = bot.Send(replyText)
}

func replyWithMessage(bot *tgbotapi.BotAPI, webhook system.WebhookRequest, message string) {
	replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, message)
	replyText.ReplyToMessageID = webhook.Message.MessageID
	_, _ = bot.Send(replyText)
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
	log.Printf("机器人名:%v\n", bot.Self.UserName)
	fmt.Println("-----------------------------------------------------")
	var update system.WebhookRequest // telegram消息响应结构体
	// 解析请求体
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	GetAdminLink(bot, update)
	BuildJobsWithText(bot, update)
	GetProjectParams(bot, update)
	restartService(bot, update)
}
