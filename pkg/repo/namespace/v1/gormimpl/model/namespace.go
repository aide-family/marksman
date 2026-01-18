// Package model is the model package for the namespace service.
package model

import (
	"errors"

	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"gorm.io/gorm"
)

const TableNameNamespaces = "namespaces"

type Namespace struct {
	gormmodel.BaseModel

	Name     string                      `gorm:"column:name;type:varchar(100);not null;uniqueIndex"`
	Metadata *safety.Map[string, string] `gorm:"column:metadata;type:json;"`
	Status   int8                        `gorm:"column:status;type:tinyint;not null;default:0"`
}

func (Namespace) TableName() string {
	return TableNameNamespaces
}

func (n *Namespace) BeforeCreate(tx *gorm.DB) (err error) {
	if err = n.BaseModel.BeforeCreate(tx); err != nil {
		return
	}
	if n.Status <= 0 {
		return errors.New("status is required")
	}
	return
}
