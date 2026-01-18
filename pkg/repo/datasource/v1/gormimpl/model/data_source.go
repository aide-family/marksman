// Package model is the model package for the datasource service.
package model

import (
	"github.com/aide-family/magicbox/safety"
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	"github.com/bwmarrin/snowflake"
)

const TableNameDataSources = "data_sources"

type DataSource struct {
	gormmodel.BaseModel

	NamespaceUID snowflake.ID               `gorm:"column:namespace_uid;not null;index"`
	Type         string                     `gorm:"column:type;type:varchar(20);not null;index"`
	Engine       string                     `gorm:"column:engine;type:varchar(50);not null;index"`
	Name         string                     `gorm:"column:name;type:varchar(100);not null;"`
	Status       int8                       `gorm:"column:status;type:tinyint;not null;default:0"`
	Endpoint     string                     `gorm:"column:endpoint;type:varchar(255);not null;"`
	Description  string                     `gorm:"column:description;type:varchar(255);"`
	Config       *safety.Map[string, string] `gorm:"column:config;type:json;"`
	Metadata     *safety.Map[string, string] `gorm:"column:metadata;type:json;"`
}

func (DataSource) TableName() string {
	return TableNameDataSources
}
