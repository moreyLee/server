package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type jsonData struct {
	Name  string `json:"name"`
	Start string `json:"jump_start"`
}

func (b *BaseApi) Domain(c *gin.Context) {
	var data jsonData
	err := c.ShouldBind(&data)
	fmt.Println("前端过来的短域名", data.Name)
	// 构建请求 将struct转换为json 数据
	requestBody := new(bytes.Buffer)
	err = json.NewEncoder(requestBody).Encode(data)
	if err != nil {
		fmt.Println("json映射结构体失败")
		return
	}
	//postData := []byte(`{"name": ss36.vip,"jump_start": "true"}`)
	req, err := http.NewRequest("POST",
		global.GVA_CONFIG.Cloudflare.ApiUrl,
		bytes.NewBuffer(requestBody.Bytes()))
	if err != nil {
		fmt.Println("创建post请求失败:", err)
		return
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Email", global.GVA_CONFIG.Cloudflare.CfApiLogin)
	req.Header.Set("X-Auth-Key", global.GVA_CONFIG.Cloudflare.GlobalKey)
	// 发送post请求并获取响应
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送post请求失败:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("defer 关闭连接失败", err)
		}
	}(resp.Body)
	// 读取响应体
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("解析响应体失败:", err)
		return
	}
	// 打印post 响应的内容
	response.OkWithDetailed(result, "域名创建完成", c)

}

/**
*		修改DNS 记录
 */
type jsonDnsData struct {
	Type    string `json:"type"`
	Name    string `json:"name"`    //@ 代表域名
	Content string `json:"content"` // 记录 114.114.114.114
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}
type ZoneID struct {
	ZoneID string `uri:"zoneID" binding:"required"`
}

func (b *BaseApi) CreateDnsRecord(c *gin.Context) {
	var data jsonDnsData
	err := c.ShouldBind(&data) // 绑定json body 中的数据
	var zoneID ZoneID
	err = c.ShouldBindUri(&zoneID)                // 绑定uri 请求参数
	zoneIDStr := fmt.Sprintf("%v", zoneID.ZoneID) // 结构体转换为字符串
	fmt.Println("zoneID值", zoneID)
	fmt.Println("zoneIDStr值", zoneIDStr)
	url := global.GVA_CONFIG.Cloudflare.ApiUrl + zoneIDStr + "/dns_records"
	if zoneIDStr == "" {
		response.FailWithMessage("zoneID不能为空", c)
		return
	}

	// 构建请求 将struct转换为json 数据
	requestBody := new(bytes.Buffer)
	err = json.NewEncoder(requestBody).Encode(data)
	if err != nil {
		fmt.Println("json映射结构体失败")
		return
	}
	//postData := []byte(`{"name": ss36.vip,"jump_start": "true"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody.Bytes()))
	if err != nil {
		fmt.Println("创建post请求失败:", err)
		return
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Email", global.GVA_CONFIG.Cloudflare.CfApiLogin)
	req.Header.Set("X-Auth-Key", global.GVA_CONFIG.Cloudflare.GlobalKey)
	// 发送post请求并获取响应
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送post请求失败:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("defer 关闭连接失败", err)
		}
	}(resp.Body)
	// 读取响应体
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("解析响应体失败:", err)
		return
	}
	// 打印post 响应的内容
	response.OkWithDetailed(result, "DNS修改记录", c)
}

type Constraint struct {
	Operator string `json:"operator"`
	Value    string `json:"value"`
}
type Action struct {
	ID    string `json:"id"`
	Value struct {
		URL        string `json:"url"`
		StatusCode int    `json:"status_code"`
	} `json:"value"`
}
type pageRuleData struct {
	Targets []struct {
		Target     string     `json:"target"`
		Constraint Constraint `json:"constraint"`
	} `json:"targets"`
	Actions  []Action `json:"actions"`
	Priority int      `json:"priority"`
	Status   string   `json:"status"`
}

func (b *BaseApi) PageRule(c *gin.Context) {
	var data pageRuleData
	err := c.ShouldBind(&data) // 绑定json body 中的数据
	var zoneID ZoneID
	err = c.ShouldBindUri(&zoneID)                // 绑定uri 请求参数
	zoneIDStr := fmt.Sprintf("%v", zoneID.ZoneID) // 结构体转换为字符串
	fmt.Println("zoneID值", zoneID)
	fmt.Println("zoneIDStr值", zoneIDStr)
	url := global.GVA_CONFIG.Cloudflare.ApiUrl + zoneIDStr + "/pagerules"
	if zoneIDStr == "" {
		response.FailWithMessage("zoneID不能为空", c)
		return
	}

	// 构建请求 将struct转换为json 数据
	requestBody := new(bytes.Buffer)
	err = json.NewEncoder(requestBody).Encode(data)
	if err != nil {
		fmt.Println("json映射结构体失败")
		return
	}
	//postData := []byte(`{"name": ss36.vip,"jump_start": "true"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody.Bytes()))
	if err != nil {
		fmt.Println("创建post请求失败:", err)
		return
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Email", global.GVA_CONFIG.Cloudflare.CfApiLogin)
	req.Header.Set("X-Auth-Key", global.GVA_CONFIG.Cloudflare.GlobalKey)
	// 发送post请求并获取响应
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送post请求失败:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("defer 关闭连接失败", err)
		}
	}(resp.Body)
	// 读取响应体
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("解析响应体失败:", err)
		return
	}
	// 打印post 响应的内容
	response.OkWithDetailed(result, "添加页面规则", c)
}
