// Package model is the model package for the strategy service.
package model

import (
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameStrategyGroupReceivers = "strategy_group_receivers"

type StrategyGroupReceiver struct {
	gormmodel.BaseModel

	GroupUID    snowflake.ID `gorm:"column:group_uid;not null;index"`
	ReceiverUID snowflake.ID `gorm:"column:receiver_uid;not null;index"`

	Type string `gorm:"column:type;type:varchar(20);not null"` // common/label_match
}

func (StrategyGroupReceiver) TableName() string {
	return TableNameStrategyGroupReceivers
}

