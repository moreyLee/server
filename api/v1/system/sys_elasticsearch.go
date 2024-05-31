package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

// 索引
const indexName = "lottery"

func (b *BaseApi) SearchById(c *gin.Context) {
	id, _ := c.GetQuery("id")
	res, err := global.GVA_ELASTIC.Get().Index(indexName).Type("_doc").Id(id).Do(c)
	if err != nil {
		response.FailWithMessage("查询失败", c)
	}
	response.OkWithDetailed(res.Source, "查询到的记录", c)
}
