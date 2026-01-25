// Package model is the model package for the datasource service.
package model

import (
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameDataSourceTests = "data_source_tests"

type DataSourceTest struct {
	gormmodel.BaseModel

	DataSourceUID snowflake.ID `gorm:"column:data_source_uid;not null;index"`
	Status        int8         `gorm:"column:status;type:tinyint;not null;default:0"`
	LatencyMs     int64        `gorm:"column:latency_ms;type:bigint;not null;default:0"`
	Message       string       `gorm:"column:message;type:varchar(255);"`
}

func (DataSourceTest) TableName() string {
	return TableNameDataSourceTests
}
