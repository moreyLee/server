package system

type YzSiteConfig struct {
	SiteID     uint   `gorm:"primaryKey; autoIncrement comment:站点ID" json:"site_id"`
	SiteName   string `json:"site_name" gorm:"comment:站点名称"`
	Domains    string `json:"domains" gorm:"type:json comment:H5域名"`
	ApiDomains string `json:"api_domains" gorm:"type:json comment:管理后台域名"`
	BeginTime  int64  `json:"begin_time" gorm:"type int comment:开版日期"`
	CreateTime int64  `json:"create_time" gorm:"type int comment:交站时间 起租时间"`
	EndTime    int64  `json:"end_time" gorm:"type int comment:点到期时间"`
}

func (YzSiteConfig) TableName() string {
	return "yz_site_config"
}
