package convert

import (
	"context"

	"github.com/aide-family/magicbox/contextx"
	"github.com/aide-family/magicbox/enum"

	"github.com/aide-family/marksman/internal/biz/bo"
	"github.com/aide-family/marksman/internal/data/impl/do"
)

func ToDatasourceItemBo(m *do.Datasource) *bo.DatasourceItemBo {
	if m == nil {
		return nil
	}
	return &bo.DatasourceItemBo{
		UID:       m.UID,
		Name:      m.Name,
		Type:      m.Type,
		Driver:    m.Driver,
		Metadata:  m.Metadata,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToDatasourceDo(ctx context.Context, req *bo.CreateDatasourceBo) *do.Datasource {
	m := &do.Datasource{
		Name:     req.Name,
		Type:     req.Type,
		Driver:   req.Driver,
		Metadata: req.Metadata,
		Status:   enum.GlobalStatus_ENABLED,
	}
	m.WithCreator(contextx.GetUserUID(ctx))
	m.WithNamespace(contextx.GetNamespace(ctx))
	return m
}
