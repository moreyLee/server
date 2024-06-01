package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"time"
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

const (
	JenkinsURL      = "http://119.8.127.96:8001/"
	JenkinsUrl      = "http://192.168.217.128:8082/"
	JenkinsUser     = "admin"
	JenkinPassword  = "$C3gR01NiKWLKkg8"
	JenkinsAPIToken = "11d2d3cd4784aa28379905bf13988ad50e"
	DevelopAPIToken = "11c9bc0d6ea88891f45ee4cfe5bd218287"
)

func JenkinsCommon() {

}

// JenkinsBuildJob  无参数构建/**
func JenkinsBuildJob(JobName string) {
	jenkinsUrl := JenkinsUrl + "job/" + JobName + "/build"
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

// JenkinsBuildJobWithParam 有参数构建
func JenkinsBuildJobWithParam(JobName string) {
	jenkinsUrl := JenkinsUrl + "job/" + JobName + "/buildWithParameters"
	GetBuildJobParam(JobName)
	// 表单数据
	data := url.Values{}
	data.Set("key1", "web")
	data.Set("key2", "test")
	data.Set("key3", "www.test.com")
	req, err := http.NewRequest("POST", jenkinsUrl, bytes.NewBufferString(data.Encode()))
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
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			global.GVA_LOG.Error("关闭响应体失败:", zap.Error(err))
			return
		}
	}(resp.Body)
}

// GetBuildJobParam 获取JobName构建参数 @param GET
func GetBuildJobParam(JobName string) {
	jenkinsUrl := JenkinsUrl + "job/" + JobName + "/api/json?pretty=true"
	req, err := http.NewRequest("GET", jenkinsUrl, nil)
	if err != nil {
		global.GVA_LOG.Error("创建post请求失败:", zap.Error(err))
		return
	}
	// 设置Auth Basic Auth
	req.SetBasicAuth(JenkinsUser, DevelopAPIToken)
	// 创建GET请求并获取响应
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		global.GVA_LOG.Error("发送GET请求失败:", zap.Error(err))
		return
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	//fmt.Printf(string(body))
	// 解析json 响应
	var job JenkinsJob
	err = json.Unmarshal(body, &job)
	// 输出构建参数的名称和值
	for _, action := range job.Actions {
		for _, param := range action.ParameterDefinitions {
			fmt.Printf("参数名称:%s, 默认值%v\n", param.Name, param.DefaultParameterValue.Value)
		}
	}
}
