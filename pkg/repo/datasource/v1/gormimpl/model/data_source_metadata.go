package model

import (
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameDataSourceMetadata = "data_source_metadata"

// DataSourceMetadata 数据源元数据模型
type DataSourceMetadata struct {
	gormmodel.BaseModel
	DataSourceUID snowflake.ID `gorm:"column:datasource_uid;type:bigint;not null;index"`
	Key           string       `gorm:"column:key;type:varchar(255);not null;index"`
	Value         string       `gorm:"column:value;type:text"`
	MetadataType  string       `gorm:"column:metadata_type;type:varchar(50);index"`
	Description   string       `gorm:"column:description;type:varchar(500)"`
}

func (DataSourceMetadata) TableName() string {
	return TableNameDataSourceMetadata
}
