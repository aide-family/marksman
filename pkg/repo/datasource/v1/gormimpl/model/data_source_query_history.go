// Package model is the model package for the datasource service.
package model

import (
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameDataSourceQueryHistory = "data_source_query_history"

type DataSourceQueryHistory struct {
	gormmodel.BaseModel

	DataSourceUID snowflake.ID `gorm:"column:data_source_uid;not null;index"`
	Query         string       `gorm:"column:query;type:text;not null;"`
	Favorite      bool         `gorm:"column:favorite;type:tinyint;not null;default:0"`
}

func (DataSourceQueryHistory) TableName() string {
	return TableNameDataSourceQueryHistory
}
