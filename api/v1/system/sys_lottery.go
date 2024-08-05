package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lottery"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (b *BaseApi) GetDataByID(c *gin.Context) {
	id, _ := c.GetQuery("id")
	var sysLottery lottery.SysLottery
	err := global.GVA_DB.Where("id=?", id).First(&sysLottery).Error
	//err := gamService.LoadAllData(sysLottery)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(sysLottery, c)
}

func (b *BaseApi) GetAllData(c *gin.Context) {
	var sysLottery []lottery.SysLottery
	err := global.GVA_DB.Find(&sysLottery).Error
	//err = gamService.LoadAllData(sysLottery)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithData(sysLottery, c)
}

// GetData 测试 连通性
func (b *BaseApi) GetData(c *gin.Context) {
	response.OkWithMessage("查询成功", c)

}
