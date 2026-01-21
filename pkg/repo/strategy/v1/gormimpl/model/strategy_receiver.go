// Package model is the model package for the strategy service.
package model

import (
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameStrategyReceivers = "strategy_receivers"

type StrategyReceiver struct {
	gormmodel.BaseModel

	StrategyUID snowflake.ID `gorm:"column:strategy_uid;not null;index"`
	ReceiverUID snowflake.ID `gorm:"column:receiver_uid;not null;index"`

	Type string `gorm:"column:type;type:varchar(20);not null"` // common/label_match
}

func (StrategyReceiver) TableName() string {
	return TableNameStrategyReceivers
}

