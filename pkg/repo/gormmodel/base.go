// Package gormmodel provides shared GORM base models.
package gormmodel

import (
	"errors"
	"time"

	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

// BaseModel is embedded by GORM models that need standard audit fields.
type BaseModel struct {
	ID        uint32         `gorm:"column:id;primaryKey;autoIncrement"`
	UID       snowflake.ID   `gorm:"column:uid;not null;uniqueIndex"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index"`
	Creator   snowflake.ID   `gorm:"column:creator;not null;index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.Creator == 0 {
		return errors.New("creator is required")
	}

	node, err := snowflake.NewNode(hello.NodeID())
	if err != nil {
		return err
	}
	b.WithUID(node.Generate())

	return
}

func (b *BaseModel) WithCreator(creator snowflake.ID) *BaseModel {
	b.Creator = creator
	return b
}

func (b *BaseModel) WithUID(uid snowflake.ID) *BaseModel {
	b.UID = uid
	return b
}
