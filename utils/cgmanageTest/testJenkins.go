package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/bndr/gojenkins"
	"net/http"
)

const (
	DeveloperToken  = "11c9bc0d6ea88891f45ee4cfe5bd218287"
	ProductTokenAPI = "11d2d3cd4784aa28379905bf13988ad50e"
	JenkinsURL      = "http://119.8.127.96:8001"
	DeveloperURL    = "http://192.168.217.128:8080"
	JenkinsUser     = "admin"
	JenkinPassword  = "$C3gR01NiKWLKkg8"
	JenkinsAPIToken = "11d2d3cd4784aa28379905bf13988ad50e"
	JobName         = "0898guoji_adminapi"
	Name            = "web"
	ctx
)

func JenkinsBuildJob() error {
	jenkinsUrl := fmt.Sprintf("%s/job/%s/build", DeveloperURL, Name)
	req, err := http.NewRequest("POST", jenkinsUrl, nil)
	if err != nil {
		return err
	}
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", JenkinsUser, JenkinsAPIToken)))
	req.Header.Add("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to trigger job, status code: %d", resp.StatusCode)
	}

	return nil
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
	err := JenkinsBuildJob()
	if err != nil {
		return
	}

}
