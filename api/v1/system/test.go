package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (b *BaseApi) TestT(c *gin.Context) {
	response.Ok(c)
	fmt.Println("执行命令")
}

type jsonData struct {
	Name  string `json:"name"`
	Start string `json:"jump_start"`
}

func (b *BaseApi) Domain(c *gin.Context) {
	var data jsonData
	err := c.ShouldBind(&data)
	fmt.Println("前端过来的短域名", data.Name)
	//account_id :=ce7ca80686b3787313165855f53c401e
	CfApiLogin := "djpt36@163.com"
	globalKey := "0237bd44ec3b541e622d6aa1b187aac9193f0"
	//zone_id := "f09f2f527f41da9b5f2c100c4ff61fe9"
	url := "https://api.cloudflare.com/client/v4/zones"

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
	req.Header.Set("X-Auth-Email", CfApiLogin)
	req.Header.Set("X-Auth-Key", globalKey)
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
	fmt.Println(result)
	response.OkWithDetailed(result, "域名创建完成", c)
	//c.JSON(200, gin.H{
	//	"code":       200,
	//	"name":       data.Name,
	//	"jump_start": data.Start,
	//	"result":     result,
	//})

}
