package system

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	system2 "github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/task"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/selenium"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strings"
	"time"
)

var jkBuild system.JenkinsBuild

func BuildJenkins(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) ([]string, bool) {
	msg := webhook.Message.Text
	botUsername := bot.Self.UserName
	// 检查消息中是否包含机器人用户名
	if !strings.Contains(msg, "@"+botUsername) {
		return nil, false
	}
	// 获取外部参数
	params := strings.Fields(msg)
	log.Printf("参数数量: %d, 参数内容: %v", len(params), params)
	// 引入标志变量 确保只输出一次
	// 检查参数数量
	if len(params) < 5 { // 不要改 容易数组内存越界
		// 只返回是否有效 不做消息返回 又调用者输出消息
		global.GVA_LOG.Error("输入参数不足，请参考用法:  @机器人 用法")
		ReplyWithMessage(bot, webhook, "输入参数不足，请参考用法:  @机器人 用法 ")
		return nil, false
	}
	// 检查环境参数  第一个参数  参数内容
	envName := params[1]
	validEnvs := []string{"正式", "正式环境", "生产", "生产环境", "测试", "测试环境"}
	if !validateEnv(envName, validEnvs, bot, webhook) {
		expectedEnvs := strings.Join(validEnvs, " | ")
		global.GVA_LOG.Info("第一个参数 环境参数无效，请任选其一输入 '" + expectedEnvs + "'")
		ReplyWithMessage(bot, webhook, "第一个参数 环境参数无效，请任选其一输入 '"+expectedEnvs+"'")
		return nil, false
	}
	prodEnvs := []string{"正式", "正式环境", "生产", "生产环境"}
	if validateEnv(envName, prodEnvs, bot, webhook) {
		BuildJobsWithProd(bot, webhook)

	}
	testEnvs := []string{"测试", "测试环境"}
	if validateEnv(envName, testEnvs, bot, webhook) {
		BuildJobsWithEnvTest(bot, webhook)

	}

	//if validateEnv()
	return params, true
}

// BuildJobsWithText 基于文本消息构建jenkins 任务   生产环境
func BuildJobsWithProd(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	msg := webhook.Message.Text
	params := strings.Fields(msg)
	// 通道done 确保任务状态反馈
	done := make(chan bool)
	if params == nil || len(params) == 0 {
		return
	}
	fmt.Printf("生产params: %v, valid: %v", params)
	//// 验证生产环境参数 第一个参数内容
	//envName := params[1]
	//validEnvs := []string{"正式", "正式环境", "生产", "生产环境"}
	//// 检查生产环境参数是否有效
	//isValidEnv := false
	//for _, env := range validEnvs {
	//	if env == envName {
	//		isValidEnv = true
	//		break
	//	}
	//}
	//if !isValidEnv {
	//	global.GVA_LOG.Error("生产环境 参数不正确")
	//	ReplyWithMessage(bot, webhook, "生产环境 参数不正确")
	//	return
	//}
	// 验证站点参数  第二个参数 检查生产站点名称是否存在
	SiteName := params[2]
	tgService := system2.TgService{}
	// 创建实例并赋值
	jkBuild := system.JenkinsBuild{
		ViewName: SiteName,
	}
	exists, err := tgService.CheckJenView(jkBuild.ViewName, "jenkins_env_prod", "prod_site_name")
	if err != nil {
		global.GVA_LOG.Error(" 数据库错误: ", zap.Error(err))
		ReplyWithMessage(bot, webhook, "站点名称未找到")
		return
	}
	if !exists {
		ReplyWithMessage(bot, webhook, fmt.Sprintf("站点名 %s 不存在，请检查输入是否正确", SiteName))
		return
	}
	// 验证任务类型 第三个参数
	jkBuild.TaskType = params[3] // 任务类型(后台API)
	if !validateType(jkBuild.TaskType, bot, webhook) {
		return
	}
	// 验证构建类型   第四个参数 内容必须为 更新 查分支
	action := params[4] // 指令
	if !validateAction(action, bot, webhook) {
		return
	}
	if action != "更新" {
		global.GVA_LOG.Error("构建参数必须为 更新！！")
		return
	}

	// 启动协程执行构建任务
	go func() {
		defer func() {
			if err := recover(); err != nil {
				// 异常捕获 业务逻辑处理
				global.GVA_LOG.Error("构建任务中出现异常错误: %v", zap.Any("error", err))
			}
		}()
		// 检查用户权限
		if isUserAuthorized(webhook.Message.From.UserName) {
			task.JenkinsBuildJobWithView(bot, webhook, jkBuild.ViewName, jkBuild.TaskType, done)
		} else {
			global.GVA_LOG.Error("该用户无权限执行构建任务！！")
		}
	}()
	// 主协程 也要等待任务完成并获取构建状态
	// select 语句监听 done 通道
	select {
	case success := <-done:
		if success {
			// 加业务逻辑容易导致死循环
			global.GVA_LOG.Info("协程返回成功，进入构建状态发送\n")
			go func() {
				time.Sleep(30 * time.Second)
				//sendBuildStatus(bot, webhook)
			}()
		} else {
			global.GVA_LOG.Error("主协程通知 视图: " + jkBuild.ViewName + " 不存在构建任务失败，无法获取构建任务状态\n") // 别删
		}
	case <-time.After(15 * time.Second):
		global.GVA_LOG.Info("主协程通知 构建任务超时\n")
	}
}

// BuildJobsWithEnvTest  测试环境 更新Jenkins
func BuildJobsWithEnvTest(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	msg := webhook.Message.Text
	params := strings.Fields(msg)
	// 通道done 确保任务状态反馈
	done := make(chan bool)
	if params == nil || len(params) == 0 {
		return
	}
	fmt.Printf("测试params: %v, valid: %v", params)
	// 验证生产环境参数 第一个参数内容
	//envName := params[1]
	//validEnvs := []string{"测试", "测试环境"}
	//// 检查生产环境参数是否有效
	//isValidEnv := false
	//for _, env := range validEnvs {
	//	if env == envName {
	//		isValidEnv = true
	//		break
	//	}
	//}
	//if !isValidEnv {
	//	global.GVA_LOG.Error("测试环境 参数不正确")
	//	ReplyWithMessage(bot, webhook, "测试环境 参数不正确")
	//	return
	//}
	// 验证站点信息 第二个参数

	// 验证任务类型
	jkBuild.ViewName = params[2] // 站点名
	jkBuild.TaskType = params[3] // 任务类型(后台API)    注意一下
	if !validateType(jkBuild.TaskType, bot, webhook) {
		return
	}
	// 验证触发动作
	action := params[4] // 指令
	if !validateAction(action, bot, webhook) {
		return
	}

	// 启动协程执行构建任务
	go func() {
		defer func() {
			if err := recover(); err != nil {
				// 异常捕获 业务逻辑处理
				global.GVA_LOG.Error("构建任务中出现异常错误: %v", zap.Any("error", err))
			}
		}()
		if isUserAuthorized(webhook.Message.From.UserName) {
			task.JenkinsJobsWithTest(bot, webhook, jkBuild.ViewName, jkBuild.TaskType)
		}
	}()
	// 主协程 也要等待任务完成并获取构建状态
	// select 语句监听 done 通道
	select {
	case success := <-done:
		if success {
			// 加业务逻辑容易导致死循环
			global.GVA_LOG.Info("协程返回成功，进入构建状态发送\n")
			go func() {
				time.Sleep(30 * time.Second)
			}()
		} else {
			global.GVA_LOG.Error(jkBuild.ViewName + ": 主协程通知 测试环境 视图名不存在构建任务失败，无法获取构建任务状态\n") // 别删
		}
	case <-time.After(15 * time.Second):
		global.GVA_LOG.Info("主协程通知 构建任务超时\n")
	}
}

// isUserAuthorized 检查用户是否被授权
func isUserAuthorized(username string) bool {
	bot := &tgbotapi.BotAPI{}
	webhook := system.WebhookRequest{}
	//authorizedUsers := []string{"David5886", "nikon_aaa", "tank9999999"}
	authorizedUsers := global.GVA_CONFIG.Telegram.AuthorizedUser
	for _, user := range authorizedUsers {
		if username == user {
			return true
		}
	}
	ReplyWithMessage(bot, webhook, "用户未授权 请联系运维 @David5886")
	return false
}

// 环境参数验证  验证参数内容
func validateEnv(envName string, validEnvs []string, bot *tgbotapi.BotAPI, webhook system.WebhookRequest) bool {
	// 检查输入的环境参数是否有效
	for _, validEnv := range validEnvs {
		if envName == validEnv {
			return true
		}
	}
	//expectedEnvs := strings.Join(validEnvs, "或")
	//global.GVA_LOG.Info("第一个参数 环境参数无效，请任选其一输入 '" + expectedEnvs + "'")
	return false
}

// 构建任务类型 验证参数内容
func validateType(taskType string, bot *tgbotapi.BotAPI, webhook system.WebhookRequest) bool {
	validTypes := []string{"后台API", "后台api", "前台API", "前台api", "后台H5", "后台h5", "前台H5", "前台h5", "定时任务"}
	for _, validType := range validTypes {
		if taskType == validType {
			return true
		}
	}
	expectedEnvs := strings.Join(validTypes, "或")
	ReplyWithMessage(bot, webhook, "第三个参数 类型参数无效，请任选其一输入 '"+expectedEnvs+"'")
	return false
}
func validateAction(actionName string, bot *tgbotapi.BotAPI, webhook system.WebhookRequest) bool {
	validActions := []string{"更新", "查分支"}
	// 检查输入的环境参数是否有效
	for _, action := range validActions {
		if actionName == action {
			return true
		}
	}
	expectedEnvs := strings.Join(validActions, "或")
	ReplyWithMessage(bot, webhook, "第五个构建动作参数无效，请任选其一输入 '"+expectedEnvs+"'")
	return false
}

// sendBuildStatus 异步发送构建状态信息 避免阻塞主协程的其他操作  主协程不需要等待构建状态发送函数的执行 可以更快响应其他操作
func sendBuildStatus(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	go func() {
		if jkBuild.ViewName != "" && jkBuild.JobName != "" {
			task.GetLastBuildStatus(bot, webhook, jkBuild.ViewName, jkBuild.JobName, false)
		}
	}()
}

// GetBranch 获取构建任务 git 分支
func GetBranch(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	msg := webhook.Message.Text
	botUsername := bot.Self.UserName
	if !strings.Contains(msg, "@"+botUsername) {
		return
	}
	// 获取外部参数
	params := strings.Fields(msg)
	log.Printf("分支参数: %d, 参数内容: %v", len(params), params)

	// 检验环境参数
	envName := params[1]
	validEnvs := []string{"正式环境", "生产环境", "生产", "正式"}
	if !validateEnv(envName, validEnvs, bot, webhook) {
		return
	}
	breachName := params[4]
	if !validateAction(breachName, bot, webhook) {
		return
	}
	jkBuild.ViewName = params[2] // 参数2  视图名
	jkBuild.JobName = params[3]  //
	if breachName == "查分支" {
		jobConfig := task.GetBranch(bot, webhook, jkBuild.ViewName, jkBuild.JobName)
		if len(jobConfig.SCM.UserRemoteConfigs.URLs) > 0 && len(jobConfig.SCM.Branches) > 0 && len(jobConfig.SCM.UserRemoteConfigs.URLs) > 0 {
			replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, "生产环境\n"+jkBuild.ViewName+" "+jkBuild.JobName+"git分支:  "+jobConfig.SCM.Branches[0]+"\n"+
				jkBuild.ViewName+" "+jkBuild.JobName+"仓库地址:  "+jobConfig.SCM.UserRemoteConfigs.URLs[0])
			replyText.ReplyToMessageID = webhook.Message.MessageID
			_, _ = bot.Send(replyText)
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
					siteLink := selenium.GetAdminLinkTools(bot, webhook, siteName)
					message := siteName + "站点地址: " + siteLink
					replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, message)
					replyText.ReplyToMessageID = webhook.Message.MessageID
					_, _ = bot.Send(replyText)
				}()
			}
			if params[3] == "后台图片" {
				go func() {
					defer func() {
						if err := recover(); err != nil {
							replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, "后台登录异常！！")
							replyText.ReplyToMessageID = webhook.Message.MessageID
							_, _ = bot.Send(replyText)
						}
					}()
					err := selenium.GetAdminLinkTool(bot, webhook, siteName)
					if err != nil {
						global.GVA_LOG.Error("获取后台链接图片失败", zap.Any("err", err))
					}
					ReplyWithMessage(bot, webhook, "已发送")
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
				ReplyWithMessage(bot, webhook, "这是测试环境"+siteName+EnvName+command)
				//task.ManageService(siteName, EnvName, command)
			}

		}
	}
}

// Help 打印帮助信息
func Help(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	msg := webhook.Message.Text
	botUsername := bot.Self.UserName
	if !strings.Contains(msg, "@"+botUsername) {
		return
	}
	params := strings.Fields(msg)
	if len(params) > 1 && params[0] == global.GVA_CONFIG.Telegram.BotName && params[1] == "用法" {
		ReplyWithMessage(bot, webhook, "使用说明:\n"+
			"@机器人 <视图名> <任务名> 更新 - 指定视图名和任务名触发构建\n"+
			"示例: @机器人  28国际 生产环境 后台API 更新\n"+
			"@机器人 AK国际 生产环境 后台API 查分支 \n"+
			"@CG33333_bot 28国际 生产 后台地址 更新 \n"+
			"@CG33333_bot 28国际 测试  亿万-T1(自营) 更新")
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
	BuildJenkins(bot, update)
	GetAdminLink(bot, update) // 获取总后台站点链接
	//BuildJobsWithText(bot, update)    // 生产环境  更新jenkins代码
	//BuildJobsWithEnvTest(bot, update) // 测试环境 更新jenkins代码
	//GetBranch(bot, update)            //  生产环境  获取项目分支与仓库
	//RestartService(bot, update)
	//Help(bot, update)
	//sendBuildStatus(bot, update)
}
