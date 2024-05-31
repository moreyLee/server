package initialize

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
)

func InitES() {
	elasticConfig := global.GVA_CONFIG.ElasticSearch
	if elasticConfig.Enable {
		fmt.Printf("elasticsearch: %v\n", elasticConfig)
		client, err := elastic.NewClient(
			elastic.SetURL(elasticConfig.URL),
			elastic.SetSniff(elasticConfig.Sniff),
			elastic.SetHealthcheckInterval(elasticConfig.HealthcheckInterval),
		)
		if err != nil {
			global.GVA_LOG.Error("创建ElasticSearch 客户端错误", zap.Error(err))
		}
		global.GVA_ELASTIC = client
	}
}
