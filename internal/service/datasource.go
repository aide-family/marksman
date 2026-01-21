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

func (s *DataSourceService) TestConnection(ctx context.Context, req *apiv1.TestConnectionRequest) (*apiv1.TestConnectionReply, error) {
	err := s.dataSourceService.TestConnection(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return &apiv1.TestConnectionReply{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	return &apiv1.TestConnectionReply{
		Success: true,
		Message: "connection test successful",
	}, nil
}

// 及时查询相关接口实现（占位实现，需要后续完善 biz 层）
func (s *DataSourceService) QueryDataSource(ctx context.Context, req *apiv1.QueryDataSourceRequest) (*apiv1.QueryDataSourceReply, error) {
	// TODO: 实现数据源查询逻辑
	return &apiv1.QueryDataSourceReply{
		Result: "{}",
	}, nil
}

func (s *DataSourceService) ListQueryHistory(ctx context.Context, req *apiv1.ListQueryHistoryRequest) (*apiv1.ListQueryHistoryReply, error) {
	// TODO: 实现查询历史记录逻辑
	return &apiv1.ListQueryHistoryReply{
		Items: []*apiv1.QueryHistoryItem{},
		Total:  0,
	}, nil
}

func (s *DataSourceService) SaveQueryFavorite(ctx context.Context, req *apiv1.SaveQueryFavoriteRequest) (*apiv1.SaveQueryFavoriteReply, error) {
	// TODO: 实现保存查询收藏逻辑
	return &apiv1.SaveQueryFavoriteReply{
		FavoriteId: 0,
	}, nil
}

func (s *DataSourceService) ListQueryFavorites(ctx context.Context, req *apiv1.ListQueryFavoritesRequest) (*apiv1.ListQueryFavoritesReply, error) {
	// TODO: 实现查询收藏列表逻辑
	return &apiv1.ListQueryFavoritesReply{
		Items: []*apiv1.QueryFavoriteItem{},
	}, nil
}

func (s *DataSourceService) DeleteQueryFavorite(ctx context.Context, req *apiv1.DeleteQueryFavoriteRequest) (*apiv1.DeleteQueryFavoriteReply, error) {
	// TODO: 实现删除查询收藏逻辑
	return &apiv1.DeleteQueryFavoriteReply{}, nil
}

// 告警模板相关接口实现（占位实现，需要后续完善 biz 层）
func (s *DataSourceService) CreateAlertTemplate(ctx context.Context, req *apiv1.CreateAlertTemplateRequest) (*apiv1.CreateAlertTemplateReply, error) {
	// TODO: 实现创建告警模板逻辑
	return &apiv1.CreateAlertTemplateReply{}, nil
}

func (s *DataSourceService) UpdateAlertTemplate(ctx context.Context, req *apiv1.UpdateAlertTemplateRequest) (*apiv1.UpdateAlertTemplateReply, error) {
	// TODO: 实现更新告警模板逻辑
	return &apiv1.UpdateAlertTemplateReply{}, nil
}

func (s *DataSourceService) DeleteAlertTemplate(ctx context.Context, req *apiv1.DeleteAlertTemplateRequest) (*apiv1.DeleteAlertTemplateReply, error) {
	// TODO: 实现删除告警模板逻辑
	return &apiv1.DeleteAlertTemplateReply{}, nil
}

func (s *DataSourceService) GetAlertTemplate(ctx context.Context, req *apiv1.GetAlertTemplateRequest) (*apiv1.AlertTemplateItem, error) {
	// TODO: 实现获取告警模板逻辑
	return &apiv1.AlertTemplateItem{}, nil
}

func (s *DataSourceService) ListAlertTemplate(ctx context.Context, req *apiv1.ListAlertTemplateRequest) (*apiv1.ListAlertTemplateReply, error) {
	// TODO: 实现列表查询告警模板逻辑
	return &apiv1.ListAlertTemplateReply{
		Templates: []*apiv1.AlertTemplateItem{},
		Total:     0,
	}, nil
}

func (s *DataSourceService) UpdateAlertTemplateStatus(ctx context.Context, req *apiv1.UpdateAlertTemplateStatusRequest) (*apiv1.UpdateAlertTemplateStatusReply, error) {
	// TODO: 实现更新告警模板状态逻辑
	return &apiv1.UpdateAlertTemplateStatusReply{}, nil
}

func (s *DataSourceService) ApplyAlertTemplate(ctx context.Context, req *apiv1.ApplyAlertTemplateRequest) (*apiv1.ApplyAlertTemplateReply, error) {
	// TODO: 实现应用告警模板逻辑
	return &apiv1.ApplyAlertTemplateReply{}, nil
}

// 数据源代理相关接口实现（占位实现，需要后续完善 biz 层）
func (s *DataSourceService) CreateDataSourceProxy(ctx context.Context, req *apiv1.CreateDataSourceProxyRequest) (*apiv1.CreateDataSourceProxyReply, error) {
	// TODO: 实现创建数据源代理逻辑
	return &apiv1.CreateDataSourceProxyReply{}, nil
}

func (s *DataSourceService) UpdateDataSourceProxy(ctx context.Context, req *apiv1.UpdateDataSourceProxyRequest) (*apiv1.UpdateDataSourceProxyReply, error) {
	// TODO: 实现更新数据源代理逻辑
	return &apiv1.UpdateDataSourceProxyReply{}, nil
}

func (s *DataSourceService) DeleteDataSourceProxy(ctx context.Context, req *apiv1.DeleteDataSourceProxyRequest) (*apiv1.DeleteDataSourceProxyReply, error) {
	// TODO: 实现删除数据源代理逻辑
	return &apiv1.DeleteDataSourceProxyReply{}, nil
}

func (s *DataSourceService) GetDataSourceProxy(ctx context.Context, req *apiv1.GetDataSourceProxyRequest) (*apiv1.DataSourceProxyItem, error) {
	// TODO: 实现获取数据源代理逻辑
	return &apiv1.DataSourceProxyItem{}, nil
}

func (s *DataSourceService) ListDataSourceProxy(ctx context.Context, req *apiv1.ListDataSourceProxyRequest) (*apiv1.ListDataSourceProxyReply, error) {
	// TODO: 实现列表查询数据源代理逻辑
	return &apiv1.ListDataSourceProxyReply{
		Proxies: []*apiv1.DataSourceProxyItem{},
		Total:   0,
	}, nil
}

// 元数据管理相关接口实现（占位实现，需要后续完善 biz 层）
func (s *DataSourceService) ListDataSourceMetadata(ctx context.Context, req *apiv1.ListDataSourceMetadataRequest) (*apiv1.ListDataSourceMetadataReply, error) {
	// TODO: 实现列表查询元数据逻辑
	return &apiv1.ListDataSourceMetadataReply{
		Items: []*apiv1.DataSourceMetadataItem{},
		Total: 0,
	}, nil
}

func (s *DataSourceService) GetDataSourceMetadata(ctx context.Context, req *apiv1.GetDataSourceMetadataRequest) (*apiv1.DataSourceMetadataItem, error) {
	// TODO: 实现获取元数据详情逻辑
	return &apiv1.DataSourceMetadataItem{}, nil
}

func (s *DataSourceService) RefreshDataSourceMetadata(ctx context.Context, req *apiv1.RefreshDataSourceMetadataRequest) (*apiv1.RefreshDataSourceMetadataReply, error) {
	// TODO: 实现刷新元数据逻辑
	return &apiv1.RefreshDataSourceMetadataReply{
		MetricCount: 0,
	}, nil
}
