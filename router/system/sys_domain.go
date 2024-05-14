package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DomainRouter struct {
}

func (s *DomainRouter) DomainRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("sys").Use(middleware.OperationRecord())
	var baseApi = v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		userRouter.POST("domain", baseApi.Domain) // CF 创建域名
		userRouter.POST("test1", baseApi.Test1)
	}

}
