package service

import (
	"context"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/sovereign/internal/biz"
	"github.com/aide-family/sovereign/internal/biz/bo"
	apiv1 "github.com/aide-family/sovereign/pkg/api/v1"
)

func NewDataSourceService(dataSourceBiz *biz.DataSource) *DataSourceService {
	return &DataSourceService{
		dataSourceBiz: dataSourceBiz,
	}
}

type DataSourceService struct {
	apiv1.UnimplementedDataSourceServer

	dataSourceBiz *biz.DataSource
}

func (s *DataSourceService) CreateDataSource(ctx context.Context, req *apiv1.CreateDataSourceRequest) (*apiv1.CreateDataSourceReply, error) {
	createBo := bo.NewCreateDataSourceBo(req)
	if err := s.dataSourceBiz.CreateDataSource(ctx, createBo); err != nil {
		return nil, err
	}
	return &apiv1.CreateDataSourceReply{}, nil
}

func (s *DataSourceService) UpdateDataSource(ctx context.Context, req *apiv1.UpdateDataSourceRequest) (*apiv1.UpdateDataSourceReply, error) {
	updateBo := bo.NewUpdateDataSourceBo(req)
	if err := s.dataSourceBiz.UpdateDataSource(ctx, updateBo); err != nil {
		return nil, err
	}
	return &apiv1.UpdateDataSourceReply{}, nil
}

func (s *DataSourceService) UpdateDataSourceStatus(ctx context.Context, req *apiv1.UpdateDataSourceStatusRequest) (*apiv1.UpdateDataSourceStatusReply, error) {
	updateBo := bo.NewUpdateDataSourceStatusBo(req)
	if err := s.dataSourceBiz.UpdateDataSourceStatus(ctx, updateBo); err != nil {
		return nil, err
	}
	return &apiv1.UpdateDataSourceStatusReply{}, nil
}

func (s *DataSourceService) DeleteDataSource(ctx context.Context, req *apiv1.DeleteDataSourceRequest) (*apiv1.DeleteDataSourceReply, error) {
	if err := s.dataSourceBiz.DeleteDataSource(ctx, snowflake.ParseInt64(req.Uid)); err != nil {
		return nil, err
	}
	return &apiv1.DeleteDataSourceReply{}, nil
}

func (s *DataSourceService) GetDataSource(ctx context.Context, req *apiv1.GetDataSourceRequest) (*apiv1.DataSourceItem, error) {
	itemBo, err := s.dataSourceBiz.GetDataSource(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return itemBo.ToAPIV1DataSourceItem(), nil
}

func (s *DataSourceService) ListDataSource(ctx context.Context, req *apiv1.ListDataSourceRequest) (*apiv1.ListDataSourceReply, error) {
	listBo := bo.NewListDataSourceBo(req)
	pageResponseBo, err := s.dataSourceBiz.ListDataSource(ctx, listBo)
	if err != nil {
		return nil, err
	}
	return bo.ToAPIV1ListDataSourceReply(pageResponseBo), nil
}

func (s *DataSourceService) SelectDataSource(ctx context.Context, req *apiv1.SelectDataSourceRequest) (*apiv1.SelectDataSourceReply, error) {
	selectBo := bo.NewSelectDataSourceBo(req)
	result, err := s.dataSourceBiz.SelectDataSource(ctx, selectBo)
	if err != nil {
		return nil, err
	}
	return bo.ToAPIV1SelectDataSourceReply(result), nil
}
