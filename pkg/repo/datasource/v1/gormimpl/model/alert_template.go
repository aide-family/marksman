package model

import (
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

// AlertTemplate 告警模板模型
type AlertTemplate struct {
	gormmodel.BaseModel
	DataSourceUID   snowflake.ID `gorm:"column:datasource_uid;type:bigint;not null;index"`
	Name            string       `gorm:"column:name;type:varchar(100);not null"`
	TitleTemplate   string       `gorm:"column:title_template;type:text"`
	ContentTemplate string       `gorm:"column:content_template;type:text"`
	Status          int8         `gorm:"column:status;type:tinyint;not null;default:1;comment:状态:1=启用,2=禁用"`
}

func (AlertTemplate) TableName() string {
	return "alert_template"
}
