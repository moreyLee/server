package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ElasticRouter struct {
}

func (s *ElasticRouter) ElasticOpsRouter(Router *gin.RouterGroup) {
	EsRouter := Router.Group("es").Use(middleware.OperationRecord())
	var baseApi = v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		EsRouter.GET("searchById", baseApi.SearchById) //

	}

}
