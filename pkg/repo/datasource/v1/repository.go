// Package datasourcev1 is the datasource repository service.
package datasourcev1

import (
	context "context"
)

type Repository interface {
	CreateDataSource(ctx context.Context, req *CreateDataSourceRequest) (*DataSourceModel, error)
	GetDataSource(ctx context.Context, req *GetDataSourceRequest) (*DataSourceModel, error)
	UpdateDataSource(ctx context.Context, req *UpdateDataSourceRequest) (*ResultInfo, error)
	UpdateDataSourceStatus(ctx context.Context, req *UpdateDataSourceStatusRequest) (*ResultInfo, error)
	DeleteDataSource(ctx context.Context, req *DeleteDataSourceRequest) (*ResultInfo, error)
	ListDataSource(ctx context.Context, req *ListDataSourceRequest) (*ListDataSourceResponse, error)
	SelectDataSource(ctx context.Context, req *SelectDataSourceRequest) (*SelectDataSourceResponse, error)
}
