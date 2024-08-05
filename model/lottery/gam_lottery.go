package lottery

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type SysLottery struct {
	global.GVA_MODEL
	Lose            string `json:"lose" gorm:"comment:输"`
	Win             string `json:"win" gorm:"comment:赢"`
	Stake           string `json:"stake" gorm:"type:decimal(10,2) comment:本金"`
	Odds            string `json:"odds" gorm:"type:decimal(10,2) comment:赔率"`
	Commission      string `json:"commission" gorm:"type:decimal(10,2) comment:佣金"`
	ExpectedValue   string `json:"expectedValue" gorm:"type:decimal(20,10) comment:数学期望"`
	RestartIndex    string `json:"restartIndex" gorm:"comment:重起位置"`
	CashFlowIndex   string `json:"CashFlowIndex" gorm:"type:decimal(10,2) comment:现金流水位置"`
	BetCount        string `json:"betCount" gorm:"type:decimal(10,2) comment:投注额"`
	OriginalWinLose string `json:"originalWinLose" gorm:"type:decimal(10,2) comment:输赢值"`
	AdjustedWinLose string `json:"AdjustedWinLose" gorm:"type:decimal(10,2) comment:消数后的输赢值 净输赢值"`
	WinLossRecord   string `json:"winLossRecord" gorm:"comment:胜负路"`
	Result          string `json:"result" gorm:"comment:开出的结果 庄赢或闲赢"`
	Balance         string `json:"balance" gorm:"type:decimal(10,2) comment:当前账户余额"`
	Refresh         string `json:"refresh" gorm:"comment:刷新"`
}

func (SysLottery) TableName() string {
	return "sys_lottery"
}
