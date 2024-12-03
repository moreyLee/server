package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	// 数据库连接
	db, err := sql.Open("mysql", "root:rOYkHEc#jOesowLL@tcp(47.243.51.88:3306)/cg_devops")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 确保数据库连接有效
	err = db.Ping()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	resp, err := http.Get("https://domain.3333d.vip/site_info/")
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()
	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch URL: HTTP %d", resp.StatusCode)
	}
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	//fmt.Println(string(body))  调试使用

	// 使用 goquery 加载 HTML 内容
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	// 提取 <th> 和 <td> 标签内容
	doc.Find("table tbody tr").Each(func(i int, row *goquery.Selection) {
		// 获取服务器组名称
		groupName := strings.TrimSpace(row.Find("th").Text())
		if groupName == "" {
			return
		}
		// 插入服务器组到 server_groups 表 并获取group_id
		var groupID int64
		result, err := db.Exec("INSERT INTO server_groups (group_name) VALUES (?) ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id)", groupName)
		if err != nil {
			log.Printf("Failed to insert group: %s, error: %v", groupName, err)
		}
		groupID, err = result.LastInsertId()
		if err != nil {
			log.Printf("Failed to retrieve group ID for group: %s, error: %v", groupName, err)
			return
		}
		fmt.Printf("组名: %s ,group_id: %d\n", groupName, groupID)
		// 获取 <td> 内容   站点信息插入到 servers 表
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			text := strings.TrimSpace(cell.Text())
			parts := strings.Split(text, "(")
			if len(parts) == 2 {
				siteName := parts[0]
				siteCode := strings.TrimSuffix(parts[1], ")")
				status := "active"
				if strings.Contains(text, "停服") {
					status = "inactive"
				}
				// 插入站点名称到 servers 表
				_, err := db.Exec("INSERT INTO servers (group_id, site_name, site_code, status) VALUES (?, ?, ?, ?)",
					groupID, siteName, siteCode, status)
				if err != nil {
					log.Printf("Failed to insert server: %s, error: %v", siteName, err)
				}
				fmt.Printf(" %s     %s\n", siteName, siteCode)
			} else {
				//fmt.Printf("  原始内容: %s\n", text)
			}
			//fmt.Printf(" %s\n", strings.TrimSpace(cell.Text()))
		})
	})
}
