package system

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

func (b *BaseApi) TestT(c *gin.Context) {
	response.Ok(c)
	fmt.Println("执行命令")
}
