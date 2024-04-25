package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TestRouter struct {
}

func (s *TestRouter) TestUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("test").Use(middleware.OperationRecord())
	var baseApi = v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		userRouter.POST("testT", baseApi.TestT) // 测试api

	}
}
