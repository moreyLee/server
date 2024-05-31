package task

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/bndr/gojenkins"
	"net/http"
)

const (
	JenkinsURL      = "http://119.8.127.96:8001/"
	JenkinsUser     = "admin"
	JenkinPassword  = "$C3gR01NiKWLKkg8"
	JenkinsAPIToken = "11d2d3cd4784aa28379905bf13988ad50e"
	JobName         = "0898guoji_adminapi"
	ctx
)

func JenkinsBuildJob() error {
	jenkinsUrl := fmt.Sprintf("%s/job/%s/build", JenkinsURL, JobName)
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
	n, err = jenkins.BuildJob(ctx, JobName, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("构建完成：", n) // n是序号

	return
}
