package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type GambleRouter struct {
}

func (g *GambleRouter) InitGambleRouter(Router *gin.RouterGroup) {
	gambleRouter := Router.Group("gamble").Use(middleware.OperationRecord())
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		gambleRouter.GET("searchById", baseApi.GetDataByID)
		gambleRouter.GET("getAllData", baseApi.GetAllData)
		gambleRouter.POST("test", baseApi.GetData)
	}
}
