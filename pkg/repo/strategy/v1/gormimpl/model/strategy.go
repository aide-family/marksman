// Package model is the model package for the strategy service.
package model

import (
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameStrategies = "strategies"

type Strategy struct {
	gormmodel.BaseModel

	NamespaceUID snowflake.ID `gorm:"column:namespace_uid;not null;index"`
	GroupUID     snowflake.ID `gorm:"column:group_uid;index"` // 策略组ID，可为空

	// 基本信息
	Type        string `gorm:"column:type;type:varchar(20);not null;index"` // basic/dial_test/suppress
	Name        string `gorm:"column:name;type:varchar(100);not null;index"`
	Description string `gorm:"column:description;type:varchar(255);"`
	Status      int8   `gorm:"column:status;type:tinyint;not null;default:0"` // 使用 vobj.GlobalStatus

	// 基础策略配置
	DataSourceUIDs *safety.Map[snowflake.ID, bool] `gorm:"column:datasource_uids;type:json"`
	Query          string                          `gorm:"column:query;type:text"`
	DataSourceType string                         `gorm:"column:datasource_type;type:varchar(20)"` // metric/logs/trace/event

	// 拨测策略配置
	DialTestType    string                       `gorm:"column:dial_test_type;type:varchar(20)"` // ping/cert/port/http
	DialTestTargets *safety.Map[string, string]  `gorm:"column:dial_test_targets;type:json"`

	// 抑制策略配置
	SuppressType   string `gorm:"column:suppress_type;type:varchar(20)"`
	SuppressConfig string `gorm:"column:suppress_config;type:json"`

	// 告警配置
	AlertTitle   string `gorm:"column:alert_title;type:varchar(255)"`
	AlertContent string `gorm:"column:alert_content;type:text"`
	AlertLevel   int8   `gorm:"column:alert_level;type:tinyint;index"` // 使用 vobj.AlertLevel
	AlertPages   string `gorm:"column:alert_pages;type:json"`

	// 规则配置
	Rules string `gorm:"column:rules;type:json"`

	// 标签和元数据
	Labels   *safety.Map[string, string] `gorm:"column:labels;type:json"`
	Metadata *safety.Map[string, string] `gorm:"column:metadata;type:json"`
}

func (Strategy) TableName() string {
	return TableNameStrategies
}

