package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Config 结构体用于存储 Jenkins 的配置信息
type Config struct {
	JenkinsURL string
	Username   string
	APIToken   string
}

// SCM 结构体用于存储 SCM 配置信息
type SCM struct {
	ConfigVersion     string `xml:"configVersion"`
	UserRemoteConfigs struct {
		URLs []string `xml:"hudson.plugins.git.UserRemoteConfig>url"`
	} `xml:"userRemoteConfigs"`
	Branches []string `xml:"branches>hudson.plugins.git.BranchSpec>name"`
}

// JobConfig 结构体用于存储 Jenkins Job 的配置信息
type JobConfig struct {
	XMLName xml.Name `xml:"project"`
	SCM     SCM      `xml:"scm"`
}

// GetJobConfig 函数用于获取 Jenkins Job 的配置信息
func GetJobConfig(viewName string, jobName string) (*JobConfig, error) {
	// 构建 API URL
	apiURL := fmt.Sprintf("http://jenkins1.3333d.vip/view/%s/job/%s/config.xml", viewName, jobName)

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建 HTTP 请求失败: %v", err)
	}

	// 设置 Basic Auth 认证信息
	req.SetBasicAuth("admin", "11d2d3cd4784aa28379905bf13988ad50e")

	// 发送 HTTP 请求并获取响应
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送 HTTP 请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}
	bodyStr := strings.Replace(string(body), `<?xml version='1.1' encoding='UTF-8'?>`, `<?xml version='1.0' encoding='UTF-8'?>`, -1)
	bodyStr = strings.Replace(bodyStr, `<?xml version="1.1" encoding="UTF-8"?>`, `<?xml version="1.0" encoding="UTF-8"?>`, -1)
	//bodyStr := strings.Replace(string(body), `<?xml version='1.1' encoding='UTF-8'?>`, `<?xml version='1.0' encoding='UTF-8'?>`, 1)

	// 解析 XML 响应
	var jobConfig JobConfig
	if err := xml.Unmarshal([]byte(bodyStr), &jobConfig); err != nil {
		return nil, fmt.Errorf("解析 XML 响应失败: %v", err)
	}

	return &jobConfig, nil
}

func main() {
	jobName := "28guoji_quartz"
	viewName := "28国际"
	// 获取 Job 配置信息
	jobConfig, err := GetJobConfig(viewName, jobName)
	if err != nil {
		log.Fatalf("获取 Job 配置信息失败: %v", err)
	}

	// 打印 Git 仓库和分支信息
	if len(jobConfig.SCM.UserRemoteConfigs.URLs) > 0 && len(jobConfig.SCM.Branches) > 0 {
		fmt.Printf(viewName+jobName+"Git代码仓库URL: %s\n", jobConfig.SCM.UserRemoteConfigs.URLs[0])
		fmt.Printf(jobName+"Git分支: %s\n", jobConfig.SCM.Branches[0])
	} else {
		fmt.Println("未找到 Git 仓库或分支信息")
	}

}
