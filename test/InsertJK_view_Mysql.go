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

type JenkinsView struct {
	Name string `json:"name"`
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
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	return viewNames.Views, nil
}
func TestView(url string, user string, token string) ([]JenkinsView, error) {
	// 获取 Jenkins 视图名
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(user, token)
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	//resp, err := http.Get(jenkinsURL)
	if err != nil {
		log.Fatalf("Failed to get Jenkins views: %v", err)
	}
	defer resp.Body.Close()

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
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	return viewNames.Views, nil
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
	tokenApi := "11d2d3cd4784aa28379905bf13988ad50e"
	jkTestURL := "https://jenkins.qiyinyun.com/api/json?tree=views[name]"
	TestUser := "root"
	TestToken := "117a9f29e2793cb262426c8fbbb39b27cd"
	prodViews, _ := ProdView(jenkinsURL, user, tokenApi)

	testViews, _ := TestView(jkTestURL, TestUser, TestToken)
	// 数据库连接
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
	err = InsertViewsIntoDB(db, testViews, "jenkins_env_test", "test_site_name")
	if err != nil {
		log.Fatalf("Error inserting test views into database: %v", err)
	}
}
