package system

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"os/exec"
)

func (b *BaseApi) TestT(c *gin.Context) {
	response.Ok(c)
	fmt.Println("执行命令")
}
func (b *BaseApi) Domain(c *gin.Context) {
	shortUrl := c.PostForm("shortUrl")
	//invite_url := c.PostForm("invite_url")
	//account_id :=ce7ca80686b3787313165855f53c401e
	CfApiLogin := "djpt36@163.com"
	//domain := "ss36.vip"
	//globalKey := "0237bd44ec3b541e622d6aa1b187aac9193f0"
	//zone_id := "f09f2f527f41da9b5f2c100c4ff61fe9"
	commandStr := "curl -X POST -H \"X-Auth-Key:\"0237bd44ec3b541e622d6aa1b187aac9193f0\" -H \"X-Auth-Email:\"" + CfApiLogin + "\" -H \"Content-Type: application/json\" \"https://api.cloudflare.com/client/v4/zones\" --data '{\"name\":\"" + shortUrl + "\",\"jump_start\":\"true\"}'"
	// 创建域名
	//curl -X POST -H "X-Auth-Key:0237bd44ec3b541e622d6aa1b187aac9193f0" -H "X-Auth-Email:djpt36@163.com" -H "Content-Type: application/json" "https://api.cloudflare.com/client/v4/zones" --data '{"name":"ss36.vip","jump_start":"true"}'
	//cmd := exec.Command("pwd")
	//cmd := exec.Command("curl", "-X", "POST", "-H", "X-Auth-Key:"+globalKey, "-H", "X-Auth-Email:"+CfApiLogin, "-H", "Content-Type: application/json", "https://api.cloudflare.com/client/v4/zones", " --data '{\"name\":\"ss36.vip\",\"jump_start\":\"true\"}'")
	cmd := exec.Command("/bin/bash", "-c", commandStr)
	// 执行命令 并返回输出 获取执行结果
	output, err := cmd.Output()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		fmt.Println("命令报错信息:", err.Error())
		return
	}
	// 将命令输出转换为字符串并返回
	response.OkWithMessage(string(output), c)
	fmt.Println("命令正确执行结果:", string(output))
	return
}
