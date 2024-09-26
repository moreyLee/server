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

// BuildJobsWithText 基于文本消息构建jenkins 任务   生产环境
func BuildJobsWithText(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	// 通道done 确保任务状态反馈
	done := make(chan bool)
	msg := webhook.Message.Text
	botUsername := bot.Self.UserName
	// 检查消息中是否包含机器人用户名
	if !strings.Contains(msg, "@"+botUsername) {
		return
	}
	// 获取外部参数
	params := strings.Fields(msg)
	log.Printf("参数: %v", params)
	if len(params) < 4 {
		ReplyWithMessage(bot, webhook, "输入参数不足")
		return
	}
	log.Printf("参数%s", params)
	jkBuild.ViewName = params[1]
	jkBuild.JobName = params[2]
	if params[3] != "更新" {
		return
	}
	// 启动协程执行构建任务
	go func() {
		defer func() {
			if err := recover(); err != nil {
				// 异常捕获 业务逻辑处理
				fmt.Printf("%v", err)
			}
		}()
		if isUserAuthorized(webhook.Message.From.UserName) {
			task.JenkinsBuildJobWithView(bot, webhook, jkBuild.ViewName, jkBuild.JobName, done)
		}
	}()
	// 主协程 也要等待任务完成并获取构建状态
	// select 语句监听 done 通道
	select {
	case success := <-done:
		if success {
			// 加业务逻辑容易导致死循环
			fmt.Printf("协程返回成功，进入构建状态发送\n")
			go func() {
				time.Sleep(30 * time.Second)
				sendBuildStatus(bot, webhook)
			}()
		} else {
			fmt.Printf(jkBuild.ViewName + ": 主协程通知 视图名不存在构建任务失败，无法获取构建任务状态\n") // 别删
		}
	case <-time.After(30 * time.Second):
		fmt.Printf("主协程通知 构建任务超时\n")
	}
}

// isUserAuthorized 检查用户是否被授权
func isUserAuthorized(username string) bool {
	authorizedUsers := []string{"David5886", "nikon_aaa", "tank9999999"}
	for _, user := range authorizedUsers {
		if username == user {
			return true
		}
	}
	return false
}

// sendBuildStatus 异步发送构建状态信息 避免阻塞主协程的其他操作  主协程不需要等待构建状态发送函数的执行 可以更快响应其他操作
func sendBuildStatus(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	go func() {
		if jkBuild.ViewName != "" && jkBuild.JobName != "" {
			status, err := task.GetLastBuildStatus(jkBuild.ViewName, jkBuild.JobName)
			if err != nil {
				ReplyWithMessage(bot, webhook, "获取构建任务失败: %v\n"+err.Error())
				return
			}
			formattedTime := time.Unix(status.Timestamp/1000, (status.Timestamp%1000)*1000000).UTC().In(func() *time.Location { loc, _ := time.LoadLocation("Asia/Dubai"); return loc }()).Format("2006-01-02 15:04:05")
			result := tgbotapi.NewMessage(webhook.Message.Chat.ID, "执行结果: "+status.Result+"\n最近的构建号: "+strconv.Itoa(status.Number)+"\n构建时间: "+formattedTime)
			result.ReplyToMessageID = webhook.Message.MessageID
			_, _ = bot.Send(result)
		}
	}()
}

// BuildJobsWithEnvTest  测试环境 更新Jenkins
func BuildJobsWithEnvTest(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	msg := webhook.Message.Text
	botUsername := bot.Self.UserName
	if strings.Contains(msg, "@"+botUsername) {
		// 获取外部参数
		params := strings.Fields(msg)
		if len(params) > 4 {
			jkBuild.ViewName = params[1] //  视图名
			// Env = params[2]              // 环境

		}
	}
}

// GetBranch 获取构建任务 git 分支
func GetBranch(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
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
					defer func() {
						if err := recover(); err != nil {
							fmt.Printf("GetLinkNoLogin 函数异常捕获%v", err)
						}
					}()
					siteLink := selenium.GetLinkNoLogin()
					message := siteName + "站点地址: " + siteLink
					replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, message)
					replyText.ReplyToMessageID = webhook.Message.MessageID
					_, _ = bot.Send(replyText)
				}()
			}
			if params[3] == "后台登录" {
				go func() {
					defer func() {
						if err := recover(); err != nil {
							replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, "后台登录异常！！")
							replyText.ReplyToMessageID = webhook.Message.MessageID
							_, _ = bot.Send(replyText)
						}
					}()
					selenium.AdminLoginSaveToken(bot, webhook)
					replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, "已登录")
					replyText.ReplyToMessageID = webhook.Message.MessageID
					_, _ = bot.Send(replyText)
				}()
			}
		}
	}
}

// RestartService 测试环境 重启服务
func RestartService(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
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
				ReplyWithMessage(bot, webhook, "这是测试环境")
				task.ManageService(siteName, EnvName, command)
			}

		}
	}
}

// Help 打印帮助信息
func Help(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	msg := webhook.Message.Text
	botUsername := bot.Self.UserName
	if strings.Contains(msg, "@"+botUsername) {
		params := strings.Fields(msg)
		if len(params) > 1 && params[1] == "用法" {
			ReplyWithMessage(bot, webhook, "使用说明:\n"+
				"@机器人 <视图名> <任务名> 更新 - 指定视图名和任务名触发构建\n"+
				"示例: @CG33333_bot 28国际 后台API 更新\n"+
				"@CG33333_bot 28国际 后台API 获取分支\n"+
				"@CG33333_bot 28国际 生产 后台地址\n"+
				"@CG33333_bot 28国际 测试  后台API 更新")
		}
	}
}

// ReplyWithMessage 全局引用 用于小飞机发送消息
func ReplyWithMessage(bot *tgbotapi.BotAPI, webhook system.WebhookRequest, message string) {
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
	// 遍历收到的消息
	//update := bot.ListenForWebhook("/")
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
	GetBranch(bot, update)
	RestartService(bot, update)
}
