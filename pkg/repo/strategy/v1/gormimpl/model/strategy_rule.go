// Package model is the model package for the strategy service.
package model

import (
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameStrategyRules = "strategy_rules"

type StrategyRule struct {
	gormmodel.BaseModel

	StrategyUID snowflake.ID `gorm:"column:strategy_uid;not null;index"`

	RuleDetail string `gorm:"column:rule_detail;type:text"`
	Status     int8   `gorm:"column:status;type:tinyint;not null;default:0"` // 使用 vobj.GlobalStatus

	AlertLevel int8  `gorm:"column:alert_level;type:tinyint;index"` // 使用 vobj.AlertLevel
	AlertPages string `gorm:"column:alert_pages;type:json"`

	Order int32 `gorm:"column:order;type:int;default:0"`
}

func (StrategyRule) TableName() string {
	return TableNameStrategyRules
}

