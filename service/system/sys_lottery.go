package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/lottery"
)

type GamService struct {
}

func (gamService *GamService) LoadAllData(l lottery.SysLottery) (err error) {
	var data lottery.SysLottery
	err = global.GVA_DB.Where("id=?", l.ID).First(&data).Error
	return
}

func (gamService *GamService) HouseWins(sysRecord lottery.SysRecord) (err error) {
	err = global.GVA_DB.Create(&sysRecord).Error
	return err
}

// GetAllData 获取本金所有记录
func (gamService *GamService) GetAllData(records lottery.SysLottery) (err error) {
	//result := global.GVA_DB.Find(&records)
	return
}
