// Package model is the model package for the subscription service.
package model

import (
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameSubscriptions = "subscriptions"

type Subscription struct {
	gormmodel.BaseModel

	UserID      snowflake.ID `gorm:"column:user_id;not null;index"`
	NamespaceUID snowflake.ID `gorm:"column:namespace_uid;not null;index"`

	Name        string `gorm:"column:name;type:varchar(100);not null;index"`
	Type        string `gorm:"column:type;type:varchar(20);not null"` // strategy_group/datasource
	Description string `gorm:"column:description;type:varchar(255);"`
	Status      int8   `gorm:"column:status;type:tinyint;not null;default:0"` // 使用 vobj.GlobalStatus

	// 订阅目标
	StrategyGroupUIDs *safety.Map[snowflake.ID, bool] `gorm:"column:strategy_group_uids;type:json"`
	DataSourceUIDs    *safety.Map[snowflake.ID, bool] `gorm:"column:datasource_uids;type:json"`

	// 策略等级筛选（全部/某一个/某几个等级）
	AlertLevels string `gorm:"column:alert_levels;type:json"` // 存储 AlertLevel 数组

	// 订阅方式（短信/电话/邮件）
	NotifyTypes string `gorm:"column:notify_types;type:json"` // 存储 NotifyType 数组

	// 元数据
	Metadata *safety.Map[string, string] `gorm:"column:metadata;type:json"`
}

func (Subscription) TableName() string {
	return TableNameSubscriptions
}

