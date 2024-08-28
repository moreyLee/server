package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

type TgService struct {
}

func (service *TgService) SelectBySiteName(config *system.YzSiteConfig) (err error) {
	//err = global.GVA_DB.Where("site_id=? OR site_name=?", config.SiteID, config.SiteName).First(&config).Error
	err = global.GVA_DB.Where("site_name=?", config.SiteName).First(&config).Error
	return
}
