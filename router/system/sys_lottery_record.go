package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

func (g *GambleRouter) InitSysLotteryRecordRouter(Router *gin.RouterGroup) {
	gambleRouter := Router.Group("sysRecord").Use(middleware.OperationRecord())
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		gambleRouter.POST("houseWins", baseApi.HouseWins)
	}
}
