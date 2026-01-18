// Package model is the model package for the datasource service.
package model

import (
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameDataSourceMetadata = "data_source_metadata"

type DataSourceMetadata struct {
	gormmodel.BaseModel

	DataSourceUID snowflake.ID               `gorm:"column:data_source_uid;not null;index"`
	EntityType    string                     `gorm:"column:entity_type;type:varchar(20);not null;index"`
	Name          string                     `gorm:"column:name;type:varchar(200);not null;"`
	Kind          string                     `gorm:"column:kind;type:varchar(50);"`
	Description   string                     `gorm:"column:description;type:varchar(255);"`
	Labels        *safety.Map[string, string] `gorm:"column:labels;type:json;"`
	Extra         *safety.Map[string, string] `gorm:"column:extra;type:json;"`
}

func (DataSourceMetadata) TableName() string {
	return TableNameDataSourceMetadata
}
