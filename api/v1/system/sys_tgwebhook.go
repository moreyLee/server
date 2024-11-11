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
	switch {
	case len(params) == 2:
		usage := params[1]
		if usage == "用法" {
			Help(bot, webhook)
			return nil, true
		} else {
			return nil, false
		}
	case len(params) == 4: //@CG33333_bot 正式 速8娱乐  后台链接
		envName := params[1]
		validEnvs := []string{"正式", "正式环境", "生产", "生产环境"}
		if !validateEnv(envName, validEnvs, bot, webhook) {
			expectedEnvs := strings.Join(validEnvs, " | ")
			global.GVA_LOG.Info("第一个参数 环境参数无效，请任选其一输入 '" + expectedEnvs + "'")
			ReplyWithMessage(bot, webhook, "第一个参数 环境参数无效，请任选其一输入 '"+expectedEnvs+"'")
			return nil, false
		}
		linkEnvs := []string{"正式", "正式环境", "生产", "生产环境"}
		if validateEnv(envName, linkEnvs, bot, webhook) {
			GetAdminLink(bot, webhook)
		}
		return nil, true
	case len(params) == 5: // @CG33333_bot 正式 AK国际  后台API 更新
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
			GetBranch(bot, webhook) // 待优化
		}
		testEnvs := []string{"测试", "测试环境"}
		if validateEnv(envName, testEnvs, bot, webhook) {
			BuildJobsWithEnvTest(bot, webhook)
			RestartService(bot, webhook)
			ExecShell(bot, webhook)
		}
	case len(params) < 5:
		global.GVA_LOG.Error("输入参数不足，请参考用法:  @机器人 用法")
		ReplyWithMessage(bot, webhook, "输入参数不足，请参考用法:  @机器人 用法 ")
		return nil, false
	default:
		//global.GVA_LOG.Error("不满足触发条件")
	}
	return params, true
}

// BuildJobsWithProd 基于文本消息构建jenkins 任务   生产环境
func BuildJobsWithProd(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {
	msg := webhook.Message.Text
	params := strings.Fields(msg)
	// 通道done 确保任务状态反馈
	done := make(chan system.JenkinsBuild)
	if params == nil || len(params) == 0 {
		return
	}
	global.GVA_LOG.Info(fmt.Sprintf("生产params: %v, valid: %v", params))
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
		ReplyWithMessage(bot, webhook, "数据库中未找到站点名称")
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
		global.GVA_LOG.Info("构建参数必须为 更新！！")
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
		if isUserAuthorized(bot, webhook, webhook.Message.From.UserName) {
			task.JenkinsBuildJobWithView(bot, webhook, jkBuild.ViewName, jkBuild.TaskType, done)
		} else {
			global.GVA_LOG.Error("该用户无权限执行构建任务！！")
			done <- system.JenkinsBuild{Success: false}
		}
	}()
	// 主协程 也要等待任务完成并获取构建状态
	// select 语句监听 done 通道
	select {
	case result := <-done:
		if result.Success {
			// 加业务逻辑容易导致死循环
			global.GVA_LOG.Info("协程返回成功，进入构建状态发送\n")
			go func() {
				time.Sleep(30 * time.Second)
				task.GetJobBuildStatus(bot, webhook, jkBuild.ViewName, result.FullJobName, true)
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
	status := make(chan system.JenkinsBuild)
	// 通道done 确保任务状态反馈
	done := make(chan bool)
	if params == nil || len(params) == 0 {
		return
	}
	fmt.Printf("测试params: %v, valid: %v", params)
	// 验证站点信息 第二个参数

	// 验证任务类型
	jkBuild.ViewName = params[2] // 站点名
	jkBuild.TaskType = params[3] // 任务类型(后台API)    注意一下
	if !validateType(jkBuild.TaskType, bot, webhook) {
		return
	}
	// 验证触发动作
	action := params[4] // 指令

	if action != "更新" {
		ReplyWithMessage(bot, webhook, "测试环境  触发构建参数 更新 ")
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
		if isUserAuthorized(bot, webhook, webhook.Message.From.UserName) {
			task.JenkinsJobsWithTest(bot, webhook, jkBuild.ViewName, jkBuild.TaskType, status)
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
func isUserAuthorized(bot *tgbotapi.BotAPI, webhook system.WebhookRequest, username string) bool {
	//authorizedUsers := []string{"David5886", "nikon_aaa", "tank9999999"}
	authorizedUsers := global.GVA_CONFIG.Telegram.AuthorizedUser
	for _, user := range authorizedUsers {
		if username == user {
			return true
		}
	}
	ReplyWithMessage(bot, webhook, "用户 @"+username+" 未授权 请联系夜班运维 @David5886")
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
	expectedEnvs := strings.Join(validTypes, " | ")
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
	expectedEnvs := strings.Join(validActions, " | ")
	ReplyWithMessage(bot, webhook, "第五个参数 构建参数无效，请任选其一输入 '"+expectedEnvs+"'")
	return false
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
	if breachName != "查分支" {
		return
	}
	jkBuild.ViewName = params[2] // 参数2  视图名
	jkBuild.JobName = params[3]  //
	jobConfig := task.GetBranch(bot, webhook, jkBuild.ViewName, jkBuild.JobName)
	if len(jobConfig.SCM.UserRemoteConfigs.URLs) > 0 && len(jobConfig.SCM.Branches) > 0 && len(jobConfig.SCM.UserRemoteConfigs.URLs) > 0 {
		replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, "生产环境\n"+jkBuild.ViewName+" "+jkBuild.JobName+"git分支:  "+jobConfig.SCM.Branches[0]+"\n"+
			jkBuild.ViewName+" "+jkBuild.JobName+"仓库地址:  "+jobConfig.SCM.UserRemoteConfigs.URLs[0])
		replyText.ReplyToMessageID = webhook.Message.MessageID
		_, _ = bot.Send(replyText)
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
			siteName := params[2]
			if params[3] == "后台链接" {
				go func() {
					defer func() {
						if err := recover(); err != nil {
							code, _ := selenium.GetCaptchaCode(bot, webhook)
							global.GVA_LOG.Error("获取后台链接异常", zap.Any("err", err))
							ReplyWithMessage(bot, webhook, "有效验证码 4位数字"+code)
						}
					}()
					selenium.GetAdminLinkTools(bot, webhook, siteName)
				}()
			}
			if params[3] == "后台图片" {
				go func() {
					defer func() {
						if err := recover(); err != nil {
							ReplyWithMessage(bot, webhook, "后台登录")
						}
					}()
					err := selenium.GetAdminLinkPhoto(bot, webhook, siteName)
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
	if !strings.Contains(msg, "@"+botUsername) {
		return
	}
	params := strings.Fields(msg)
	log.Printf("分支参数: %d, 参数内容: %v", len(params), params)

	//siteName := params[2] // 测试环境 站点名称

	Optype := params[3]                      // 执行类型
	if !validateType(Optype, bot, webhook) { // 验证执行类型参数
		return
	}

}

// ExecShell  测试环境 执行脚本
func ExecShell(bot *tgbotapi.BotAPI, webhook system.WebhookRequest) {

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
			"示例: @机器人  环境 站点 任务类型  触发动作\n"+
			"构建Jenkins 需授权用户"+
			"@机器人  生产环境 28国际 后台API 查分支 \n"+
			"@CG33333_bot 生产环境 28国际 后台API 更新 \n"+
			"@CG33333_bot 测试环境 T1 定时任务 更新 \n"+
			"@CG33333_bot 测试环境 T1 定时任务 查分支")
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
}
