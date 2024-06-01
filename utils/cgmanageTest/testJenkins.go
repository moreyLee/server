package main

import (
	"context"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

const (
	DeveloperToken  = "11c9bc0d6ea88891f45ee4cfe5bd218287"
	ProductTokenAPI = "11d2d3cd4784aa28379905bf13988ad50e"
	JenkinsURL      = "http://119.8.127.96:8001"
	DeveloperURL    = "http://192.168.217.128:8082/"
	JenkinsUser     = "admin"
	JenkinPassword  = "$C3gR01NiKWLKkg8"
	JenkinsAPIToken = "11d2d3cd4784aa28379905bf13988ad50e"
	JobName         = "0898guoji_adminapi"
	Name            = "web"
	ctx
)

func JenkinsBuildJob(JobName string) {
	jenkinsUrl := DeveloperURL + "job/" + JobName + "/build"
	req, err := http.NewRequest("POST", jenkinsUrl, nil)
	if err != nil {
		fmt.Println("创建post请求失败:", err)
		return
	}
	// 设置Auth Basic Auth
	req.SetBasicAuth(JenkinsUser, DeveloperToken)
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
			global.GVA_LOG.Error("关闭resp失败:", zap.Error(err))
		}
	}(resp.Body)
	return
}

// 构建指定任务
func buildJob(ctx context.Context, jenkins *gojenkins.Jenkins, name string) (n int64) {
	var err error
	n, err = jenkins.BuildJob(ctx, name, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("构建完成：", n) // n是序号
	return
}
func main() {
	//ctx := context.Background()
	//jenkins, _ := gojenkins.CreateJenkins(nil, "http://192.168.217.128:8082", "admin", "123456").Init(ctx)
	//
	//// 构建server
	//buildJob(ctx, jenkins, "web")
	JenkinsBuildJob("web")

}
