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
	"T1":    "(gz)亿万-T1",
	"T2":    "(gz)狗子-T2",
	"T3":    "(gz)多多-T3",
	"T4":    "(gz)旺财-T4",
	"28国际":  "28国际",
	"鼎尚国际":  "鼎尚国际",
	"圆梦娱乐2": "圆梦娱乐2",
	"东升国际":  "东升国际",
	"大利集团":  "大利集团",
	"大满贯":   "大满贯",
	"88娱乐":  "88娱乐",
	"ABC28": "ABC28",
	"AK国际":  "AK国际",
	"CF游戏":  "CF游戏",
	"DT28":  "DT28",
	"新直属":   "新直属模版",
	"恒旺28":  "恒旺28",
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
