// Package model is the model package for the strategy service.
package model

import (
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameStrategyGroups = "strategy_groups"

type StrategyGroup struct {
	gormmodel.BaseModel

	NamespaceUID snowflake.ID `gorm:"column:namespace_uid;not null;index"`

	Name        string `gorm:"column:name;type:varchar(100);not null;index"`
	Description string `gorm:"column:description;type:varchar(255);"`
	Status      int8   `gorm:"column:status;type:tinyint;not null;default:0"` // 使用 vobj.GlobalStatus

	UpgradeMode   int8  `gorm:"column:upgrade_mode;type:tinyint;default:0"` // 0-不升级 1-自动升级 2-手动升级
	UpgradeConfig string `gorm:"column:upgrade_config;type:json"`
}

func (StrategyGroup) TableName() string {
	return TableNameStrategyGroups
}

