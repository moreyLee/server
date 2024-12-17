package task

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	modelSystem "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/cg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ReplyWithMessage 全局引用 用于小飞机发送消息
func ReplyWithMessage(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, message string) {
	replyText := tgbotapi.NewMessage(webhook.Message.Chat.ID, message)
	replyText.ReplyToMessageID = webhook.Message.MessageID
	_, _ = bot.Send(replyText)
	return
}

// GetBuildJobParam 获取JobName构建参数  根据参数选择环境   正式构建在调用
func GetBuildJobParam(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, JobName string, isProduction bool) map[interface{}]interface{} {
	// 根据 isProduction 参数选择URL
	var baseUrl, user, token string
	// 调用解密函数获取 Jenkins 配置
	jenkinsConfig, err := GetDecryptedJenkinsConfig()
	if err != nil {
		global.GVA_LOG.Error("获取 Jenkins 配置失败:", zap.Error(err))
		return nil
	}
	if isProduction {
		baseUrl = jenkinsConfig.Url
		user = jenkinsConfig.User
		token = jenkinsConfig.ApiToken
		fmt.Printf("生产url: " + baseUrl)
	} else {
		baseUrl = jenkinsConfig.TestUrl
		user = jenkinsConfig.TestUser
		token = jenkinsConfig.TestToken
		ReplyWithMessage(bot, webhook, "测试URL"+baseUrl)
		fmt.Printf("测试url " + baseUrl)
	}
	Url := baseUrl + "job/" + JobName + "/api/json?pretty=true"
	fmt.Println("合成的Url: " + Url)
	// 定义一个map  用于保存获取的构建参数
	//params := make(map[interface{}]interface{})
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		global.GVA_LOG.Error("创建GET请求失败:", zap.Error(err))
	}
	// 设置Auth Basic Auth  user &token
	req.SetBasicAuth(user, token)
	// 创建GET请求并获取响应
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("发送GET请求失败:", zap.Error(err))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			global.GVA_LOG.Error("关闭响应体失败:", zap.Error(err))
			return
		}
	}(resp.Body)
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	//fmt.Printf(string(body)) // 输出响应体
	return parseJenkinsJobParams(body) // 解析 JSON 数据
}

// parseJenkinsJobParams  解析 jenkins Job 的参数
func parseJenkinsJobParams(body []byte) map[interface{}]interface{} {
	// 定义一个 map 用于保存解析结果
	params := make(map[interface{}]interface{})
	// 解析 JSON 响应  结构体
	var job modelSystem.JenkinsJob
	err := json.Unmarshal(body, &job)
	if err != nil {
		global.GVA_LOG.Error("解析 JSON 失败:", zap.Error(err))
		return params
	}
	// 处理 Actions 中的参数定义
	for _, action := range job.Actions {
		for _, param := range action.ParameterDefinitions {
			// 构建参数赋值给 params map
			params[param.Name] = param.DefaultParameterValue.Value
			fmt.Printf("参数名称: %s, 默认值:  %v\n", param.Name, param.DefaultParameterValue.Value)
		}
	}
	// 处理 Property 中参数定义
	for _, property := range job.Property {
		for _, param := range property.ParameterDefinitions {
			// 构建参数写入map 中
			params[param.Name] = param.DefaultParameterValue.Value
			//fmt.Printf(" key: %s,  value: %v\n", param.Name, param.DefaultParameterValue.Value)
		}
	}
	return params
}

// GetExtName 获取JobName构建参数 获取jenkins job 的前缀名
func GetExtName(ViewName string) string {
	// 调用解密函数获取 Jenkins 配置
	jenkinsConfig, err := GetDecryptedJenkinsConfig()
	if err != nil {
		global.GVA_LOG.Error("获取 Jenkins 配置失败:", zap.Error(err))
		return ""
	}
	jenkinsUrl := jenkinsConfig.Url + "view/" + ViewName + "/api/json"
	// 定义一个map  用于保存获取的构建参数
	//params := make(map[interface{}]interface{})
	req, err := http.NewRequest("GET", jenkinsUrl, nil)
	if err != nil {
		global.GVA_LOG.Error("创建post请求失败:", zap.Error(err))
	}
	// 设置Auth Basic Auth  user &token
	req.SetBasicAuth(jenkinsConfig.User, jenkinsConfig.ApiToken)
	// 创建GET请求并获取响应
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("发送GET请求失败:Jenkins url 异常\n,", zap.Error(err))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			global.GVA_LOG.Error("关闭响应体失败:", zap.Error(err))
			return
		}
	}(resp.Body)
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	//fmt.Printf(string(body)) // 输出响应体
	// 解析json 响应
	var jenkinsView modelSystem.JenkinsView
	err = json.Unmarshal(body, &jenkinsView)
	// 检查 Jobs 是否为空
	if len(jenkinsView.Jobs) == 0 {
		global.GVA_LOG.Error("Jenkins 视图中没有任何任务")
		return ""
	}
	// 获取任务名前缀
	firstJob := jenkinsView.Jobs[0]
	global.GVA_LOG.Info("拼接后的名字%s" + firstJob.Name)
	extName := strings.SplitN(firstJob.Name, "_", 2)[0]
	fmt.Printf("后缀名%s", extName)
	return extName
}

// GetBranch  获取jenkins git仓库和代码分支
func GetBranch(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, ViewName string, typeName string) *modelSystem.JobConfig {
	if typeName == "" || ViewName == "" {
		log.Printf("视图名和任务名不能为空")
		ReplyWithMessage(bot, webhook, "视图名，或任务名不能为空")
		return nil
	}
	global.GVA_LOG.Info("映射名称来源: " + typeName)
	caseInsensitiveMap := createCaseInsensitiveMap(ExtensionMap)
	log.Printf("映射名称来源: " + typeName)
	// 使用不区分大小写的映射Map
	extName, exists := getCaseInsensitiveValue(caseInsensitiveMap, typeName)
	if !exists {
		ReplyWithMessage(bot, webhook, "未找到构建任务名！类型名: "+typeName+"视图名"+ViewName)
		return nil
	}
	global.GVA_LOG.Info("后缀: " + extName)
	preName := GetExtName(ViewName)
	global.GVA_LOG.Info("前缀: " + preName)
	fullName := preName + extName
	global.GVA_LOG.Info("拼接后任务名: %s: " + fullName)

	jenkinsUrl := global.GVA_CONFIG.Jenkins.Url + "view/" + ViewName + "/job/" + fullName + "/config.xml"
	global.GVA_LOG.Info("jenkins URL: " + jenkinsUrl)
	// 创建http 请求
	req, _ := http.NewRequest("GET", jenkinsUrl, nil)
	// 设置Auth Basic Auth  user &token 认证信息
	req.SetBasicAuth(global.GVA_CONFIG.Jenkins.User, global.GVA_CONFIG.Jenkins.ApiToken)
	// 创建GET请求并获取响应
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("发送GET请求失败 \n,", zap.Error(err))
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	bodyStr := strings.Replace(string(body), `<?xml version='1.1' encoding='UTF-8'?>`, `<?xml version='1.0' encoding='UTF-8'?>`, -1)
	bodyStr = strings.Replace(bodyStr, `<?xml version="1.1" encoding="UTF-8"?>`, `<?xml version="1.0" encoding="UTF-8"?>`, -1)

	// 解析xml 响应
	var jobConfig modelSystem.JobConfig
	if err := xml.Unmarshal([]byte(bodyStr), &jobConfig); err != nil {
		global.GVA_LOG.Error("解析xml响应失败, :", zap.Error(err))
	}

	return &jobConfig
}

// JenkinsBuildJobWithView 有参数构建 根据视图名构建
func JenkinsBuildJobWithView(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, ViewName string, JobName string, done chan modelSystem.JenkinsBuild) {
	defer func() {
		close(done) // 确保协程结束后关闭通道
	}()
	if JobName == "" || ViewName == "" {
		ReplyWithMessage(bot, webhook, "视图名或任务名不能为空")
		done <- modelSystem.JenkinsBuild{Success: false} // 通知失败
		return
	}
	// 获取Jenkins 的视图列表
	views := GetJenkinsViews()
	viewExists := false
	for _, view := range views {
		if view.Name == ViewName {
			fmt.Printf("%s\n", view.Name)
			viewExists = true
			break
		}
	}
	// 如果视图不存在，发送消息并退出
	if !viewExists {
		global.GVA_LOG.Info("视图名: %s 生产Jenkins中不存在\n" + ViewName)
		done <- modelSystem.JenkinsBuild{Success: false} // 管道通知任务失败
		return
	}
	// 引入不区分大小写的映射
	caseInsensitiveMap := createCaseInsensitiveMap(ExtensionMap)
	log.Printf("映射名称来源: " + JobName)
	// 使用不区分大小写的映射Map
	extName, exists := getCaseInsensitiveValue(caseInsensitiveMap, JobName)
	if !exists {
		ReplyWithMessage(bot, webhook, "未找到构建任务名！")
		done <- modelSystem.JenkinsBuild{Success: false}
		return
	}
	// 获取Job 前缀名称
	preName := GetExtName(ViewName)
	// 映射完 合成jobName名称
	fullJobName := preName + extName
	log.Printf("映射后的任务名Name: %s", fullJobName)
	//调用解密函数获取 Jenkins 配置
	jenkinsConfig, err := GetDecryptedJenkinsConfig()
	if err != nil {
		global.GVA_LOG.Error("获取 Jenkins 配置失败:", zap.Error(err))
		return
	}
	jenkinsUrl := jenkinsConfig.Url + "view/" + ViewName + "/job/" + fullJobName + "/buildWithParameters"
	// 判断一下 jenkins URL 是否有效
	_, err = url.ParseRequestURI(jenkinsUrl)
	if err != nil {
		global.GVA_LOG.Error("无效的jenkins URL:\n", zap.Error(err))
		done <- modelSystem.JenkinsBuild{Success: false}
		return
	}
	// 获取构建参数 params为map类型
	params := GetBuildJobParam(bot, webhook, fullJobName, true)
	if len(params) == 0 {
		JenkinsBuildJob(bot, webhook, ViewName, fullJobName, true, done) //  无参数构建jenkins
		global.GVA_LOG.Info(fmt.Sprintf("%s %s : 无参数Job任务已触发构建\n", ViewName, fullJobName))
		done <- modelSystem.JenkinsBuild{Success: true} // 通知成功
		return
	}
	// 表单数据 将获取的参数转换 表单数据    有参数构建
	data := url.Values{}
	for key, value := range params {
		fmt.Printf("参数名称: %s, 默认值: %v\n", key, value)
		req, err := http.NewRequest("POST", jenkinsUrl, bytes.NewBufferString(data.Encode()))
		if err != nil {
			global.GVA_LOG.Error("创建post请求失败:", zap.Error(err))
			done <- modelSystem.JenkinsBuild{Success: false} // 通知失败
			return
		}
		// 设置Auth Basic Auth
		req.SetBasicAuth(jenkinsConfig.User, jenkinsConfig.ApiToken)
		// 发送post请求并获取响应
		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			global.GVA_LOG.Error("发送post请求失败:", zap.Error(err))
			done <- modelSystem.JenkinsBuild{Success: false}
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			global.GVA_LOG.Info("构建请求成功，任务已启动 \n", zap.String("status", resp.Status))
			ReplyWithMessage(bot, webhook, ViewName+": 已成功触发构建 "+fullJobName+" 任务，30秒后获取构建状态")
			done <- modelSystem.JenkinsBuild{Success: true, FullJobName: fullJobName}
			return
		}
	}
}

// JenkinsBuildJob  无参数构建 支持生产环境和测试环境
func JenkinsBuildJob(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, ViewName string, Name string, isProduction bool, done chan modelSystem.JenkinsBuild) {
	// 输入视图名和任务名  只接受两个参数
	if Name == "" || ViewName == "" {
		fmt.Printf("错误: 需要传入视图名或任务名 至少一个参数！！")
		return
	}
	var jenkinsUrl string
	// 调用解密函数获取 Jenkins 配置
	jenkinsConfig, err := GetDecryptedJenkinsConfig()
	if err != nil {
		global.GVA_LOG.Error("获取 Jenkins 配置失败:", zap.Error(err))
		return
	}
	if isProduction {
		// 拼接生产环境 URL    变量屏蔽问题 :=  无法赋值给外层url
		jenkinsUrl = jenkinsConfig.Url + "view/" + ViewName + "/job/" + Name + "/build"
		global.GVA_LOG.Info(fmt.Sprintf("生产环境URL: %s\n", jenkinsUrl))
	} else {
		// 拼接测试环境 URL
		jenkinsUrl = jenkinsConfig.TestUrl + "view/" + ViewName + "/job/" + Name + "/build"
		global.GVA_LOG.Info(fmt.Sprintf("测试环境URL: %s\n", jenkinsUrl))
	}
	req, err := http.NewRequest("POST", jenkinsUrl, nil)
	if err != nil {
		global.GVA_LOG.Error("创建post请求失败:", zap.Error(err))
		return
	}
	// 设置Auth Basic Auth 生产环境 测试环境
	if isProduction {
		req.SetBasicAuth(jenkinsConfig.User, jenkinsConfig.ApiToken)
	} else {
		req.SetBasicAuth(jenkinsConfig.TestUser, jenkinsConfig.TestToken)
	}
	// 发送post请求并获取响应
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("发送post请求失败:", zap.Error(err))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		global.GVA_LOG.Info("构建请求成功，任务已启动 \n", zap.String("status", resp.Status))
		ReplyWithMessage(bot, webhook, ViewName+": 已成功触发构建 "+Name+" 无参任务，30秒后获取构建状态")
		done <- modelSystem.JenkinsBuild{Success: true, FullJobName: Name}
		return
	}
	global.GVA_LOG.Error("构建请求失败， 状态码:\n", zap.String("status", resp.Status))
	done <- modelSystem.JenkinsBuild{Success: false, FullJobName: Name}
}

// GetJenkinsViews 获取视图名和任务名 用于确定视图名是否存在
func GetJenkinsViews() []modelSystem.JenkinsView {
	// 调用解密函数 获取 Jenkins 配置
	jenkinsConfig, err := GetDecryptedJenkinsConfig()
	if err != nil {
		global.GVA_LOG.Error("获取 Jenkins 配置失败:", zap.Error(err))
		return nil
	}
	jenkinsUrl := jenkinsConfig.Url + "/api/json"
	req, err := http.NewRequest("GET", jenkinsUrl, nil)
	if err != nil {
		global.GVA_LOG.Error("创建GET请求失败:", zap.Error(err))
	}
	// 设置Auth Basic Auth
	req.SetBasicAuth(jenkinsConfig.User, jenkinsConfig.ApiToken)
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("发送post请求失败:", zap.Error(err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // 读取响应内容
		global.GVA_LOG.Error(fmt.Sprintf("请求失败状态码: %d,  response body: %s ", resp.StatusCode, string(body)))
	}
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Views []modelSystem.JenkinsView `json:"views"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		global.GVA_LOG.Error("未解析到数据:", zap.Error(err))
	}
	return data.Views
}

// GetDecryptedJenkinsConfig 解密 Jenkins 配置
func GetDecryptedJenkinsConfig() (*config.Jenkins, error) {
	decryptedUrl, err := cg.Decrypt(global.GVA_CONFIG.Jenkins.Url, global.GVA_CONFIG.Telegram.AesKey)
	if err != nil {
		global.GVA_LOG.Error("解密 Jenkins URL 失败:", zap.Error(err))
		return nil, err
	}

	decryptedUser, err := cg.Decrypt(global.GVA_CONFIG.Jenkins.User, global.GVA_CONFIG.Telegram.AesKey)
	if err != nil {
		global.GVA_LOG.Error("解密 Jenkins 用户名失败:", zap.Error(err))
		return nil, err
	}

	decryptedApiToken, err := cg.Decrypt(global.GVA_CONFIG.Jenkins.ApiToken, global.GVA_CONFIG.Telegram.AesKey)
	if err != nil {
		global.GVA_LOG.Error("解密 Jenkins API Token 失败:", zap.Error(err))
		return nil, err
	}
	decryptedTestUrl, err := cg.Decrypt(global.GVA_CONFIG.Jenkins.TestUrl, global.GVA_CONFIG.Telegram.AesKey)
	if err != nil {
		global.GVA_LOG.Error("解密 Jenkins URL 失败:", zap.Error(err))
		return nil, err
	}

	decryptedTestUser, err := cg.Decrypt(global.GVA_CONFIG.Jenkins.TestUser, global.GVA_CONFIG.Telegram.AesKey)
	if err != nil {
		global.GVA_LOG.Error("解密 Jenkins 用户名失败:", zap.Error(err))
		return nil, err
	}

	decryptedTestToken, err := cg.Decrypt(global.GVA_CONFIG.Jenkins.TestToken, global.GVA_CONFIG.Telegram.AesKey)
	if err != nil {
		global.GVA_LOG.Error("解密 Jenkins API Token 失败:", zap.Error(err))
		return nil, err
	}
	return &config.Jenkins{
		Url:       decryptedUrl,
		User:      decryptedUser,
		ApiToken:  decryptedApiToken,
		TestUrl:   decryptedTestUrl,
		TestUser:  decryptedTestUser,
		TestToken: decryptedTestToken,
	}, nil
}

// GetJobBuildStatus  获取jenkins Job 任务状态
func GetJobBuildStatus(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, ViewName string, JobName string, isProduction bool) (modelSystem.Build, error) {
	global.GVA_LOG.Info("视图名: " + ViewName + "任务名: " + JobName)
	// 引入的结构体
	var build modelSystem.Build
	var job modelSystem.JenkinsJob
	var jenkinsUrl string
	var body []byte //定义为 byte 类型  需要的时候进行类型转换
	var err error
	// 调用解密函数获取 Jenkins 配置
	jenkinsConfig, err := GetDecryptedJenkinsConfig()
	if err != nil {
		global.GVA_LOG.Error("获取 Jenkins 配置失败:", zap.Error(err))
		return modelSystem.Build{}, nil
	}
	if isProduction {
		// 生产环境
		jenkinsUrl = jenkinsConfig.Url + "view/" + ViewName + "/job/" + JobName + "/api/json"
		global.GVA_LOG.Info("生产环境URL: " + jenkinsUrl)

		body, err = GetRequestJkBody(jenkinsUrl, isProduction)
		global.GVA_LOG.Info("参数Body" + string(body))
		if err != nil {
			return modelSystem.Build{}, err
		}
	} else {
		// 测试环境
		global.GVA_LOG.Info("测试环境视图名%s,%v\n")
		jenkinsUrl = jenkinsConfig.TestUrl + "view/" + ViewName + "/job/" + JobName + "/api/json"
		global.GVA_LOG.Info("测试环境URL: " + jenkinsUrl)
		body, err = GetRequestJkBody(jenkinsUrl, isProduction)
		global.GVA_LOG.Info("测试环境参数Body" + string(body))
		if err != nil {
			return modelSystem.Build{}, err
		}
	}
	global.GVA_LOG.Info("实际的URL: \n" + jenkinsUrl)

	// 解析 Jenkins job 信息
	if err := json.Unmarshal(body, &job); err != nil {
		global.GVA_LOG.Error("解析失败 \n") // 记录错误日志
		return modelSystem.Build{}, fmt.Errorf("解析响应失败: %v", err)
	}
	// 检查是否有构建记录
	if len(job.Builds) == 0 {
		return modelSystem.Build{}, fmt.Errorf("未找到任何构建")
	}

	// 获取到最近的构建号
	recentlyNumber := job.LastBuild.Number
	var JobUrl string
	if isProduction {
		// 生产环境 构建号
		JobUrl = jenkinsConfig.Url + "view/" + ViewName + "/job/" + JobName + "/" + strconv.Itoa(recentlyNumber) + "/api/json"
		global.GVA_LOG.Info("生产环境 构建号URL: \n" + JobUrl)

	} else {
		// 测试环境 构建号
		JobUrl = jenkinsConfig.TestUrl + "view/" + ViewName + "/job/" + JobName + "/" + strconv.Itoa(recentlyNumber) + "/api/json"
		global.GVA_LOG.Info("测试环境 构建号URL: " + JobUrl + "\n")
	}
	global.GVA_LOG.Info("实际构建号URL: " + JobUrl + "\n")
	// 获取最近的构建详细信息
	if body, err = GetRequestJkBody(JobUrl, isProduction); err != nil {
		global.GVA_LOG.Info("构建详细信息: %v")
		return modelSystem.Build{}, fmt.Errorf("无法获取构建信息：: %v", err)

	}
	global.GVA_LOG.Info("\n Body 信息:  " + string(body) + "\n")
	// 解析构建信息
	if err := json.Unmarshal(body, &build); err != nil {
		return modelSystem.Build{},
			fmt.Errorf("解析响应失败: %v", zap.Error(err))
	}
	global.GVA_LOG.Info("构建信息: " + build.Result)
	// jenkins 变更记录
	var changes []string
	for _, item := range build.ChangeSet.Items {
		changes = append(changes, fmt.Sprintf("Commit: %s by %s - %s", item.CommitID, item.Author.FullName, item.Msg))
	}
	changeLog := "无变更记录"
	if len(changes) > 0 {
		changeLog = strings.Join(changes, "\n")
	}
	// 格式化时间
	formattedTime := FormatTimestamp(build.Timestamp, "Asia/Dubai")
	ReplyWithMessage(bot, webhook, fmt.Sprintf(
		"构建信息: %s\n最新构建编号: %d\n构建时间: %s\n更新描述: \n%s",
		build.Result, build.Number, formattedTime, changeLog))

	return build, nil
}

// FormatTimestamp 格式化时间戳为指定时区的时间字符串
func FormatTimestamp(timestamp int64, timezone string) string {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		global.GVA_LOG.Warn("加载时区失败，使用 UTC 代替:", zap.Error(err))
		location = time.UTC
	}
	return time.Unix(timestamp/1000, (timestamp%1000)*1000000).In(location).Format("2006-01-02 15:04:05")
}
