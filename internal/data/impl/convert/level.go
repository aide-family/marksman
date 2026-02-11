// Package convert provides conversion functions for level data.
package convert

import (
	"context"

	"github.com/aide-family/magicbox/contextx"
	"github.com/aide-family/magicbox/enum"

	"github.com/aide-family/marksman/internal/biz/bo"
	"github.com/aide-family/marksman/internal/data/impl/do"
)

func ToLevelItemBo(m *do.Level) *bo.LevelItemBo {
	return &bo.LevelItemBo{
		UID:       m.UID,
		Name:      m.Name,
		Remark:    m.Remark,
		Status:    m.Status,
		Metadata:  m.Metadata,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToLevelItemSelectBo(m *do.Level) *bo.LevelItemSelectBo {
	return &bo.LevelItemSelectBo{
		Value:    m.UID.Int64(),
		Label:    m.Name,
		Disabled: m.Status != enum.GlobalStatus_ENABLED || m.DeletedAt.Valid,
		Tooltip:  m.Remark,
	}
}

func ToLevelDo(ctx context.Context, req *bo.CreateLevelBo) *do.Level {
	m := &do.Level{
		Name:     req.Name,
		Remark:   req.Remark,
		Metadata: req.Metadata,
		Status:   enum.GlobalStatus_ENABLED,
	}
	m.WithCreator(contextx.GetUserUID(ctx))
	m.WithNamespace(contextx.GetNamespace(ctx))
	return m
}
