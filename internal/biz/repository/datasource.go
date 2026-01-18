package repository

import (
	"context"

	"github.com/aide-family/sovereign/internal/biz/bo"
	"github.com/bwmarrin/snowflake"
)

type DataSource interface {
	CreateDataSource(ctx context.Context, req *bo.CreateDataSourceBo) error
	UpdateDataSource(ctx context.Context, req *bo.UpdateDataSourceBo) error
	UpdateDataSourceStatus(ctx context.Context, req *bo.UpdateDataSourceStatusBo) error
	DeleteDataSource(ctx context.Context, uid snowflake.ID) error
	GetDataSource(ctx context.Context, uid snowflake.ID) (*bo.DataSourceItemBo, error)
	ListDataSource(ctx context.Context, req *bo.ListDataSourceBo) (*bo.PageResponseBo[*bo.DataSourceItemBo], error)
	SelectDataSource(ctx context.Context, req *bo.SelectDataSourceBo) (*bo.SelectDataSourceBoResult, error)
}
