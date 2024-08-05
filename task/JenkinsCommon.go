package task

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// createCaseInsensitiveMap 创建一个不区分大小写的映射
func createCaseInsensitiveMap(m map[string]string) map[string]string {
	lowerMap := make(map[string]string)
	for k, v := range m {
		lowerMap[strings.ToLower(k)] = v
	}
	return lowerMap
}

// getCaseInsensitiveValue 从不区分大小写的映射中获取值
func getCaseInsensitiveValue(m map[string]string, key string) (string, bool) {
	val, exists := m[strings.ToLower(key)]
	return val, exists
}

// GetJenkinsData  统一封装Get请求 获取Jenkins数据
func GetJenkinsData(url string) ([]byte, error) {
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
