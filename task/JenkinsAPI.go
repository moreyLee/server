package task

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	modelSystem "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	system2 "github.com/flipped-aurora/gin-vue-admin/server/service/system"
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
}

// GetBuildJobParam 获取JobName构建参数  根据参数选择环境
func GetBuildJobParam(JobName string, isProduction bool) map[interface{}]interface{} {
	fmt.Printf("是否调用")
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
	fmt.Printf("拼接后的名字%s", firstJob.Name)
	extName := strings.SplitN(firstJob.Name, "_", 2)[0]
	fmt.Printf("后缀名%s", extName)
	return extName
}

// GetBranch  获取jenkins git仓库和代码分支
func GetBranch(ViewName string, JobName string) *modelSystem.JobConfig {
	if JobName == "" || ViewName == "" {
		log.Printf("视图名和任务名不能为空")
	}
	var extensionMap = map[string]string{
		"后台API": "_adminapi",
		"前台API": "_api",
		"前台H5":  "_h5",
		"后台H5":  "_h5admin",
		"定时任务":  "_quartz",
	}
	log.Printf("映射名称来源: " + JobName)
	extName := extensionMap[JobName]
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
		ReplyWithMessage(bot, webhook, ViewName+"站点不存在，构建失败，请检查站点名称")
		log.Printf("视图名: %s Jenkins中不存在\n", ViewName)
		done <- false // 管道通知任务失败
		return
	}
	ReplyWithMessage(bot, webhook, ViewName+" "+JobName+" 已触发构建，请稍等... ")
	var extensionMap = map[string]string{
		"后台API": "_adminapi",
		"前台API": "_api",
		"前台H5":  "_h5",
		"后台H5":  "_h5admin",
		"定时任务":  "_quartz",
	}
	// 引入不区分大小写的映射
	caseInsensitiveMap := createCaseInsensitiveMap(extensionMap)
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
	params := GetBuildJobParam(Name, true)
	if len(params) == 0 {
		JenkinsBuildJob(ViewName, Name)
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
			//case http.StatusNotFound:
			//	global.GVA_LOG.Error("构建请求失败，资源未找到", zap.String("status", resp.Status))
			//	done <- false
			//default:
			//	global.GVA_LOG.Error("构建请求失败", zap.String("status", resp.Status))
			//	done <- false
		}
	}
}

// JenkinsBuildJob  无参数构建/**
func JenkinsBuildJob(ViewName string, Name string) {
	jenkinsUrl := global.GVA_CONFIG.Jenkins.Url + "view/" + ViewName + "/job/" + Name + "/build"
	req, err := http.NewRequest("POST", jenkinsUrl, nil)
	if err != nil {
		global.GVA_LOG.Error("创建post请求失败:", zap.Error(err))
		return
	}
	// 设置Auth Basic Auth
	req.SetBasicAuth(global.GVA_CONFIG.Jenkins.User, global.GVA_CONFIG.Jenkins.ApiToken)
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
func GetLastBuildStatus(ViewName string, JobName string) (modelSystem.Build, error) {
	var extensionMap = map[string]string{
		"后台API": "_adminapi",
		"前台API": "_api",
		"前台H5":  "_h5",
		"后台H5":  "_h5admin",
		"定时任务":  "_quartz",
	}
	// 引入不区分大小写的映射
	caseInsensitiveMap := createCaseInsensitiveMap(extensionMap)
	log.Printf("映射名称来源: " + JobName)
	// 使用不区分大小写的映射Map
	extName, _ := getCaseInsensitiveValue(caseInsensitiveMap, JobName)
	// 获取前缀名称
	preName := GetExtName(ViewName)
	// 映射完 合成jobName名称
	Name := preName + extName
	log.Printf("Name: %s", Name)
	jenkinsUrl := global.GVA_CONFIG.Jenkins.Url + "/view/" + ViewName + "/job/" + Name + "/api/json"
	body, err := GetJenkinsData(jenkinsUrl)
	if err != nil {
		return modelSystem.Build{}, err
	}
	//打印原始响应以进行调试
	//fmt.Printf("API响应:%v", string(body))
	var job modelSystem.JenkinsJob
	if err := json.Unmarshal(body, &job); err != nil {
		return modelSystem.Build{}, fmt.Errorf("解析响应失败: %v", err)
	}

	if len(job.Builds) == 0 {
		return modelSystem.Build{}, fmt.Errorf("未找到任何构建")
	}

	// 获取到最近的构建号
	recentlyNumber := job.LastBuild.Number
	JobURL := global.GVA_CONFIG.Jenkins.Url + "view/" + ViewName + "/job/" + Name + "/" + strconv.Itoa(recentlyNumber) + "/api/json"
	body, err = GetJenkinsData(JobURL)
	if err != nil {
		return modelSystem.Build{}, err
	}
	var build modelSystem.Build
	if err := json.Unmarshal(body, &build); err != nil {
		return modelSystem.Build{}, fmt.Errorf("解析响应失败: %v", err)
	}
	// 返回最近的构建
	return build, nil
}

// ManageService 测试环境 重启应用服务 参数 站点名称 环境参数  命令映射
func ManageService(siteName string, EnvName string, command string) {
	type OrderedMap struct {
		Keys   []string
		Values map[string]interface{}
	}
	if siteName == "" {
		log.Printf("站点信息为空")
	}
	jobName := "exec_command"
	server := "csapi"
	// 获取测试环境执行命令job的参数   false 表示测试环境
	params := GetBuildJobParam(jobName, false)
	fmt.Printf("jenkins 参数 parms: %v\n", params)
	// 从数据库中获取站点信息 参数入参  站点名称 或 站点ID   返回值 站点名称对应的一条记录
	var siteConfig modelSystem.YzSiteConfig
	siteConfig.SiteName = siteName
	// 创建 TgService 实例
	tgService := system2.TgService{}
	err := tgService.SelectBySiteName(&siteConfig)
	if err != nil {
		log.Printf("查询数据为空")
	}
	fmt.Printf("对象 %v\n", siteConfig)
	siteID := siteConfig.SiteID
	siteName = siteConfig.SiteName
	siteDomain := siteConfig.Domains
	beginTime := siteConfig.BeginTime
	fmt.Printf("站点ID: %v\n站点名称: %v\n开站时间: %v\n建站时间: %v\n", siteID, siteName, siteDomain, beginTime)
	// 创建表单数据 数据库的值 赋值给 jenkins
	data := url.Values{} // 用于处理查询jenkins URL参数或表单数据
	for key, value := range params {
		data.Set(key.(string), value.(string))
	}
	fmt.Printf("Jenkins赋值给表单data%v\n", data)
	data["site_id"] = []string{strconv.Itoa(int(siteID))}
	data["site"] = []string{siteName}
	data["server"] = []string{server}
	fmt.Printf("jenkins修改后%v\n", data)
}
