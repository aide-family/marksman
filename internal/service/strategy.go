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
			Value:    item.UID.String(),
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

// 策略组升级模式更新接口实现
func (s *StrategyService) UpdateStrategyGroupUpgradeMode(ctx context.Context, req *apiv1.UpdateStrategyGroupUpgradeModeRequest) (*apiv1.UpdateStrategyGroupUpgradeModeReply, error) {
	if err := s.strategyService.UpdateGroupUpgradeMode(ctx, snowflake.ParseInt64(req.Uid), vobj.UpgradeMode(req.UpgradeMode)); err != nil {
		return nil, err
	}
	return &apiv1.UpdateStrategyGroupUpgradeModeReply{}, nil
}

func (s *StrategyService) UpdateStrategyGroupUpgradeConfig(ctx context.Context, req *apiv1.UpdateStrategyGroupUpgradeConfigRequest) (*apiv1.UpdateStrategyGroupUpgradeConfigReply, error) {
	if err := s.strategyService.UpdateGroupUpgradeConfig(ctx, snowflake.ParseInt64(req.Uid), req.UpgradeConfig); err != nil {
		return nil, err
	}
	return &apiv1.UpdateStrategyGroupUpgradeConfigReply{}, nil
}

// 接收对象相关接口实现
func (s *StrategyService) CreateReceiver(ctx context.Context, req *apiv1.CreateReceiverRequest) (*apiv1.CreateReceiverReply, error) {
	userIDs := make(map[snowflake.ID]bool)
	for _, uid := range req.UserIds {
		userIDs[snowflake.ParseInt64(uid)] = true
	}
	notifyTypes := make([]vobj.NotifyType, 0, len(req.NotifyTypes))
	for _, t := range req.NotifyTypes {
		notifyTypes = append(notifyTypes, vobj.NotifyType(t))
	}
	if err := s.strategyService.CreateReceiver(ctx, snowflake.ParseInt64(req.NamespaceUid), vobj.ReceiverType(req.Type), req.Name, req.Description, userIDs, req.LabelMatch, notifyTypes); err != nil {
		return nil, err
	}
	return &apiv1.CreateReceiverReply{}, nil
}

func (s *StrategyService) UpdateReceiver(ctx context.Context, req *apiv1.UpdateReceiverRequest) (*apiv1.UpdateReceiverReply, error) {
	userIDs := make(map[snowflake.ID]bool)
	for _, uid := range req.UserIds {
		userIDs[snowflake.ParseInt64(uid)] = true
	}
	notifyTypes := make([]vobj.NotifyType, 0, len(req.NotifyTypes))
	for _, t := range req.NotifyTypes {
		notifyTypes = append(notifyTypes, vobj.NotifyType(t))
	}
	if err := s.strategyService.UpdateReceiver(ctx, snowflake.ParseInt64(req.Uid), req.Name, req.Description, userIDs, req.LabelMatch, notifyTypes); err != nil {
		return nil, err
	}
	return &apiv1.UpdateReceiverReply{}, nil
}

func (s *StrategyService) DeleteReceiver(ctx context.Context, req *apiv1.DeleteReceiverRequest) (*apiv1.DeleteReceiverReply, error) {
	if err := s.strategyService.DeleteReceiver(ctx, snowflake.ParseInt64(req.Uid)); err != nil {
		return nil, err
	}
	return &apiv1.DeleteReceiverReply{}, nil
}

func (s *StrategyService) GetReceiver(ctx context.Context, req *apiv1.GetReceiverRequest) (*apiv1.ReceiverItem, error) {
	receiver, err := s.strategyService.GetReceiver(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return toAPIV1ReceiverItem(receiver), nil
}

func (s *StrategyService) ListReceiver(ctx context.Context, req *apiv1.ListReceiverRequest) (*apiv1.ListReceiverReply, error) {
	query := &strategy.ReceiverListQuery{
		PageRequest:  shared.NewPageRequest(req.Page, req.PageSize),
		NamespaceUID: snowflake.ParseInt64(req.NamespaceUid),
		Type:         vobj.ReceiverType(req.Type),
		Keyword:      req.Keyword,
	}
	page, err := s.strategyService.ListReceivers(ctx, query)
	if err != nil {
		return nil, err
	}
	return toAPIV1ListReceiverReply(page), nil
}

// 策略规则相关接口实现
func (s *StrategyService) CreateStrategyRule(ctx context.Context, req *apiv1.CreateStrategyRuleRequest) (*apiv1.CreateStrategyRuleReply, error) {
	alertPages := make([]string, 0, len(req.AlertPages))
	alertPages = append(alertPages, req.AlertPages...)
	if err := s.strategyService.CreateStrategyRule(ctx, snowflake.ParseInt64(req.StrategyUid), req.RuleDetail, vobj.AlertLevel(req.AlertLevel), alertPages, req.Order); err != nil {
		return nil, err
	}
	return &apiv1.CreateStrategyRuleReply{}, nil
}

func (s *StrategyService) UpdateStrategyRule(ctx context.Context, req *apiv1.UpdateStrategyRuleRequest) (*apiv1.UpdateStrategyRuleReply, error) {
	alertPages := make([]string, 0, len(req.AlertPages))
	alertPages = append(alertPages, req.AlertPages...)
	if err := s.strategyService.UpdateStrategyRule(ctx, snowflake.ParseInt64(req.Uid), req.RuleDetail, vobj.AlertLevel(req.AlertLevel), alertPages, req.Order); err != nil {
		return nil, err
	}
	return &apiv1.UpdateStrategyRuleReply{}, nil
}

func (s *StrategyService) DeleteStrategyRule(ctx context.Context, req *apiv1.DeleteStrategyRuleRequest) (*apiv1.DeleteStrategyRuleReply, error) {
	if err := s.strategyService.DeleteStrategyRule(ctx, snowflake.ParseInt64(req.Uid)); err != nil {
		return nil, err
	}
	return &apiv1.DeleteStrategyRuleReply{}, nil
}

func (s *StrategyService) GetStrategyRule(ctx context.Context, req *apiv1.GetStrategyRuleRequest) (*apiv1.StrategyRuleItem, error) {
	rule, err := s.strategyService.GetStrategyRule(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return toAPIV1StrategyRuleItem(rule), nil
}

func (s *StrategyService) ListStrategyRule(ctx context.Context, req *apiv1.ListStrategyRuleRequest) (*apiv1.ListStrategyRuleReply, error) {
	rules, err := s.strategyService.ListStrategyRule(ctx, snowflake.ParseInt64(req.StrategyUid), vobj.GlobalStatus(req.Status))
	if err != nil {
		return nil, err
	}
	items := make([]*apiv1.StrategyRuleItem, 0, len(rules))
	for _, rule := range rules {
		items = append(items, toAPIV1StrategyRuleItem(rule))
	}
	return &apiv1.ListStrategyRuleReply{
		Rules: items,
	}, nil
}

func (s *StrategyService) UpdateStrategyRuleStatus(ctx context.Context, req *apiv1.UpdateStrategyRuleStatusRequest) (*apiv1.UpdateStrategyRuleStatusReply, error) {
	if err := s.strategyService.UpdateStrategyRuleStatus(ctx, snowflake.ParseInt64(req.Uid), vobj.GlobalStatus(req.Status)); err != nil {
		return nil, err
	}
	return &apiv1.UpdateStrategyRuleStatusReply{}, nil
}

// toAPIV1StrategyRuleItem converts strategy rule entity to API response
func toAPIV1StrategyRuleItem(r *strategy.StrategyRule) *apiv1.StrategyRuleItem {
	return &apiv1.StrategyRuleItem{
		Uid:        r.UID().Int64(),
		StrategyUid: r.StrategyUID().Int64(),
		RuleDetail: r.RuleDetail(),
		Status:     enum.GlobalStatus(r.Status()),
		AlertLevel: int32(r.AlertLevel()),
		AlertPages: r.AlertPages(),
		Order:      r.Order(),
		CreatedAt:  r.CreatedAt().Unix(),
		UpdatedAt:  r.UpdatedAt().Unix(),
	}
}

// toAPIV1ListStrategyRuleReply converts strategy rule list to API response
func toAPIV1ListStrategyRuleReply(rules []*strategy.StrategyRule) *apiv1.ListStrategyRuleReply {
	items := make([]*apiv1.StrategyRuleItem, 0, len(rules))
	for _, r := range rules {
		items = append(items, toAPIV1StrategyRuleItem(r))
	}
	return &apiv1.ListStrategyRuleReply{
		Rules: items,
	}
}

// 策略详细配置更新接口实现
func (s *StrategyService) UpdateStrategyDataSourceConfig(ctx context.Context, req *apiv1.UpdateStrategyDataSourceConfigRequest) (*apiv1.UpdateStrategyDataSourceConfigReply, error) {
	uids := make(map[snowflake.ID]bool)
	for _, uid := range req.DatasourceUids {
		uids[snowflake.ParseInt64(uid)] = true
	}
	if err := s.strategyService.UpdateDataSourceConfig(ctx, snowflake.ParseInt64(req.Uid), uids, req.Query, req.DatasourceType); err != nil {
		return nil, err
	}
	return &apiv1.UpdateStrategyDataSourceConfigReply{}, nil
}

func (s *StrategyService) UpdateStrategyAlertConfig(ctx context.Context, req *apiv1.UpdateStrategyAlertConfigRequest) (*apiv1.UpdateStrategyAlertConfigReply, error) {
	alertPages := make([]string, 0, len(req.AlertPages))
	alertPages = append(alertPages, req.AlertPages...)
	if err := s.strategyService.UpdateAlertConfig(ctx, snowflake.ParseInt64(req.Uid), req.AlertTitle, req.AlertContent, vobj.AlertLevel(req.AlertLevel), alertPages); err != nil {
		return nil, err
	}
	return &apiv1.UpdateStrategyAlertConfigReply{}, nil
}

func (s *StrategyService) UpdateStrategyNotifyConfig(ctx context.Context, req *apiv1.UpdateStrategyNotifyConfigRequest) (*apiv1.UpdateStrategyNotifyConfigReply, error) {
	uids := make(map[snowflake.ID]bool)
	for _, uid := range req.ReceiverUids {
		uids[snowflake.ParseInt64(uid)] = true
	}
	if err := s.strategyService.UpdateNotifyConfig(ctx, snowflake.ParseInt64(req.Uid), uids); err != nil {
		return nil, err
	}
	return &apiv1.UpdateStrategyNotifyConfigReply{}, nil
}

func (s *StrategyService) UpdateStrategyDialTestConfig(ctx context.Context, req *apiv1.UpdateStrategyDialTestConfigRequest) (*apiv1.UpdateStrategyDialTestConfigReply, error) {
	if err := s.strategyService.UpdateDialTestConfig(ctx, snowflake.ParseInt64(req.Uid), vobj.DialTestType(req.DialTestType), req.DialTestTargets); err != nil {
		return nil, err
	}
	return &apiv1.UpdateStrategyDialTestConfigReply{}, nil
}

func (s *StrategyService) UpdateStrategySuppressConfig(ctx context.Context, req *apiv1.UpdateStrategySuppressConfigRequest) (*apiv1.UpdateStrategySuppressConfigReply, error) {
	if err := s.strategyService.UpdateSuppressConfig(ctx, snowflake.ParseInt64(req.Uid), req.SuppressType, req.SuppressConfig); err != nil {
		return nil, err
	}
	return &apiv1.UpdateStrategySuppressConfigReply{}, nil
}

// toAPIV1ReceiverItem converts receiver entity to API response
func toAPIV1ReceiverItem(r *strategy.Receiver) *apiv1.ReceiverItem {
	userIDs := make([]int64, 0, len(r.UserIDs()))
	for uid := range r.UserIDs() {
		userIDs = append(userIDs, uid.Int64())
	}
	notifyTypes := make([]string, 0, len(r.NotifyTypes()))
	for _, t := range r.NotifyTypes() {
		notifyTypes = append(notifyTypes, string(t))
	}
	return &apiv1.ReceiverItem{
		Uid:         r.UID().Int64(),
		NamespaceUid: r.NamespaceUID().Int64(),
		Type:        string(r.Type()),
		Name:        r.Name(),
		Description: r.Description(),
		UserIds:     userIDs,
		LabelMatch:  r.LabelMatch(),
		NotifyTypes: notifyTypes,
		CreatedAt:  r.CreatedAt().Unix(),
		UpdatedAt:   r.UpdatedAt().Unix(),
	}
}

// toAPIV1ListReceiverReply converts receiver page to API response
func toAPIV1ListReceiverReply(page *shared.Page[*strategy.Receiver]) *apiv1.ListReceiverReply {
	items := make([]*apiv1.ReceiverItem, 0, len(page.Items))
	for _, r := range page.Items {
		items = append(items, toAPIV1ReceiverItem(r))
	}
	return &apiv1.ListReceiverReply{
		Receivers: items,
		Total:     page.Total,
	}
}


