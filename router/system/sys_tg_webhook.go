package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type TelegramRouter struct {
}

func (s *TelegramRouter) InitTelegramRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	userRouter := Router.Group("jenkins")
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		userRouter.POST("telegram-webhook", baseApi.TelegramWebhook) //webhook 发送消息

	}
	return userRouter

}
