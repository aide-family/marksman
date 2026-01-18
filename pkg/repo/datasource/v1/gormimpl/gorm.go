// Package gormimpl is the implementation of the gorm repository for the datasource service.
package gormimpl

import (
	"context"
	"errors"
	"strings"

	"github.com/aide-family/magicbox/hello"
	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/safety"
	"github.com/bwmarrin/snowflake"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"

	"github.com/aide-family/sovereign/pkg/config"
	"github.com/aide-family/sovereign/pkg/connect"
	"github.com/aide-family/sovereign/pkg/enum"
	"github.com/aide-family/sovereign/pkg/merr"
	"github.com/aide-family/sovereign/pkg/repo"
	datasourcev1 "github.com/aide-family/sovereign/pkg/repo/datasource/v1"
	"github.com/aide-family/sovereign/pkg/repo/datasource/v1/gormimpl/model"
)

func init() {
	repo.RegisterDataSourceV1Factory(config.DataSourceConfig_GORM, NewGormRepository)
}

func NewGormRepository(c *config.DataSourceConfig) (datasourcev1.Repository, func() error, error) {
	ormConfig := &config.ORMConfig{}
	if pointer.IsNotNil(c) && pointer.IsNotNil(c.GetOptions()) {
		if err := anypb.UnmarshalTo(c.GetOptions(), ormConfig, proto.UnmarshalOptions{Merge: true}); err != nil {
			return nil, nil, merr.ErrorInternalServer("unmarshal orm config failed: %v", err)
		}
	}
	if ormConfig.Dialector == config.ORMConfig_TYPE_UNKNOWN {
		return nil, nil, merr.ErrorInternalServer("orm dialector is required")
	}
	db, close, err := connect.NewDB(ormConfig)
	if err != nil {
		return nil, nil, err
	}
	node, err := snowflake.NewNode(hello.NodeID())
	if err != nil {
		return nil, nil, err
	}
	return &gormRepository{repoConfig: c, db: db, node: node}, close, nil
}

type gormRepository struct {
	repoConfig *config.DataSourceConfig
	db         *gorm.DB
	node       *snowflake.Node
}

func (g *gormRepository) CreateDataSource(ctx context.Context, req *datasourcev1.CreateDataSourceRequest) (*datasourcev1.DataSourceModel, error) {
	do := &model.DataSource{
		NamespaceUID: snowflake.ParseInt64(req.NamespaceUid),
		Type:         req.Type,
		Engine:       req.Engine,
		Name:         req.Name,
		Status:       int8(req.Status),
		Endpoint:     req.Endpoint,
		Description:  req.Description,
		Config:       safety.NewMap(req.Config),
		Metadata:     safety.NewMap(req.Metadata),
	}
	do.WithCreator(1)
	do.WithUID(g.node.Generate())
	if err := g.db.WithContext(ctx).Create(do).Error; err != nil {
		return nil, merr.ErrorInternalServer("create data source failed: %v", err)
	}
	return g.GetDataSource(ctx, &datasourcev1.GetDataSourceRequest{Uid: do.UID.Int64()})
}

func (g *gormRepository) GetDataSource(ctx context.Context, req *datasourcev1.GetDataSourceRequest) (*datasourcev1.DataSourceModel, error) {
	var item model.DataSource
	if err := g.db.WithContext(ctx).Where("uid = ?", req.Uid).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, merr.ErrorNotFound("data source %d not found", req.Uid)
		}
		return nil, err
	}
	return ConvertDataSourceModel(&item), nil
}

func (g *gormRepository) UpdateDataSource(ctx context.Context, req *datasourcev1.UpdateDataSourceRequest) (*datasourcev1.ResultInfo, error) {
	updates := map[string]any{
		"name":        req.Name,
		"status":      int8(req.Status),
		"endpoint":    req.Endpoint,
		"description": req.Description,
		"config":      safety.NewMap(req.Config),
		"metadata":    safety.NewMap(req.Metadata),
	}
	result := g.db.WithContext(ctx).Model(&model.DataSource{}).Where("uid = ?", req.Uid).Updates(updates)
	if result.Error != nil {
		return nil, merr.ErrorInternalServer("update data source failed: %v", result.Error)
	}
	return convertResultInfo(result), nil
}

func (g *gormRepository) UpdateDataSourceStatus(ctx context.Context, req *datasourcev1.UpdateDataSourceStatusRequest) (*datasourcev1.ResultInfo, error) {
	result := g.db.WithContext(ctx).Model(&model.DataSource{}).Where("uid = ?", req.Uid).Update("status", int8(req.Status))
	if result.Error != nil {
		return nil, merr.ErrorInternalServer("update data source status failed: %v", result.Error)
	}
	return convertResultInfo(result), nil
}

func (g *gormRepository) DeleteDataSource(ctx context.Context, req *datasourcev1.DeleteDataSourceRequest) (*datasourcev1.ResultInfo, error) {
	result := g.db.WithContext(ctx).Where("uid = ?", req.Uid).Delete(&model.DataSource{})
	if result.Error != nil {
		return nil, merr.ErrorInternalServer("delete data source failed: %v", result.Error)
	}
	return convertResultInfo(result), nil
}

func (g *gormRepository) ListDataSource(ctx context.Context, req *datasourcev1.ListDataSourceRequest) (*datasourcev1.ListDataSourceResponse, error) {
	query := g.db.WithContext(ctx).Model(&model.DataSource{})
	query = applyFilters(query, req.Keyword, req.NamespaceUid, req.Status, req.Type, req.Engine)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, merr.ErrorInternalServer("list data source count failed: %v", err)
	}

	query = applyOrder(query, req.OrderBy, req.Order)
	if req.Page > 0 && req.PageSize > 0 {
		query = query.Limit(int(req.PageSize)).Offset(int((req.Page - 1) * req.PageSize))
	}

	items := make([]*model.DataSource, 0)
	if err := query.Find(&items).Error; err != nil {
		return nil, merr.ErrorInternalServer("list data source failed: %v", err)
	}

	models := make([]*datasourcev1.DataSourceModel, 0, len(items))
	for _, item := range items {
		models = append(models, ConvertDataSourceModel(item))
	}
	return &datasourcev1.ListDataSourceResponse{
		Items:    models,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (g *gormRepository) SelectDataSource(ctx context.Context, req *datasourcev1.SelectDataSourceRequest) (*datasourcev1.SelectDataSourceResponse, error) {
	query := g.db.WithContext(ctx).Model(&model.DataSource{})
	query = applyFilters(query, req.Keyword, req.NamespaceUid, req.Status, req.Type, req.Engine)

	order := "uid asc"
	if req.Order == datasourcev1.Order_DESC {
		order = "uid desc"
	}
	query = query.Order(order)

	if req.LastUID > 0 {
		if req.Order == datasourcev1.Order_DESC {
			query = query.Where("uid < ?", req.LastUID)
		} else {
			query = query.Where("uid > ?", req.LastUID)
		}
	}
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 20
	}
	query = query.Limit(limit)

	items := make([]*model.DataSource, 0)
	if err := query.Find(&items).Error; err != nil {
		return nil, merr.ErrorInternalServer("select data source failed: %v", err)
	}

	resultItems := make([]*datasourcev1.DataSourceItemSelect, 0, len(items))
	for _, item := range items {
		resultItems = append(resultItems, ConvertDataSourceItemSelect(item))
	}

	lastUID := int64(0)
	if len(items) > 0 {
		lastUID = items[len(items)-1].UID.Int64()
	}

	return &datasourcev1.SelectDataSourceResponse{
		Items:   resultItems,
		Total:   int64(len(resultItems)),
		LastUID: lastUID,
		HasMore: len(items) == limit,
	}, nil
}

func applyFilters(query *gorm.DB, keyword string, namespaceUID int64, status enum.GlobalStatus, dataType, engine string) *gorm.DB {
	if namespaceUID > 0 {
		query = query.Where("namespace_uid = ?", namespaceUID)
	}
	if status > enum.GlobalStatus_GlobalStatus_UNKNOWN {
		query = query.Where("status = ?", int8(status))
	}
	if strings.TrimSpace(dataType) != "" {
		query = query.Where("type = ?", dataType)
	}
	if strings.TrimSpace(engine) != "" {
		query = query.Where("engine = ?", engine)
	}
	if strings.TrimSpace(keyword) != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	return query
}

func applyOrder(query *gorm.DB, orderBy datasourcev1.Field, order datasourcev1.Order) *gorm.DB {
	column := "created_at"
	switch orderBy {
	case datasourcev1.Field_ID:
		column = "id"
	case datasourcev1.Field_UID:
		column = "uid"
	case datasourcev1.Field_NAME:
		column = "name"
	case datasourcev1.Field_STATUS:
		column = "status"
	case datasourcev1.Field_ENDPOINT:
		column = "endpoint"
	case datasourcev1.Field_CREATED_AT:
		column = "created_at"
	case datasourcev1.Field_UPDATED_AT:
		column = "updated_at"
	case datasourcev1.Field_DELETED_AT:
		column = "deleted_at"
	case datasourcev1.Field_CREATOR:
		column = "creator"
	case datasourcev1.Field_NAMESPACE_UID:
		column = "namespace_uid"
	case datasourcev1.Field_ENGINE:
		column = "engine"
	case datasourcev1.Field_TYPE:
		column = "type"
	}
	orderDir := "asc"
	if order == datasourcev1.Order_DESC {
		orderDir = "desc"
	}
	return query.Order(column + " " + orderDir)
}
