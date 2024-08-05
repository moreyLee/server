package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lottery"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HouseWins 庄赢
func (b *BaseApi) HouseWins(c *gin.Context) {
	var sysRecord lottery.SysRecord
	err := c.ShouldBindJSON(&sysRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = gamService.HouseWins(sysRecord)
	if err != nil {
		global.GVA_LOG.Error("插入数据记录报错:", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// GetSmallData 展示小表数据
func (b *BaseApi) GetSmallData(c *gin.Context) {
	//var records []lottery.SysRecord
	//err, sysRecord := gamService.GetAllData(records)
}
