package config

type Server struct {
	JWT     JWT     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	Mongo   Mongo   `mapstructure:"mongo" json:"mongo" yaml:"mongo"`
	Email   Email   `mapstructure:"email" json:"email" yaml:"email"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	// auto
	AutoCode Autocode `mapstructure:"autocode" json:"autocode" yaml:"autocode"`
	// gorm
	Mysql  Mysql           `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Mssql  Mssql           `mapstructure:"mssql" json:"mssql" yaml:"mssql"`
	Pgsql  Pgsql           `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	Oracle Oracle          `mapstructure:"oracle" json:"oracle" yaml:"oracle"`
	DBList []SpecializedDB `mapstructure:"db-list" json:"db-list" yaml:"db-list"`
	// oss
	Local     Local     `mapstructure:"local" json:"local" yaml:"local"`
	Qiniu     Qiniu     `mapstructure:"qiniu" json:"qiniu" yaml:"qiniu"`
	AliyunOSS AliyunOSS `mapstructure:"aliyun-oss" json:"aliyun-oss" yaml:"aliyun-oss"`
	HuaWeiObs HuaWeiObs `mapstructure:"hua-wei-obs" json:"hua-wei-obs" yaml:"hua-wei-obs"`
	AwsS3     AwsS3     `mapstructure:"aws-s3" json:"aws-s3" yaml:"aws-s3"`

	Excel Excel `mapstructure:"excel" json:"excel" yaml:"excel"`

	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`
	// ElasticSearch
	ElasticSearch ElasticSearch `mapstructure:"elasticsearch" json:"elasticsearch" yaml:"elasticsearch"`
	// TG 机器人
	Telegram Telegram `mapstructure:"telegram" json:"telegram" yaml:"telegram"`

	// CloudFlare
	Cloudflare Cloudflare `mapstructure:"cloudflare" json:"cloudflare" yaml:"cloudflare"`
	// Jenkins
	Jenkins Jenkins `mapstructure:"jenkins" json:"jenkins" yaml:"jenkins"`
	//	Operator AdminApi
	OpsLink OpAdminLink `mapstructure:"admin-link" json:"admin-link" yaml:"admin-link"`
}
