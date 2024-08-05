package lottery

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type SysRecord struct {
	global.GVA_MODEL
	Stake      string `json:"stake" gorm:"type:decimal(10,2) comment:本金"`
	NetCapital string `json:"net_capital" gorm:"type:decimal(10,2) comment:净本金"`
}

func (SysRecord) TableName() string {
	return "sys_lottery_record"
}
