package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type JKbuild struct {
	Number    int    `json:"number"`
	Result    string `json:"result"`
	Timestamp int64  `json:"timestamp"`
	URL       string `json:"url"`
}

type Job struct {
	LastBuild JKbuild   `json:"lastBuild"`
	Name      string    `json:"name"`
	Builds    []JKbuild `json:"builds"`
}

func fetchJenkinsData(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.SetBasicAuth("admin", "11d2d3cd4784aa28379905bf13988ad50e")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func getLastBuildStatus(viewName string, jobName string) (JKbuild, error) {
	jenkinsUrl := fmt.Sprintf("http://jenkins1.3333d.vip/view/%s/job/%s/api/json", viewName, jobName)
	body, err := fetchJenkinsData(jenkinsUrl)
	if err != nil {
		return JKbuild{}, err
	}
	//打印原始响应以进行调试
	//fmt.Printf("API响应:%v", string(body))
	var job Job
	if err := json.Unmarshal(body, &job); err != nil {
		return JKbuild{}, fmt.Errorf("解析响应失败: %v", err)
	}

	if len(job.Builds) == 0 {
		return JKbuild{}, fmt.Errorf("未找到任何构建")
	}

	// 获取到最近的构建号
	recentlyNumber := job.LastBuild.Number
	JobURL := "http://jenkins1.3333d.vip/view/" + viewName + "/job/" + jobName + "/" + strconv.Itoa(recentlyNumber) + "/api/json"
	body, err = fetchJenkinsData(JobURL)
	if err != nil {
		return JKbuild{}, err
	}
	var build JKbuild
	if err := json.Unmarshal(body, &build); err != nil {
		return JKbuild{}, fmt.Errorf("解析响应失败: %v", err)
	}
	// 返回最近的构建
	return build, nil
}

func main() {
	viewName := "AK国际"
	jobName := "akguoji_h5"
	build, err := getLastBuildStatus(viewName, jobName)
	if err != nil {
		fmt.Printf("获取最近的构建状态失败: %v\n", err)
		return
	}
	//fmt.Printf("所有参数:%v\n", build)
	fmt.Printf("最近的构建号: %d\n", build.Number)
	fmt.Printf("构建结果: %s\n", build.Result)
	fmt.Printf("构建URL: %s\n", build.URL)
	seconds := build.Timestamp / 1000
	nanoseconds := (build.Timestamp % 1000) * 1000000
	t := time.Unix(seconds, nanoseconds).UTC()
	location, _ := time.LoadLocation("Asia/Dubai")
	beijingTime := t.In(location)
	formattedTime := beijingTime.Format("2006-01-02 15:04:05")
	fmt.Printf("构建时间: %v\n", formattedTime)
	fmt.Printf("构建日志URL: %s\n", build.URL+"console")
}
