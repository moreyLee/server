package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type QingStorRouter struct {
}

func (q *QingStorRouter) InitQingStorRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("qingstor").Use(middleware.OperationRecord())
	var baseApi = v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		userRouter.GET("dirList", baseApi.QingCloud) // 获取青云存储空间列表APK
	}

}
