package service

import (
	"context"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/internal/biz/strategy"
	"github.com/aide-family/sovereign/internal/biz/vobj"
	apiv1 "github.com/aide-family/sovereign/pkg/api/v1"
	"github.com/aide-family/sovereign/pkg/enum"
)

func NewStrategyService(strategyService *strategy.Service) *StrategyService {
	return &StrategyService{
		strategyService: strategyService,
	}
}

type StrategyService struct {
	apiv1.UnimplementedStrategyServer

	strategyService *strategy.Service
}

func (s *StrategyService) CreateStrategy(ctx context.Context, req *apiv1.CreateStrategyRequest) (*apiv1.CreateStrategyReply, error) {
	if err := s.strategyService.Create(ctx, snowflake.ParseInt64(req.NamespaceUid), vobj.StrategyType(req.Type), req.Name); err != nil {
		return nil, err
	}
	return &apiv1.CreateStrategyReply{}, nil
}

func (s *StrategyService) UpdateStrategy(ctx context.Context, req *apiv1.UpdateStrategyRequest) (*apiv1.UpdateStrategyReply, error) {
	if err := s.strategyService.Update(ctx, snowflake.ParseInt64(req.Uid), req.Name, req.Description); err != nil {
		return nil, err
	}
	return &apiv1.UpdateStrategyReply{}, nil
}

func (s *StrategyService) UpdateStrategyStatus(ctx context.Context, req *apiv1.UpdateStrategyStatusRequest) (*apiv1.UpdateStrategyStatusReply, error) {
	if err := s.strategyService.UpdateStatus(ctx, snowflake.ParseInt64(req.Uid), vobj.GlobalStatus(req.Status)); err != nil {
		return nil, err
	}
	return &apiv1.UpdateStrategyStatusReply{}, nil
}

func (s *StrategyService) DeleteStrategy(ctx context.Context, req *apiv1.DeleteStrategyRequest) (*apiv1.DeleteStrategyReply, error) {
	if err := s.strategyService.Delete(ctx, snowflake.ParseInt64(req.Uid)); err != nil {
		return nil, err
	}
	return &apiv1.DeleteStrategyReply{}, nil
}

func (s *StrategyService) GetStrategy(ctx context.Context, req *apiv1.GetStrategyRequest) (*apiv1.StrategyItem, error) {
	strategyEntity, err := s.strategyService.Get(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return toAPIV1StrategyItem(strategyEntity), nil
}

func (s *StrategyService) ListStrategy(ctx context.Context, req *apiv1.ListStrategyRequest) (*apiv1.ListStrategyReply, error) {
	query := &strategy.ListQuery{
		PageRequest:  shared.NewPageRequest(req.Page, req.PageSize),
		NamespaceUID: snowflake.ParseInt64(req.NamespaceUid),
		GroupUID:      snowflake.ParseInt64(req.GroupUid),
		Type:          vobj.StrategyType(req.Type),
		Keyword:       req.Keyword,
		Status:        vobj.GlobalStatus(req.Status),
	}
	page, err := s.strategyService.List(ctx, query)
	if err != nil {
		return nil, err
	}
	return toAPIV1ListStrategyReply(page), nil
}

func (s *StrategyService) SelectStrategy(ctx context.Context, req *apiv1.SelectStrategyRequest) (*apiv1.SelectStrategyReply, error) {
	var nextUID snowflake.ID
	if req.NextUid > 0 {
		nextUID = snowflake.ParseInt64(req.NextUid)
	}
	query := &strategy.SelectQuery{
		NamespaceUID: snowflake.ParseInt64(req.NamespaceUid),
		Keyword:      req.Keyword,
		Limit:        req.Limit,
		NextUID:      nextUID,
		Status:       vobj.GlobalStatus(req.Status),
	}
	result, err := s.strategyService.Select(ctx, query)
	if err != nil {
		return nil, err
	}
	return toAPIV1SelectStrategyReply(result), nil
}

// toAPIV1StrategyItem converts strategy entity to API response
func toAPIV1StrategyItem(s *strategy.Strategy) *apiv1.StrategyItem {
	return &apiv1.StrategyItem{
		Uid:         s.UID().Int64(),
		NamespaceUid: s.NamespaceUID().Int64(),
		GroupUid:    s.GroupUID().Int64(),
		Type:        string(s.Type()),
		Name:        s.Name(),
		Description: s.Description(),
		Status:      enum.GlobalStatus(s.Status()),
		CreatedAt:   s.CreatedAt().Unix(),
		UpdatedAt:   s.UpdatedAt().Unix(),
	}
}

// toAPIV1ListStrategyReply converts strategy page to API response
func toAPIV1ListStrategyReply(page *shared.Page[*strategy.Strategy]) *apiv1.ListStrategyReply {
	items := make([]*apiv1.StrategyItem, 0, len(page.Items))
	for _, s := range page.Items {
		items = append(items, toAPIV1StrategyItem(s))
	}
	return &apiv1.ListStrategyReply{
		Strategies: items,
		Total:      page.Total,
	}
}

// toAPIV1SelectStrategyReply converts strategy select result to API response
func toAPIV1SelectStrategyReply(result *strategy.SelectResult) *apiv1.SelectStrategyReply {
	selectItems := make([]*apiv1.SelectItem, 0, len(result.Items))
	for _, item := range result.Items {
		selectItems = append(selectItems, &apiv1.SelectItem{
			Value:    item.UID.Int64(),
			Label:    item.Name,
			Disabled: item.Disabled,
			Tooltip:  item.Tooltip,
		})
	}
	return &apiv1.SelectStrategyReply{
		Items:   selectItems,
		Total:   result.Total,
		NextUid: result.NextUID.Int64(),
		HasMore: result.HasMore,
	}
}

// 策略分组相关接口实现
func (s *StrategyService) CreateStrategyGroup(ctx context.Context, req *apiv1.CreateStrategyGroupRequest) (*apiv1.CreateStrategyGroupReply, error) {
	if err := s.strategyService.CreateGroup(ctx, snowflake.ParseInt64(req.NamespaceUid), req.Name); err != nil {
		return nil, err
	}
	return &apiv1.CreateStrategyGroupReply{}, nil
}

func (s *StrategyService) UpdateStrategyGroup(ctx context.Context, req *apiv1.UpdateStrategyGroupRequest) (*apiv1.UpdateStrategyGroupReply, error) {
	if err := s.strategyService.UpdateGroup(ctx, snowflake.ParseInt64(req.Uid), req.Name, req.Description); err != nil {
		return nil, err
	}
	return &apiv1.UpdateStrategyGroupReply{}, nil
}

func (s *StrategyService) UpdateStrategyGroupStatus(ctx context.Context, req *apiv1.UpdateStrategyGroupStatusRequest) (*apiv1.UpdateStrategyGroupStatusReply, error) {
	if err := s.strategyService.UpdateGroupStatus(ctx, snowflake.ParseInt64(req.Uid), vobj.GlobalStatus(req.Status)); err != nil {
		return nil, err
	}
	return &apiv1.UpdateStrategyGroupStatusReply{}, nil
}

func (s *StrategyService) DeleteStrategyGroup(ctx context.Context, req *apiv1.DeleteStrategyGroupRequest) (*apiv1.DeleteStrategyGroupReply, error) {
	if err := s.strategyService.DeleteGroup(ctx, snowflake.ParseInt64(req.Uid)); err != nil {
		return nil, err
	}
	return &apiv1.DeleteStrategyGroupReply{}, nil
}

func (s *StrategyService) GetStrategyGroup(ctx context.Context, req *apiv1.GetStrategyGroupRequest) (*apiv1.StrategyGroupItem, error) {
	group, err := s.strategyService.GetGroup(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return toAPIV1StrategyGroupItem(group), nil
}

func (s *StrategyService) ListStrategyGroup(ctx context.Context, req *apiv1.ListStrategyGroupRequest) (*apiv1.ListStrategyGroupReply, error) {
	query := &strategy.GroupListQuery{
		PageRequest:  shared.NewPageRequest(req.Page, req.PageSize),
		NamespaceUID: snowflake.ParseInt64(req.NamespaceUid),
		Keyword:      req.Keyword,
		Status:       vobj.GlobalStatus(req.Status),
	}
	page, err := s.strategyService.ListGroups(ctx, query)
	if err != nil {
		return nil, err
	}
	return toAPIV1ListStrategyGroupReply(page), nil
}

func (s *StrategyService) GetGroupStrategies(ctx context.Context, req *apiv1.GetGroupStrategiesRequest) (*apiv1.ListStrategyReply, error) {
	groupUID := snowflake.ParseInt64(req.GroupUid)
	query := &strategy.ListQuery{
		PageRequest: shared.NewPageRequest(req.Page, req.PageSize),
		Keyword:     req.Keyword,
		Status:      vobj.GlobalStatus(req.Status),
	}
	page, err := s.strategyService.GetGroupStrategies(ctx, groupUID, query)
	if err != nil {
		return nil, err
	}
	return toAPIV1ListStrategyReply(page), nil
}

func (s *StrategyService) AddStrategyToGroup(ctx context.Context, req *apiv1.AddStrategyToGroupRequest) (*apiv1.AddStrategyToGroupReply, error) {
	if err := s.strategyService.AddStrategyToGroup(ctx, snowflake.ParseInt64(req.StrategyUid), snowflake.ParseInt64(req.GroupUid)); err != nil {
		return nil, err
	}
	return &apiv1.AddStrategyToGroupReply{}, nil
}

func (s *StrategyService) RemoveStrategyFromGroup(ctx context.Context, req *apiv1.RemoveStrategyFromGroupRequest) (*apiv1.RemoveStrategyFromGroupReply, error) {
	if err := s.strategyService.RemoveStrategyFromGroup(ctx, snowflake.ParseInt64(req.StrategyUid)); err != nil {
		return nil, err
	}
	return &apiv1.RemoveStrategyFromGroupReply{}, nil
}

// toAPIV1StrategyGroupItem converts strategy group entity to API response
func toAPIV1StrategyGroupItem(g *strategy.StrategyGroup) *apiv1.StrategyGroupItem {
	return &apiv1.StrategyGroupItem{
		Uid:          g.UID().Int64(),
		NamespaceUid: g.NamespaceUID().Int64(),
		Name:         g.Name(),
		Description:  g.Description(),
		Status:       enum.GlobalStatus(g.Status()),
		UpgradeMode:  int32(g.UpgradeMode()),
		UpgradeConfig: g.UpgradeConfig(),
		CreatedAt:    g.CreatedAt().Unix(),
		UpdatedAt:    g.UpdatedAt().Unix(),
	}
}

// toAPIV1ListStrategyGroupReply converts strategy group page to API response
func toAPIV1ListStrategyGroupReply(page *shared.Page[*strategy.StrategyGroup]) *apiv1.ListStrategyGroupReply {
	items := make([]*apiv1.StrategyGroupItem, 0, len(page.Items))
	for _, g := range page.Items {
		items = append(items, toAPIV1StrategyGroupItem(g))
	}
	return &apiv1.ListStrategyGroupReply{
		Groups: items,
		Total:  page.Total,
	}
}

