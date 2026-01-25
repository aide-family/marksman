// Package model is the model package for the strategy service.
package model

import (
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameReceivers = "receivers"

type Receiver struct {
	gormmodel.BaseModel

	NamespaceUID snowflake.ID `gorm:"column:namespace_uid;not null;index"`

	Name        string `gorm:"column:name;type:varchar(100);not null;index"`
	Type        string `gorm:"column:type;type:varchar(20);not null"` // common/label_match
	Description string `gorm:"column:description;type:varchar(255);"`

	UserIDs    *safety.Map[snowflake.ID, bool] `gorm:"column:user_ids;type:json"`
	LabelMatch *safety.Map[string, string]     `gorm:"column:label_match;type:json"`

	NotifyTypes string `gorm:"column:notify_types;type:json"` // 短信/电话/邮件
}

func (Receiver) TableName() string {
	return TableNameReceivers
}

