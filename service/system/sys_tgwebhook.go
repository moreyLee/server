package system

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
	"time"
)

type TgService struct {
}

func (service *TgService) SelectBySiteName(config *system.YzSiteConfig) (err error) {
	//err = global.GVA_DB.Where("site_id=? OR site_name=?", config.SiteID, config.SiteName).First(&config).Error
	err = global.GVA_DB.Where("site_name=?", config.SiteName).First(&config).Error
	return
}
func (service *TgService) CheckJenView(viewName string, tableName string, columnName string) (bool, error) {
	query := fmt.Sprintf("%s = ?", columnName)
	err := global.GVA_DB.Table(tableName).Where(query, viewName).First(&system.JenkinsBuild{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
		return false, err
	}
	return true, err
}
func (service *TgService) SaveAdminLoginToken(token string) (err error) {
	adminLoginToken := &system.AdminLoginToken{
		ID:        1,
		HttpToken: token,
		CreatedAt: time.Now(),
	}
	err = global.GVA_DB.Save(&adminLoginToken).Error
	return
}

func (service *TgService) GetAdminLoginToken() (adminLoginToken system.AdminLoginToken, err error) {
	//var adminLoginToken system.AdminLoginToken
	err = global.GVA_DB.Model(&system.AdminLoginToken{}).Where("id = ?", 1).First(&adminLoginToken).Error
	if err != nil {
		fmt.Println("token 报错信息", err)
	}
	return adminLoginToken, err
}
