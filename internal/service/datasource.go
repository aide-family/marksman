package service

import (
	"context"
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/sovereign/internal/biz/datasource"
	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/internal/biz/vobj"
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

// 及时查询相关接口实现
func (s *DataSourceService) QueryDataSource(ctx context.Context, req *apiv1.QueryDataSourceRequest) (*apiv1.QueryDataSourceReply, error) {
	result, err := s.dataSourceService.QueryDataSource(ctx, snowflake.ParseInt64(req.Uid), req.Query, req.Format)
	if err != nil {
		return nil, err
	}
	return &apiv1.QueryDataSourceReply{
		Result: result,
	}, nil
}

func (s *DataSourceService) ListQueryHistory(ctx context.Context, req *apiv1.ListQueryHistoryRequest) (*apiv1.ListQueryHistoryReply, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	histories, total, err := s.dataSourceService.ListQueryHistory(ctx, snowflake.ParseInt64(req.Uid), page, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]*apiv1.QueryHistoryItem, 0, len(histories))
	for _, h := range histories {
		items = append(items, &apiv1.QueryHistoryItem{
			Id:        int64(h.ID()),
			Query:     h.Query(),
			CreatedAt: h.CreatedAt().Unix(),
		})
	}

	return &apiv1.ListQueryHistoryReply{
		Items: items,
		Total: total,
	}, nil
}

func (s *DataSourceService) SaveQueryFavorite(ctx context.Context, req *apiv1.SaveQueryFavoriteRequest) (*apiv1.SaveQueryFavoriteReply, error) {
	id, err := s.dataSourceService.SaveQueryFavorite(ctx, snowflake.ParseInt64(req.Uid), req.Name, req.Query)
	if err != nil {
		return nil, err
	}
	return &apiv1.SaveQueryFavoriteReply{
		FavoriteId: int64(id),
	}, nil
}

func (s *DataSourceService) ListQueryFavorites(ctx context.Context, req *apiv1.ListQueryFavoritesRequest) (*apiv1.ListQueryFavoritesReply, error) {
	favorites, err := s.dataSourceService.ListQueryFavorites(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}

	items := make([]*apiv1.QueryFavoriteItem, 0, len(favorites))
	for _, f := range favorites {
		items = append(items, &apiv1.QueryFavoriteItem{
			Id:        int64(f.ID()),
			Name:      f.Name(),
			Query:     f.Query(),
			CreatedAt: f.CreatedAt().Unix(),
		})
	}

	return &apiv1.ListQueryFavoritesReply{
		Items: items,
	}, nil
}

func (s *DataSourceService) DeleteQueryFavorite(ctx context.Context, req *apiv1.DeleteQueryFavoriteRequest) (*apiv1.DeleteQueryFavoriteReply, error) {
	if err := s.dataSourceService.DeleteQueryFavorite(ctx, uint32(req.FavoriteId)); err != nil {
		return nil, err
	}
	return &apiv1.DeleteQueryFavoriteReply{}, nil
}

// 告警模板相关接口实现
func (s *DataSourceService) CreateAlertTemplate(ctx context.Context, req *apiv1.CreateAlertTemplateRequest) (*apiv1.CreateAlertTemplateReply, error) {
	if err := s.dataSourceService.CreateAlertTemplate(ctx, snowflake.ParseInt64(req.DatasourceUid), req.Name, req.TitleTemplate, req.ContentTemplate); err != nil {
		return nil, err
	}
	return &apiv1.CreateAlertTemplateReply{}, nil
}

func (s *DataSourceService) UpdateAlertTemplate(ctx context.Context, req *apiv1.UpdateAlertTemplateRequest) (*apiv1.UpdateAlertTemplateReply, error) {
	if err := s.dataSourceService.UpdateAlertTemplate(ctx, snowflake.ParseInt64(req.Uid), req.Name, req.TitleTemplate, req.ContentTemplate); err != nil {
		return nil, err
	}
	return &apiv1.UpdateAlertTemplateReply{}, nil
}

func (s *DataSourceService) DeleteAlertTemplate(ctx context.Context, req *apiv1.DeleteAlertTemplateRequest) (*apiv1.DeleteAlertTemplateReply, error) {
	if err := s.dataSourceService.DeleteAlertTemplate(ctx, snowflake.ParseInt64(req.Uid)); err != nil {
		return nil, err
	}
	return &apiv1.DeleteAlertTemplateReply{}, nil
}

func (s *DataSourceService) GetAlertTemplate(ctx context.Context, req *apiv1.GetAlertTemplateRequest) (*apiv1.AlertTemplateItem, error) {
	template, err := s.dataSourceService.GetAlertTemplate(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return toAPIV1AlertTemplateItem(template), nil
}

func (s *DataSourceService) ListAlertTemplate(ctx context.Context, req *apiv1.ListAlertTemplateRequest) (*apiv1.ListAlertTemplateReply, error) {
	query := datasource.NewAlertTemplateListQuery(
		snowflake.ParseInt64(req.DatasourceUid),
		req.Page,
		req.PageSize,
		req.Keyword,
		vobj.GlobalStatus(req.Status),
	)
	page, err := s.dataSourceService.ListAlertTemplate(ctx, query)
	if err != nil {
		return nil, err
	}
	return toAPIV1ListAlertTemplateReply(page), nil
}

func (s *DataSourceService) UpdateAlertTemplateStatus(ctx context.Context, req *apiv1.UpdateAlertTemplateStatusRequest) (*apiv1.UpdateAlertTemplateStatusReply, error) {
	if err := s.dataSourceService.UpdateAlertTemplateStatus(ctx, snowflake.ParseInt64(req.Uid), vobj.GlobalStatus(req.Status)); err != nil {
		return nil, err
	}
	return &apiv1.UpdateAlertTemplateStatusReply{}, nil
}

func (s *DataSourceService) ApplyAlertTemplate(ctx context.Context, req *apiv1.ApplyAlertTemplateRequest) (*apiv1.ApplyAlertTemplateReply, error) {
	if err := s.dataSourceService.ApplyAlertTemplate(ctx, snowflake.ParseInt64(req.Uid), snowflake.ParseInt64(req.StrategyUid)); err != nil {
		return nil, err
	}
	return &apiv1.ApplyAlertTemplateReply{}, nil
}

// 辅助函数
func toAPIV1AlertTemplateItem(template *datasource.AlertTemplate) *apiv1.AlertTemplateItem {
	return &apiv1.AlertTemplateItem{
		Uid:             template.UID().Int64(),
		DatasourceUid:   template.DataSourceUID().Int64(),
		Name:            template.Name(),
		TitleTemplate:   template.TitleTemplate(),
		ContentTemplate: template.ContentTemplate(),
		Status:          enum.GlobalStatus(template.Status()),
		CreatedAt:       template.CreatedAt().Unix(),
		UpdatedAt:       template.UpdatedAt().Unix(),
	}
}

func toAPIV1ListAlertTemplateReply(page *shared.Page[*datasource.AlertTemplate]) *apiv1.ListAlertTemplateReply {
	items := make([]*apiv1.AlertTemplateItem, 0, len(page.Items))
	for _, template := range page.Items {
		items = append(items, toAPIV1AlertTemplateItem(template))
	}
	return &apiv1.ListAlertTemplateReply{
		Templates: items,
		Total:     page.Total,
	}
}

// 数据源代理相关接口实现
func (s *DataSourceService) CreateDataSourceProxy(ctx context.Context, req *apiv1.CreateDataSourceProxyRequest) (*apiv1.CreateDataSourceProxyReply, error) {
	if err := s.dataSourceService.CreateDataSourceProxy(ctx, snowflake.ParseInt64(req.NamespaceUid), snowflake.ParseInt64(req.DatasourceUid), req.Type, req.Name, req.Config); err != nil {
		return nil, err
	}
	return &apiv1.CreateDataSourceProxyReply{}, nil
}

func (s *DataSourceService) UpdateDataSourceProxy(ctx context.Context, req *apiv1.UpdateDataSourceProxyRequest) (*apiv1.UpdateDataSourceProxyReply, error) {
	if err := s.dataSourceService.UpdateDataSourceProxy(ctx, snowflake.ParseInt64(req.Uid), req.Name, req.Config); err != nil {
		return nil, err
	}
	return &apiv1.UpdateDataSourceProxyReply{}, nil
}

func (s *DataSourceService) DeleteDataSourceProxy(ctx context.Context, req *apiv1.DeleteDataSourceProxyRequest) (*apiv1.DeleteDataSourceProxyReply, error) {
	if err := s.dataSourceService.DeleteDataSourceProxy(ctx, snowflake.ParseInt64(req.Uid)); err != nil {
		return nil, err
	}
	return &apiv1.DeleteDataSourceProxyReply{}, nil
}

func (s *DataSourceService) GetDataSourceProxy(ctx context.Context, req *apiv1.GetDataSourceProxyRequest) (*apiv1.DataSourceProxyItem, error) {
	proxy, err := s.dataSourceService.GetDataSourceProxy(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return toAPIV1DataSourceProxyItem(proxy), nil
}

func (s *DataSourceService) ListDataSourceProxy(ctx context.Context, req *apiv1.ListDataSourceProxyRequest) (*apiv1.ListDataSourceProxyReply, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	proxies, total, err := s.dataSourceService.ListDataSourceProxy(ctx, snowflake.ParseInt64(req.NamespaceUid), page, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]*apiv1.DataSourceProxyItem, 0, len(proxies))
	for _, p := range proxies {
		items = append(items, toAPIV1DataSourceProxyItem(p))
	}

	return &apiv1.ListDataSourceProxyReply{
		Proxies: items,
		Total:   total,
	}, nil
}

// 辅助函数
func toAPIV1DataSourceProxyItem(proxy *datasource.DataSourceProxy) *apiv1.DataSourceProxyItem {
	return &apiv1.DataSourceProxyItem{
		Uid:           proxy.UID().Int64(),
		NamespaceUid:  proxy.NamespaceUID().Int64(),
		DatasourceUid: proxy.DataSourceUID().Int64(),
		Type:          proxy.Type(),
		Name:          proxy.Name(),
		Config:        proxy.Config(),
		CreatedAt:     proxy.CreatedAt().Unix(),
		UpdatedAt:     proxy.UpdatedAt().Unix(),
	}
}

// 元数据管理相关接口实现
func (s *DataSourceService) ListDataSourceMetadata(ctx context.Context, req *apiv1.ListDataSourceMetadataRequest) (*apiv1.ListDataSourceMetadataReply, error) {
	metadataList, err := s.dataSourceService.ListDataSourceMetadata(ctx, snowflake.ParseInt64(req.Uid), req.MetricType)
	if err != nil {
		return nil, err
	}

	items := make([]*apiv1.DataSourceMetadataItem, 0, len(metadataList))
	for _, m := range metadataList {
		// TODO: 根据实际的metadata结构转换为proto格式
		// 当前实现是简化的key-value结构，proto需要metric_name, metric_type, labels等
		items = append(items, &apiv1.DataSourceMetadataItem{
			MetricName:  m.Key(),
			MetricType:  m.MetadataType(),
			Description: m.Description(),
			LabelCount:  0,
			Labels:      make(map[string]*apiv1.MetadataLabelValues),
		})
	}

	return &apiv1.ListDataSourceMetadataReply{
		Items: items,
		Total: int64(len(items)),
	}, nil
}

func (s *DataSourceService) GetDataSourceMetadata(ctx context.Context, req *apiv1.GetDataSourceMetadataRequest) (*apiv1.DataSourceMetadataItem, error) {
	metadata, err := s.dataSourceService.GetDataSourceMetadata(ctx, snowflake.ParseInt64(req.Uid), req.MetricName)
	if err != nil {
		return nil, err
	}

	// TODO: 根据实际的metadata结构转换为proto格式
	return &apiv1.DataSourceMetadataItem{
		MetricName:  metadata.Key(),
		MetricType:  metadata.MetadataType(),
		Description: metadata.Description(),
		LabelCount:  0,
		Labels:      make(map[string]*apiv1.MetadataLabelValues),
	}, nil
}

func (s *DataSourceService) RefreshDataSourceMetadata(ctx context.Context, req *apiv1.RefreshDataSourceMetadataRequest) (*apiv1.RefreshDataSourceMetadataReply, error) {
	if err := s.dataSourceService.RefreshDataSourceMetadata(ctx, snowflake.ParseInt64(req.Uid)); err != nil {
		return nil, err
	}
	// TODO: 返回实际刷新的metric数量
	return &apiv1.RefreshDataSourceMetadataReply{
		MetricCount: 0,
	}, nil
}
