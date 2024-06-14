package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TelegramRouter struct {
}

func (s *TelegramRouter) InitTelegramRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("telegram").Use(middleware.OperationRecord())
	var baseApi = v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		userRouter.POST("webhook", baseApi.TelegramWebhook) //webhook 发送消息

	}

}
