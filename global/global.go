package global

import (
	"github.com/olivere/elastic"
	"github.com/qiniu/qmgo"
	"sync"

	"github.com/flipped-aurora/gin-vue-admin/server/utils/timer"
	"github.com/songzhibin97/gkit/cache/local_cache"

	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	"github.com/flipped-aurora/gin-vue-admin/server/config"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	GVA_DB     *gorm.DB
	GVA_DBList map[string]*gorm.DB
	GVA_REDIS  redis.UniversalClient
	GVA_MONGO  *qmgo.QmgoClient
	GVA_CONFIG config.Server
	GVA_VP     *viper.Viper
	// GVA_LOG    *oplogging.Logger
	GVA_LOG                 *zap.Logger
	GVA_Timer               timer.Timer = timer.NewTimerTask()
	GVA_Concurrency_Control             = &singleflight.Group{}

	BlackCache  local_cache.Cache
	lock        sync.RWMutex
	GVA_ELASTIC *elastic.Client
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return GVA_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := GVA_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}

// globalMap 定义全局的映射关系 jenkins job 后缀对应的项目job名称后缀
var globalMap = map[string]string{
	"后台API":  "_adminapi",
	"前台API":  "_api",
	"前台H5":   "_h5",
	"后台H5":   "_h5admin",
	"定时任务":   "_quartz",
	"重启报表":   "/etc/init.d/admin",
	"重启API":  "/etc/init.d/api",
	"重启定时任务": "/etc/init.d/quartz",
	"重启游戏拉单": "/etc/init.d/thirdOrder",
	"重启机器人":  "systemctl status robot",
}
