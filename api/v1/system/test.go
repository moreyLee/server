package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/dto"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (b *BaseApi) TestT(c *gin.Context) {
	response.Ok(c)
	fmt.Println("执行命令")
}

func (b *BaseApi) Domain(c *gin.Context) {
	var data dto.DomainDto
	err := c.ShouldBind(&data)
	fmt.Println("前端过来的短域名", data.ShortUrl)

	//account_id :=ce7ca80686b3787313165855f53c401e
	CfApiLogin := "djpt36@163.com"
	//domain := "ss36.vip"
	globalKey := "0237bd44ec3b541e622d6aa1b187aac9193f0"
	//zone_id := "f09f2f527f41da9b5f2c100c4ff61fe9"
	//var list domainGetData
	//shortUrl := c.PostForm("short_url")
	//inviteUrl := c.PostForm("invite_url")
	url := "https://api.cloudflare.com/client/v4/zones"
	// 实例化一个结构体 将struct 转化为json
	//request := dto.DomainDto{ShortUrl: data.ShortUrl}
	// 绑定结构体解析为json格式数据
	//err :=
	if err != nil {
		response.Fail(c)
		return
	}
	// 创建post 请求
	jsonData := fmt.Sprintf("{\"name\":\"%s\",\"jump_start\":\"true\"}", data.ShortUrl)
	//postData := []byte(`{"name": ,"jump_start": "true"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
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
	//response.OkWithMessage(result)
	//response.OkWithMessage(result)
	//commandStr := "curl -X POST -H \"X-Auth-Key:\"" + globalKey + "-H \"X-Auth-Email:\"" + CfApiLogin + "\" -H \"Content-Type: application/json\" \"https://api.cloudflare.com/client/v4/zones\" --data '{\"name\":\"" + shortUrl + "\",\"jump_start\":\"true\"}'"
	// 创建域名

	//curl -X POST -H "X-Auth-Key:0237bd44ec3b541e622d6aa1b187aac9193f0" -H "X-Auth-Email:djpt36@163.com" -H "Content-Type: application/json" "https://api.cloudflare.com/client/v4/zones" --data '{"name":"ss36.vip","jump_start":"true"}'
	//cmd := exec.Command("pwd")
	//cmd := exec.Command("curl", "-X", "POST", "-H", "X-Auth-Key:"+globalKey, "-H", "X-Auth-Email:"+CfApiLogin, "-H", "Content-Type: application/json", "https://api.cloudflare.com/client/v4/zones", " --data '{\"name\":\"ss36.vip\",\"jump_start\":\"true\"}'")
	//cmd := exec.Command("/bin/bash", "-c", commandStr)
	//fmt.Println("邀请代理链接:", inviteUrl)
	//// 执行命令 并返回输出 获取执行结果
	//output, err := cmd.Output()
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	fmt.Println("命令报错信息:", err.Error())
	//	return
	//}
	//// 将命令输出转换为字符串并返回
	//response.OkWithMessage(string(output), c)
	//fmt.Println("命令正确执行结果:", string(output))
	//return
}
