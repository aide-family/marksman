package biz

import (
	"context"

	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"

	"github.com/aide-family/sovereign/internal/biz/bo"
	"github.com/aide-family/sovereign/internal/biz/repository"
	"github.com/aide-family/sovereign/pkg/merr"
)

func NewDataSource(repo repository.DataSource, helper *klog.Helper) *DataSource {
	return &DataSource{
		repo:   repo,
		helper: klog.NewHelper(klog.With(helper.Logger(), "biz", "datasource")),
	}
}

type DataSource struct {
	helper *klog.Helper
	repo   repository.DataSource
}

func (d *DataSource) CreateDataSource(ctx context.Context, req *bo.CreateDataSourceBo) error {
	if err := d.repo.CreateDataSource(ctx, req); err != nil {
		d.helper.Errorw("msg", "create data source failed", "error", err, "name", req.Name)
		return merr.ErrorInternal("create data source %s failed", req.Name).WithCause(err)
	}
	return nil
}

func (d *DataSource) UpdateDataSource(ctx context.Context, req *bo.UpdateDataSourceBo) error {
	if err := d.repo.UpdateDataSource(ctx, req); err != nil {
		d.helper.Errorw("msg", "update data source failed", "error", err, "uid", req.UID)
		return merr.ErrorInternal("update data source %s failed", req.UID).WithCause(err)
	}
	return nil
}

func (d *DataSource) UpdateDataSourceStatus(ctx context.Context, req *bo.UpdateDataSourceStatusBo) error {
	if err := d.repo.UpdateDataSourceStatus(ctx, req); err != nil {
		d.helper.Errorw("msg", "update data source status failed", "error", err, "uid", req.UID)
		return merr.ErrorInternal("update data source status %s failed", req.UID).WithCause(err)
	}
	return nil
}

func (d *DataSource) DeleteDataSource(ctx context.Context, uid snowflake.ID) error {
	if err := d.repo.DeleteDataSource(ctx, uid); err != nil {
		d.helper.Errorw("msg", "delete data source failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete data source %s failed", uid).WithCause(err)
	}
	return nil
}

func (d *DataSource) GetDataSource(ctx context.Context, uid snowflake.ID) (*bo.DataSourceItemBo, error) {
	item, err := d.repo.GetDataSource(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("data source %s not found", uid)
		}
		d.helper.Errorw("msg", "get data source failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get data source %s failed", uid).WithCause(err)
	}
	return item, nil
}

func (d *DataSource) ListDataSource(ctx context.Context, req *bo.ListDataSourceBo) (*bo.PageResponseBo[*bo.DataSourceItemBo], error) {
	result, err := d.repo.ListDataSource(ctx, req)
	if err != nil {
		d.helper.Errorw("msg", "list data source failed", "error", err, "req", req)
		return nil, merr.ErrorInternal("list data source failed").WithCause(err)
	}
	return bo.NewPageResponseBo(result.PageRequestBo, result.GetItems()), nil
}

func (d *DataSource) SelectDataSource(ctx context.Context, req *bo.SelectDataSourceBo) (*bo.SelectDataSourceBoResult, error) {
	result, err := d.repo.SelectDataSource(ctx, req)
	if err != nil {
		d.helper.Errorw("msg", "select data source failed", "error", err, "req", req)
		return nil, merr.ErrorInternal("select data source failed").WithCause(err)
	}
	return result, nil
}
