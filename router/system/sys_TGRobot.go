package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TelegramRouter struct {
}

func (s *TelegramRouter) TGRobotRouter(Router *gin.RouterGroup) {
	TGRouter := Router.Group("telegram").Use(middleware.OperationRecord())
	var baseApi = v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		TGRouter.POST("jenkins", baseApi.Domain) // CF 创建域名

	}

}
