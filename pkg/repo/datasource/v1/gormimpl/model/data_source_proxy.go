// Package model is the model package for the datasource service.
package model

import (
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameDataSourceProxies = "data_source_proxies"

type DataSourceProxy struct {
	gormmodel.BaseModel

	DataSourceUID snowflake.ID               `gorm:"column:data_source_uid;not null;index"`
	ProxyType     string                     `gorm:"column:proxy_type;type:varchar(20);not null;index"`
	Status        int8                       `gorm:"column:status;type:tinyint;not null;default:0"`
	Endpoint      string                     `gorm:"column:endpoint;type:varchar(255);not null;"`
	Config        *safety.Map[string, string] `gorm:"column:config;type:json;"`
}

func (DataSourceProxy) TableName() string {
	return TableNameDataSourceProxies
}
