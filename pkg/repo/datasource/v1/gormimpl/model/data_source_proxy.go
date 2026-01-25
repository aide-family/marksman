package model

import (
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameDataSourceProxy = "data_source_proxy"

// DataSourceProxy 数据源代理模型
type DataSourceProxy struct {
	gormmodel.BaseModel
	NamespaceUID  snowflake.ID               `gorm:"column:namespace_uid;type:bigint;not null;index"`
	DataSourceUID snowflake.ID               `gorm:"column:datasource_uid;type:bigint;not null;index"`
	Type          string                     `gorm:"column:type;type:varchar(50);not null"`
	Name          string                     `gorm:"column:name;type:varchar(100);not null"`
	Config        *safety.Map[string, string] `gorm:"column:config;type:json"`
}

func (DataSourceProxy) TableName() string {
	return TableNameDataSourceProxy
}
