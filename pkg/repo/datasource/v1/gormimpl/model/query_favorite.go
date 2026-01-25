package model

import (
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameQueryFavorite = "query_favorite"

// QueryFavorite 查询收藏模型
type QueryFavorite struct {
	gormmodel.BaseModel
	DataSourceUID snowflake.ID `gorm:"column:datasource_uid;type:bigint;not null;index"`
	Name          string       `gorm:"column:name;type:varchar(100);not null"`
	Query         string       `gorm:"column:query;type:text;not null"`
}

func (QueryFavorite) TableName() string {
	return TableNameQueryFavorite
}

