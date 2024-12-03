package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"time"
)

type JenkinsJob struct {
	Name string `json:"name"`
}
type JenkinsView struct {
	Name string       `json:"name"`
	Jobs []JenkinsJob `json:"jobs"`
}
type JenkinsResponse struct {
	Views []JenkinsView `json:"views"`
}

func ProdView(jenkinsUrl string, user string, token string) ([]JenkinsView, error) {
	req, err := http.NewRequest("GET", jenkinsUrl, nil)
	req.SetBasicAuth(user, token)
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	//resp, err := http.Get(jenkinsURL)
	if err != nil {
		log.Fatalf("Failed to get Jenkins views: %v", err)
	}
	defer resp.Body.Close()
	// 检查响应的状态码是否为 200
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // 读取响应内容
		return nil, fmt.Errorf("请求失败状态码:  %d, response body: %s", resp.StatusCode, body)
	}
	// 读取并解析 JSON 数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	var viewNames struct {
		Views []JenkinsView `json:"views"`
	}
	err = json.Unmarshal(body, &viewNames)
	if err != nil {
		fmt.Printf("JSON body 解析失败: %v", err)
	}
	return viewNames.Views, nil
}

// FetchJenkinsData 解析 Jenkins API 返回的数据
func FetchJenkinsData(url, user, token string) ([]JenkinsView, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.SetBasicAuth(user, token)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var data JenkinsResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return data.Views, nil
}

// InsertDataToDB 测试环境 将视图和任务名 插入到数据库
func InsertDataToDB(db *sql.DB, views []JenkinsView, tableName string) error {
	_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName))
	if err != nil {
		return fmt.Errorf("failed to truncate table %s: %v", tableName, err)
	}
	fmt.Printf("表 %s 已成功清空 truncated successfully！！！\n", tableName)
	// 插入视图名和任务名
	insertSQL := fmt.Sprintf("INSERT INTO %s (site_name, task_name) VALUES (?, ?)", tableName)
	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()
	//fmt.Printf("Retrieved views: %+v\n", views)
	for _, view := range views {
		for _, job := range view.Jobs {
			_, err := stmt.Exec(view.Name, job.Name)
			if err != nil {
				return fmt.Errorf("failed to insert data: %v", err)
			}
			fmt.Printf("jenins_env_test Inserted view_name : %s, job: %s\n", view.Name, job.Name)
		}
	}
	return nil
}

func InsertViewsIntoDB(db *sql.DB, views []JenkinsView, tableName string, columnName string) error {
	_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName))
	if err != nil {
		return fmt.Errorf("failed to truncate table %s: %v", tableName, err)
	}
	fmt.Printf("表已成功清空 %s truncated successfully！！！\n", tableName)
	// 视图名插入到数据库
	insertSql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (?)", tableName, columnName)
	for _, view := range views {
		_, err := db.Exec(insertSql, view.Name)
		if err != nil {
			return fmt.Errorf("failed to insert view name into %s: %v", tableName, err)
		}
		fmt.Printf("Inserted view name into %s: %s\n", tableName, view.Name)
	}
	return nil
}
func main() {
	jenkinsURL := "http://jenkins1.3333d.vip/api/json?tree=views[name]"
	user := "admin"
	tokenApi := "115202b0a72dadd4f89878e7d352aa8552"
	jkTestURL := "https://jenkins.qiyinyun.com/api/json?tree=views[name,jobs[name]]"
	TestUser := "root"
	TestToken := "11700ee17be3621da8bb4443e073763a69"
	prodViews, _ := ProdView(jenkinsURL, user, tokenApi)
	views, err := FetchJenkinsData(jkTestURL, TestUser, TestToken)
	// 数据库连接
	//db, err := sql.Open("mysql", "root:Devops%588@tcp(localhost:3306)/cg_devops")
	// 测试环境
	db, err := sql.Open("mysql", "root:rOYkHEc#jOesowLL@tcp(47.243.51.88:3306)/cg_devops")

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	// 插入生产视图名到数据库
	err = InsertViewsIntoDB(db, prodViews, "jenkins_env_prod", "prod_site_name")
	if err != nil {
		log.Fatalf("Error inserting production views into database: %v", err)
	}
	// 插入测试视图名到数据库

	err = InsertDataToDB(db, views, "jenkins_env_test")
	if err != nil {
		log.Fatalf("Error inserting test views into database: %v", err)
	}
}
