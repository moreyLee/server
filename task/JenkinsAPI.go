package task

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// GetBuildJobParam 获取JobName构建参数
func GetBuildJobParam(JobName string) map[interface{}]interface{} {
	jenkinsUrl := global.GVA_CONFIG.Jenkins.Url + "job/" + JobName + "/api/json?pretty=true"
	// 定义一个map  用于保存获取的构建参数
	params := make(map[interface{}]interface{})
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
	//fmt.Printf(string(body))  // 输出响应体
	// 解析json 响应
	var job system.JenkinsJob
	err = json.Unmarshal(body, &job)
	// 输出构建参数的名称和值
	for _, action := range job.Actions {
		for _, param := range action.ParameterDefinitions {
			// 构建参数赋值给 params map
			params[param.Name] = param.DefaultParameterValue.Value
			//fmt.Printf("参数名称: %s, 默认值:  %v\n", param.Name, param.DefaultParameterValue.Value)
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
	var jenkinsView system.JenkinsView
	err = json.Unmarshal(body, &jenkinsView)
	firstJob := jenkinsView.Jobs[0]
	//params[firstJob.Name] = firstJob.Name
	//fmt.Printf("%s", firstJob.Name)
	extName := strings.SplitN(firstJob.Name, "_", 2)[0]
	return extName
}

// GetBranch  获取jenkins git仓库和代码分支
func GetBranch(ViewName string, JobName string) *system.JobConfig {
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
	//bodyStr := strings.Replace(string(body), `<?xml version='1.1' encoding='UTF-8'?>`, `<?xml version='1.0' encoding='UTF-8'?>`, 1)

	// 解析xml 响应
	var jobConfig system.JobConfig
	if err := xml.Unmarshal([]byte(bodyStr), &jobConfig); err != nil {
		global.GVA_LOG.Error("解析xml响应失败:", zap.Error(err))
	}

	return &jobConfig
}

// JenkinsBuildJobWithView 有参数构建 根据视图名构建
func JenkinsBuildJobWithView(ViewName string, JobName string) {
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
	// 引入不区分大小写的映射
	caseInsensitiveMap := createCaseInsensitiveMap(extensionMap)
	log.Printf("映射名称来源: " + JobName)
	// 使用不区分大小写的映射Map
	extName, exists := getCaseInsensitiveValue(caseInsensitiveMap, JobName)
	if exists {
		// 获取前缀名称
		preName := GetExtName(ViewName)
		// 映射完 合成jobName名称
		Name := preName + extName
		log.Printf("Name: %s", Name)
		// 判断视图名和任务名是否存在
		//views := GetJenkinsViews()
		//for view := range views {
		//	fmt.Printf("视图%s\n", views[view].Name)
		//}
		jenkinsUrl := global.GVA_CONFIG.Jenkins.Url + "view/" + ViewName + "/job/" + Name + "/buildWithParameters"
		// 判断一下 jenkins URL 是否有效
		_, err := url.ParseRequestURI(jenkinsUrl)
		if err != nil {
			global.GVA_LOG.Error("无效的jenkins URL:\n", zap.Error(err))
		}
		// 获取构建参数 params为map类型
		params := GetBuildJobParam(Name)
		if len(params) == 0 {
			JenkinsBuildJob(ViewName, Name)
			return
		}
		// 表单数据 将获取的参数转换为表单数据
		data := url.Values{}
		for key, value := range params {
			fmt.Printf("参数名称: %s, 默认值: %v\n", key, value)
			//}
			req, err := http.NewRequest("POST", jenkinsUrl, bytes.NewBufferString(data.Encode()))
			if err != nil {
				global.GVA_LOG.Error("创建post请求失败:", zap.Error(err))
			}
			//设置Auth Basic Auth
			req.SetBasicAuth(global.GVA_CONFIG.Jenkins.User, global.GVA_CONFIG.Jenkins.ApiToken)
			//发送post请求并获取响应
			client := http.Client{Timeout: 5 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				global.GVA_LOG.Error("发送post请求失败:", zap.Error(err))
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					global.GVA_LOG.Error("关闭响应体失败:", zap.Error(err))
					return
				}
			}(resp.Body)
		}
	} else {
		log.Printf("未找到构建任务名")
		//system.SendMsg()
		return // 返回不在执行后续
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

// GetJenkinsViews 获取视图名和任务名 用于判断输入JobName 是否存在
func GetJenkinsViews() []system.JenkinsView {
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
		Views []system.JenkinsView `json:"views"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		global.GVA_LOG.Error("未解析到数据:", zap.Error(err))
	}
	return data.Views
}

// GetLastBuildStatus  获取jenkins Job 任务状态
func GetLastBuildStatus(ViewName string, JobName string) (system.Build, error) {
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
		return system.Build{}, err
	}
	//打印原始响应以进行调试
	//fmt.Printf("API响应:%v", string(body))
	var job system.JenkinsJob
	if err := json.Unmarshal(body, &job); err != nil {
		return system.Build{}, fmt.Errorf("解析响应失败: %v", err)
	}

	if len(job.Builds) == 0 {
		return system.Build{}, fmt.Errorf("未找到任何构建")
	}

	// 获取到最近的构建号
	recentlyNumber := job.LastBuild.Number
	JobURL := global.GVA_CONFIG.Jenkins.Url + "view/" + ViewName + "/job/" + Name + "/" + strconv.Itoa(recentlyNumber) + "/api/json"
	body, err = GetJenkinsData(JobURL)
	if err != nil {
		return system.Build{}, err
	}
	var build system.Build
	if err := json.Unmarshal(body, &build); err != nil {
		return system.Build{}, fmt.Errorf("解析响应失败: %v", err)
	}
	// 返回最近的构建
	return build, nil
}
