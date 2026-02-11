package do

import (
	"errors"

	"github.com/aide-family/magicbox/enum"
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

type Level struct {
	BaseModel
	DeletedAt    gorm.DeletedAt    `gorm:"column:deleted_at;uniqueIndex:idx__levels__namespace_uid__deleted_at__name"`
	NamespaceUID snowflake.ID      `gorm:"column:namespace_uid;default:0;uniqueIndex:idx__levels__namespace_uid__deleted_at__name"`
	Name         string            `gorm:"column:name;type:varchar(100);default:'';uniqueIndex:idx__levels__namespace_uid__deleted_at__name"`
	Remark       string            `gorm:"column:remark;type:varchar(100);default:''"`
	Metadata     map[string]string `gorm:"column:metadata;type:json;"`
	Status       enum.GlobalStatus `gorm:"column:status;type:tinyint;default:0"`
}

func (Level) TableName() string {
	return "levels"
}

func (l *Level) WithNamespace(namespace snowflake.ID) *Level {
	l.NamespaceUID = namespace
	return l
}

func (l *Level) BeforeCreate(tx *gorm.DB) (err error) {
	if l.NamespaceUID == 0 {
		return errors.New("namespace uid is required")
	}
	return nil
}
