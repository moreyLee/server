package task

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	modelSystem "github.com/flipped-aurora/gin-vue-admin/server/model/system"
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
	if isProduction {
		baseUrl = global.GVA_CONFIG.Jenkins.Url
		user = global.GVA_CONFIG.Jenkins.User
		token = global.GVA_CONFIG.Jenkins.ApiToken
		fmt.Printf("生产url: " + baseUrl)
	} else {
		baseUrl = global.GVA_CONFIG.Jenkins.TestUrl
		user = global.GVA_CONFIG.Jenkins.TestUser
		token = global.GVA_CONFIG.Jenkins.TestToken
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
	jenkinsUrl := global.GVA_CONFIG.Jenkins.Url + "view/" + ViewName + "/api/json"
	// 定义一个map  用于保存获取的构建参数
	//params := make(map[interface{}]interface{})
	req, err := http.NewRequest("GET", jenkinsUrl, nil)
	if err != nil {
		global.GVA_LOG.Error("创建post请求失败:", zap.Error(err))
	}
	// 设置Auth Basic Auth  user &token
	req.SetBasicAuth(global.GVA_CONFIG.Jenkins.User, global.GVA_CONFIG.Jenkins.ApiToken)
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
func GetBranch(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, ViewName string, JobName string) *modelSystem.JobConfig {
	if JobName == "" || ViewName == "" {
		log.Printf("视图名和任务名不能为空")
		ReplyWithMessage(bot, webhook, "视图名，或任务名不能为空")
		return nil
	}
	log.Printf("映射名称来源: " + JobName)
	extName := ExtensionMap[JobName]
	log.Printf("后缀: %s", extName)
	preName := GetExtName(ViewName)
	log.Printf("前缀: %s ", preName)
	Name := preName + extName
	log.Printf("Name: %s", Name)

	jenkinsUrl := global.GVA_CONFIG.Jenkins.Url + "view/" + ViewName + "/job/" + Name + "/config.xml"
	log.Printf("jenkins URL: " + jenkinsUrl)
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
		global.GVA_LOG.Error("解析xml响应失败:", zap.Error(err))
		ReplyWithMessage(bot, webhook, "站点名称不存在")
	}

	return &jobConfig
}

// JenkinsBuildJobWithView 有参数构建 根据视图名构建
func JenkinsBuildJobWithView(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, ViewName string, JobName string, done chan bool) {
	defer func() {
		close(done) // 确保协程结束后关闭通道
	}()
	if JobName == "" || ViewName == "" {
		ReplyWithMessage(bot, webhook, "视图名或任务名不能为空")
		done <- false // 通知失败
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
		done <- false // 管道通知任务失败
		return
	}
	// 引入不区分大小写的映射
	caseInsensitiveMap := createCaseInsensitiveMap(ExtensionMap)
	log.Printf("映射名称来源: " + JobName)
	// 使用不区分大小写的映射Map
	extName, exists := getCaseInsensitiveValue(caseInsensitiveMap, JobName)
	if !exists {
		ReplyWithMessage(bot, webhook, "未找到构建任务名！")
		done <- false
		return
	}
	// 获取Job 前缀名称
	preName := GetExtName(ViewName)
	// 映射完 合成jobName名称
	Name := preName + extName
	log.Printf("映射后的任务名Name: %s", Name)

	jenkinsUrl := global.GVA_CONFIG.Jenkins.Url + "view/" + ViewName + "/job/" + Name + "/buildWithParameters"
	// 判断一下 jenkins URL 是否有效
	_, err := url.ParseRequestURI(jenkinsUrl)
	if err != nil {
		global.GVA_LOG.Error("无效的jenkins URL:\n", zap.Error(err))
		done <- false
		return
	}
	// 获取构建参数 params为map类型
	params := GetBuildJobParam(bot, webhook, Name, true)
	if len(params) == 0 {
		JenkinsBuildJob(bot, webhook, ViewName, Name, true) //  无参数构建jenkins
		fmt.Printf("%s %s : 无参数Job任务已触发构建\n", ViewName, Name)
		done <- true // 通知成功
		return
	}
	// 表单数据 将获取的参数转换 表单数据    有参数构建
	data := url.Values{}
	for key, value := range params {
		fmt.Printf("参数名称: %s, 默认值: %v\n", key, value)
		req, err := http.NewRequest("POST", jenkinsUrl, bytes.NewBufferString(data.Encode()))
		if err != nil {
			global.GVA_LOG.Error("创建post请求失败:", zap.Error(err))
			done <- false // 通知失败
			return
		}
		// 设置Auth Basic Auth
		req.SetBasicAuth(global.GVA_CONFIG.Jenkins.User, global.GVA_CONFIG.Jenkins.ApiToken)
		// 发送post请求并获取响应
		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			global.GVA_LOG.Error("发送post请求失败:", zap.Error(err))
			done <- false
			return
		}
		defer resp.Body.Close()
		switch resp.StatusCode {
		case http.StatusOK, http.StatusCreated:
			global.GVA_LOG.Info("构建请求成功，任务已启动 \n", zap.String("status", resp.Status))
			done <- true
		}
	}
}

// JenkinsBuildJob  无参数构建 支持生产环境和测试环境
func JenkinsBuildJob(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, ViewName string, Name string, isProduction bool) {
	// 输入视图名和任务名  只接受两个参数
	if Name == "" || ViewName == "" {
		fmt.Printf("错误: 需要传入视图名或任务名 至少一个参数！！")
		//ReplyWithMessage(bot,webhook,"错误: 需要传入视图名或任务名 至少一个参数！！")
		return
	}
	var jenkinsUrl string
	if isProduction {
		// 拼接生产环境 URL    变量屏蔽问题 :=  无法赋值给外层url
		jenkinsUrl = global.GVA_CONFIG.Jenkins.Url + "view/" + ViewName + "/job/" + Name + "/build"
		fmt.Printf("生产环境URL: %s\n", jenkinsUrl)
	} else {
		// 拼接测试环境 URL
		jenkinsUrl = global.GVA_CONFIG.Jenkins.TestUrl + "view/" + ViewName + "/job/" + Name + "/build"
		fmt.Printf("测试环境URL: %s\n", jenkinsUrl)
	}
	req, err := http.NewRequest("POST", jenkinsUrl, nil)
	if err != nil {
		global.GVA_LOG.Error("创建post请求失败:", zap.Error(err))
		return
	}
	// 设置Auth Basic Auth 生产环境 测试环境
	if isProduction {
		req.SetBasicAuth(global.GVA_CONFIG.Jenkins.User, global.GVA_CONFIG.Jenkins.ApiToken)
	} else {
		req.SetBasicAuth(global.GVA_CONFIG.Jenkins.TestUser, global.GVA_CONFIG.Jenkins.TestToken)
	}
	// 发送post请求并获取响应
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("发送post请求失败:", zap.Error(err))
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			global.GVA_LOG.Error("关闭响应体失败:", zap.Error(err))
		}
	}(resp.Body)
	// 检查响应状态
	if resp.StatusCode != http.StatusCreated {
		global.GVA_LOG.Error(fmt.Sprintf("构建请求失败，状态码: %d", resp.StatusCode))
		ReplyWithMessage(bot, webhook, fmt.Sprintf("Jenkins 构建请求失败，状态码: %d\n 视图名: %s\n 任务名: %s\n", resp.StatusCode, ViewName, Name))
	} else {
		fmt.Println("Jenkins 构建任务成功触发")
		ReplyWithMessage(bot, webhook, ViewName+": "+"已成功触发构建"+Name+"任务...")
	}
}

// GetJenkinsViews 获取视图名和任务名 用于确定视图名是否存在
func GetJenkinsViews() []modelSystem.JenkinsView {
	jenkinsUrl := global.GVA_CONFIG.Jenkins.Url + "/api/json"
	req, err := http.NewRequest("GET", jenkinsUrl, nil)
	if err != nil {
		global.GVA_LOG.Error("创建GET请求失败:", zap.Error(err))
	}
	// 设置Auth Basic Auth
	req.SetBasicAuth(global.GVA_CONFIG.Jenkins.User, global.GVA_CONFIG.Jenkins.ApiToken)
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("发送post请求失败:", zap.Error(err))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Views []modelSystem.JenkinsView `json:"views"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		global.GVA_LOG.Error("未解析到数据:", zap.Error(err))
	}
	return data.Views
}

// GetLastBuildStatus  获取jenkins Job 任务状态
func GetLastBuildStatus(bot *tgbotapi.BotAPI, webhook modelSystem.WebhookRequest, ViewName string, JobName string, isProduction bool) (modelSystem.Build, error) {
	// 引入的结构体
	var build modelSystem.Build
	var job modelSystem.JenkinsJob
	// 引入不区分大小写的映射
	caseInsensitiveMap := createCaseInsensitiveMap(ExtensionMap)
	// 使用不区分大小写的映射Map
	extName, _ := getCaseInsensitiveValue(caseInsensitiveMap, JobName)
	var Name, jenkinsUrl string
	var body []byte //定义为 byte 类型  需要的时候进行类型转换
	var err error
	if isProduction {
		// 生产环境
		// 获取视图 前缀名称
		preName := GetExtName(ViewName)
		// 映射完 合成新的 jobName名称
		Name := preName + extName
		fmt.Printf("生产环境 拼接后原始 jenkins Job 名称: %s\n", Name)
	} else {
		// 测试环境
		global.GVA_LOG.Info("测试环境  视图名，任务名" + ViewName + " " + JobName)
	}
	if isProduction {
		// 生产环境
		jenkinsUrl = global.GVA_CONFIG.Jenkins.Url + "view/" + ViewName + "/job/" + Name + "/api/json"
		global.GVA_LOG.Info("生产环境URL: " + jenkinsUrl)

		body, err = GetRequestJkBody(jenkinsUrl, isProduction)
		global.GVA_LOG.Info("参数Body" + string(body))
		if err != nil {
			return modelSystem.Build{}, err
		}
	} else {
		// 测试环境
		global.GVA_LOG.Info("测试环境视图名%s,%v\n")
		jenkinsUrl = global.GVA_CONFIG.Jenkins.TestUrl + "view/" + ViewName + "/job/" + JobName + "/api/json"
		global.GVA_LOG.Info("测试环境URL: " + jenkinsUrl)
		body, err = GetRequestJkBody(jenkinsUrl, isProduction)
		global.GVA_LOG.Info("测试环境参数Body" + string(body))
		if err != nil {
			return modelSystem.Build{}, err
		}
	}
	global.GVA_LOG.Info("实际的URL: \n" + jenkinsUrl)
	if jenkinsUrl == "" {
		global.GVA_LOG.Error("没有获取到Jenkins Url \n") // 记录错误日志
		return modelSystem.Build{}, fmt.Errorf("")  // 返回错误信息
	}
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
		JobUrl = global.GVA_CONFIG.Jenkins.Url + "view/" + ViewName + "/job/" + Name + "/" + strconv.Itoa(recentlyNumber) + "/api/json"
		global.GVA_LOG.Info("生产环境 构建号URL: \n" + JobUrl)

	} else {
		// 测试环境 构建号
		JobUrl = global.GVA_CONFIG.Jenkins.TestUrl + "view/" + ViewName + "/job/" + JobName + "/" + strconv.Itoa(recentlyNumber) + "/api/json"
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
	formattedTime := time.Unix(build.Timestamp/1000, (build.Timestamp%1000)*1000000).UTC().In(func() *time.Location { loc, _ := time.LoadLocation("Asia/Dubai"); return loc }()).Format("2006-01-02 15:04:05")
	ReplyWithMessage(bot, webhook, fmt.Sprintf(
		"构建信息: %s\n最新构建编号: %d\n构建时间: %s\n更新描述: \n%s",
		build.Result, build.Number, formattedTime, changeLog))

	return build, nil
}
