package service

import (
	"context"
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/sovereign/internal/biz/datasource"
	"github.com/aide-family/sovereign/internal/biz/shared"
	apiv1 "github.com/aide-family/sovereign/pkg/api/v1"
	"github.com/aide-family/sovereign/pkg/enum"
)

func NewDataSourceService(dataSourceService *datasource.Service) *DataSourceService {
	return &DataSourceService{
		dataSourceService: dataSourceService,
	}
}

type DataSourceService struct {
	apiv1.UnimplementedDataSourceServer

	dataSourceService *datasource.Service
}

func (s *DataSourceService) CreateDataSource(ctx context.Context, req *apiv1.CreateDataSourceRequest) (*apiv1.CreateDataSourceReply, error) {
	if err := s.dataSourceService.Create(ctx, snowflake.ParseInt64(req.NamespaceUid), datasource.Type(req.Type), datasource.Engine(req.Engine), req.Name, req.Endpoint, req.Description, req.Config, req.Metadata); err != nil {
		return nil, err
	}
	return &apiv1.CreateDataSourceReply{}, nil
}

func (s *DataSourceService) UpdateDataSource(ctx context.Context, req *apiv1.UpdateDataSourceRequest) (*apiv1.UpdateDataSourceReply, error) {
	if err := s.dataSourceService.Update(ctx, snowflake.ParseInt64(req.Uid), req.Name, req.Endpoint, req.Description, req.Config, req.Metadata); err != nil {
		return nil, err
	}
	return &apiv1.UpdateDataSourceReply{}, nil
}

func (s *DataSourceService) UpdateDataSourceStatus(ctx context.Context, req *apiv1.UpdateDataSourceStatusRequest) (*apiv1.UpdateDataSourceStatusReply, error) {
	status := datasource.Status(req.Status)
	if err := s.dataSourceService.UpdateStatus(ctx, snowflake.ParseInt64(req.Uid), status); err != nil {
		return nil, err
	}
	return &apiv1.UpdateDataSourceStatusReply{}, nil
}

func (s *DataSourceService) DeleteDataSource(ctx context.Context, req *apiv1.DeleteDataSourceRequest) (*apiv1.DeleteDataSourceReply, error) {
	if err := s.dataSourceService.Delete(ctx, snowflake.ParseInt64(req.Uid)); err != nil {
		return nil, err
	}
	return &apiv1.DeleteDataSourceReply{}, nil
}

func (s *DataSourceService) GetDataSource(ctx context.Context, req *apiv1.GetDataSourceRequest) (*apiv1.DataSourceItem, error) {
	ds, err := s.dataSourceService.Get(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return toAPIV1DataSourceItem(ds), nil
}

func (s *DataSourceService) ListDataSource(ctx context.Context, req *apiv1.ListDataSourceRequest) (*apiv1.ListDataSourceReply, error) {
	query := &datasource.ListQuery{
		PageRequest:  shared.NewPageRequest(req.Page, req.PageSize),
		Keyword:      req.Keyword,
		Status:       datasource.Status(req.Status),
		Type:         datasource.Type(req.Type),
		Engine:       datasource.Engine(req.Engine),
		NamespaceUID: snowflake.ParseInt64(req.NamespaceUid),
	}
	page, err := s.dataSourceService.List(ctx, query)
	if err != nil {
		return nil, err
	}
	return toAPIV1ListDataSourceReply(page), nil
}

func (s *DataSourceService) SelectDataSource(ctx context.Context, req *apiv1.SelectDataSourceRequest) (*apiv1.SelectDataSourceReply, error) {
	var nextUID snowflake.ID
	if req.NextUID > 0 {
		nextUID = snowflake.ParseInt64(req.NextUID)
	}
	query := &datasource.SelectQuery{
		Keyword:      req.Keyword,
		Limit:        req.Limit,
		NextUID:      nextUID,
		Status:       datasource.Status(req.Status),
		Type:         datasource.Type(req.Type),
		Engine:       datasource.Engine(req.Engine),
		NamespaceUID: snowflake.ParseInt64(req.NamespaceUid),
	}
	result, err := s.dataSourceService.Select(ctx, query)
	if err != nil {
		return nil, err
	}
	return toAPIV1SelectDataSourceReply(result), nil
}

// toAPIV1DataSourceItem converts datasource entity to API response
func toAPIV1DataSourceItem(ds *datasource.DataSource) *apiv1.DataSourceItem {
	return &apiv1.DataSourceItem{
		Uid:          ds.UID().Int64(),
		NamespaceUid: ds.NamespaceUID().Int64(),
		Type:         ds.Type().String(),
		Engine:       ds.Engine().String(),
		Name:         ds.Name(),
		Status:       enum.GlobalStatus(ds.Status()),
		Endpoint:     ds.Endpoint(),
		Description:  ds.Description(),
		Config:       ds.Config(),
		Metadata:     ds.Metadata(),
		CreatedAt:    ds.CreatedAt().Format(time.DateTime),
		UpdatedAt:    ds.UpdatedAt().Format(time.DateTime),
	}
}

// toAPIV1ListDataSourceReply converts datasource page to API response
func toAPIV1ListDataSourceReply(page *shared.Page[*datasource.DataSource]) *apiv1.ListDataSourceReply {
	items := make([]*apiv1.DataSourceItem, 0, len(page.Items))
	for _, ds := range page.Items {
		items = append(items, toAPIV1DataSourceItem(ds))
	}
	return &apiv1.ListDataSourceReply{
		Items:    items,
		Total:    page.Total,
		Page:     page.Page,
		PageSize: page.PageSize,
	}
}

// toAPIV1SelectDataSourceReply converts datasource select result to API response
func toAPIV1SelectDataSourceReply(result *datasource.SelectResult) *apiv1.SelectDataSourceReply {
	selectItems := make([]*apiv1.DataSourceItemSelect, 0, len(result.Items))
	for _, item := range result.Items {
		selectItems = append(selectItems, &apiv1.DataSourceItemSelect{
			Value:    item.UID.Int64(),
			Label:    item.Name,
			Disabled: item.Disabled,
			Tooltip:  item.Tooltip,
		})
	}
	return &apiv1.SelectDataSourceReply{
		Items:   selectItems,
		Total:   result.Total,
		NextUID: result.NextUID.Int64(),
		HasMore: result.HasMore,
	}
}

// TestConnection tests datasource connection
// Note: After running `make api` to regenerate proto files, uncomment this method
// func (s *DataSourceService) TestConnection(ctx context.Context, req *apiv1.TestConnectionRequest) (*apiv1.TestConnectionReply, error) {
// 	err := s.dataSourceService.TestConnection(ctx, snowflake.ParseInt64(req.Uid))
// 	if err != nil {
// 		return &apiv1.TestConnectionReply{
// 			Success: false,
// 			Message: err.Error(),
// 		}, nil
// 	}
// 	return &apiv1.TestConnectionReply{
// 		Success: true,
// 		Message: "connection test successful",
// 	}, nil
// }
