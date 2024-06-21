package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	JenkinsURL      = "http://119.8.127.96:8001/"
	DevelopURL      = "http://192.168.217.128:8082/"
	JenkinsUser     = "admin"
	JenkinsAPIToken = "11d2d3cd4784aa28379905bf13988ad50e" //生产
	DevelopAPIToken = "11c9bc0d6ea88891f45ee4cfe5bd218287"
)

type JenkinsJob struct {
	Actions []Action `json:"actions"`
}
type Action struct {
	ParameterDefinitions []ParameterDefinitions `json:"parameterDefinitions"`
}
type ParameterDefinitions struct {
	Name                  string `json:"name"`
	DefaultParameterValue struct {
		Value interface{} `json:"value"`
	} `json:"defaultParameterValue"`
}

// JenkinsView represents a Jenkins view
type JenkinsView struct {
	Jobs []struct {
		Name string `json:"name"`
	} `json:"jobs"`
}

// GetBuildJobParam 获取JobName构建参数
func GetBuildJobParam(JobName string) map[interface{}]interface{} {
	jenkinsUrl := JenkinsURL + "job/" + JobName + "/api/json?pretty=true"
	// 定义一个map  用于保存获取的构建参数
	params := make(map[interface{}]interface{})
	req, err := http.NewRequest("GET", jenkinsUrl, nil)
	if err != nil {
		global.GVA_LOG.Error("创建post请求失败:", zap.Error(err))
	}
	// 设置Auth Basic Auth  user &token
	req.SetBasicAuth(JenkinsUser, JenkinsAPIToken)

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
	var job JenkinsJob
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

// GetExtName 获取JobName构建参数
func GetExtName(ViewName string) string {
	jenkinsUrl := JenkinsURL + "view/" + ViewName + "/api/json"
	// 定义一个map  用于保存获取的构建参数
	//params := make(map[interface{}]interface{})
	req, err := http.NewRequest("GET", jenkinsUrl, nil)
	if err != nil {
		global.GVA_LOG.Error("创建post请求失败:", zap.Error(err))
	}
	// 设置Auth Basic Auth  user &token
	req.SetBasicAuth(JenkinsUser, JenkinsAPIToken)

	fmt.Println("全局变量jenkinsUser:" + global.GVA_CONFIG.Jenkins.User)
	fmt.Println("全部变量JenkinsToken:" + global.GVA_CONFIG.Jenkins.ApiToken)
	// 创建GET请求并获取响应
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("发送GET请求失败， jenkins用户为空，token值为空 :", zap.Error(err))
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
	var jenkinsView JenkinsView
	err = json.Unmarshal(body, &jenkinsView)
	firstJob := jenkinsView.Jobs[0]
	//params[firstJob.Name] = firstJob.Name
	//fmt.Printf("%s", firstJob.Name)
	extName := strings.SplitN(firstJob.Name, "_", 2)[0]
	//fmt.Printf("后缀%s", extName)

	return extName
}

// JenkinsBuildJobWithParam 有参数构建 根据项目名构建
func JenkinsBuildJobWithParam(JobName string) {
	jenkinsUrl := global.GVA_CONFIG.Jenkins.Url + "job/" + JobName + "/buildWithParameters"
	//获取构建参数
	params := GetBuildJobParam(JobName)
	// 表单数据 将获取的参数转换为表单数据
	data := url.Values{}
	for key, value := range params {
		fmt.Printf("参数名称: %s, 默认值: %v\n", key, value)
	}
	//data.Set("PROJECT_ITEMS_SUBNAME", "0898guoji_api")
	//data.Set("PROJECT_ITEMS_SUBNAME", "0898guoji")
	//data.Set("PROJECT_ITEMS_SUBNAME", "api.l24t8e.com")
	req, err := http.NewRequest("POST", jenkinsUrl, bytes.NewBufferString(data.Encode()))
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
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			global.GVA_LOG.Error("关闭响应体失败:", zap.Error(err))
			return
		}
	}(resp.Body)
}

// JenkinsBuildJobWithView 有参数构建 根据视图名构建
func JenkinsBuildJobWithView(ViewName string, JobName string) {
	if JobName == "" || ViewName == "" {
		log.Printf("视图名和任务名不能为空")
	}
	var extensionMap = map[string]string{
		"后台API": "_adminapi",
		"前台API": "_api",
		"H5":    "_h5",
		"后台H5":  "_h5admin",
		"定时任务":  "_quartz",
	}
	log.Printf("映射名称来源: " + JobName)
	extName, exists := extensionMap[JobName]
	if exists {
		log.Printf("后缀: %s", extName)
		preName := GetExtName(ViewName)
		log.Printf("前缀: %s ", preName)
		Name := preName + extName
		log.Printf("Name: %s", Name)
		jenkinsUrl := JenkinsURL + "view/" + ViewName + "/job/" + Name + "/buildWithParameters"
		// 获取构建参数 params为map类型
		params := GetBuildJobParam(Name)
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
			req.SetBasicAuth(JenkinsUser, JenkinsAPIToken)
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
		log.Printf("请输入正确的站点信息")

		//Send("错误: 请输入正确的站点参数信息, 具体用法请参考 /help @CG88885_bot")
		return // 返回不在执行后续
	}
}

// JenkinsBuildJob  无参数构建/**
func JenkinsBuildJob(JobName string) {
	jenkinsUrl := DevelopURL + "job/" + JobName + "/build"
	req, err := http.NewRequest("POST", jenkinsUrl, nil)
	if err != nil {
		global.GVA_LOG.Error("创建post请求失败:", zap.Error(err))
		return
	}
	// 设置Auth Basic Auth
	req.SetBasicAuth(JenkinsUser, DevelopAPIToken)
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
