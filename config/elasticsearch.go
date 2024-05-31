package config

import "time"

type ElasticSearch struct {
	Enable              bool          `mapstructure:"enable" json:"enable" yaml:"enable"`
	URL                 string        `mapstructure:"url" json:"url" yaml:"url"`
	Sniff               bool          `mapstructure:"sniff" json:"sniff" yaml:"sniff"`
	HealthcheckInterval time.Duration `mapstructure:"healthcheckInterval" json:"healthcheckInterval" yaml:"healthcheckInterval"`
	IndexPrefix         string        `mapstructure:"index-prefix" json:"index-prefix" yaml:"index-prefix"`
	// mapstructure 将通用的map[string]interface{} 解码到对应的结构体中
}
