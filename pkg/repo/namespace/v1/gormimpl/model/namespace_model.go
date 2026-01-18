// Package model is the model package for the namespace service.
package model

import (
	"errors"

	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"gorm.io/gorm"
)

const TableNameNamespaceModels = "namespace_models"

type NamespaceModel struct {
	gormmodel.BaseModel

	Namespace string `gorm:"column:namespace;type:varchar(100);not null;index"`
}

func (NamespaceModel) TableName() string {
	return TableNameNamespaceModels
}

func (n *NamespaceModel) BeforeCreate(tx *gorm.DB) (err error) {
	if err = n.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}
	if strutil.IsEmpty(n.Namespace) {
		return errors.New("namespace is required")
	}
	return nil
}

func (n *NamespaceModel) WithNamespace(namespace string) *NamespaceModel {
	n.Namespace = namespace
	return n
}
