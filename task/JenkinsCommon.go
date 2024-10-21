package task

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"io"
	"net/http"
	"strings"
	"time"
)

var ExtensionMap = map[string]string{
	"后台API": "_adminapi",
	"前台API": "_api",
	"前台H5":  "_h5",
	"后台H5":  "_h5admin",
	"定时任务":  "_quartz",
	"亿万-T1": "(gz)亿万-T1",
	"狗子-T2": "(gz)狗子-T2",
	"多多-T3": "(gz)多多-T3",
	"旺财-T4": "(gz)旺财-T4",
}

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

// GetRequestJkBody  统一封装Get请求 获取Jenkins数据  返回值( 返回字节数组 []byte,err)
func GetRequestJkBody(url string, isProduction bool) ([]byte, error) {
	global.GVA_LOG.Info("Get 请求中的URL: " + url + "\n")
	var err error
	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	if isProduction {
		// 生产环境
		req.SetBasicAuth(global.GVA_CONFIG.Jenkins.User, global.GVA_CONFIG.Jenkins.ApiToken)
	} else {
		// 测试环境
		req.SetBasicAuth(global.GVA_CONFIG.Jenkins.TestUser, global.GVA_CONFIG.Jenkins.TestToken)
	}
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
