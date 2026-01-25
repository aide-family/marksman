package service

import (
	"context"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/internal/biz/subscription"
	"github.com/aide-family/sovereign/internal/biz/vobj"
	apiv1 "github.com/aide-family/sovereign/pkg/api/v1"
	"github.com/aide-family/sovereign/pkg/enum"
)

func NewSubscriptionService(subscriptionService *subscription.Service) *SubscriptionService {
	return &SubscriptionService{
		subscriptionService: subscriptionService,
	}
}

type SubscriptionService struct {
	apiv1.UnimplementedSubscriptionServer

	subscriptionService *subscription.Service
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, req *apiv1.CreateSubscriptionRequest) (*apiv1.CreateSubscriptionReply, error) {
	if err := s.subscriptionService.Create(ctx, snowflake.ParseInt64(req.UserId), snowflake.ParseInt64(req.NamespaceUid), subscription.SubscriptionType(req.Type), req.Name); err != nil {
		return nil, err
	}
	return &apiv1.CreateSubscriptionReply{}, nil
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, req *apiv1.UpdateSubscriptionRequest) (*apiv1.UpdateSubscriptionReply, error) {
	if err := s.subscriptionService.Update(ctx, snowflake.ParseInt64(req.Uid), req.Name, req.Description); err != nil {
		return nil, err
	}
	return &apiv1.UpdateSubscriptionReply{}, nil
}

func (s *SubscriptionService) UpdateSubscriptionStatus(ctx context.Context, req *apiv1.UpdateSubscriptionStatusRequest) (*apiv1.UpdateSubscriptionStatusReply, error) {
	if err := s.subscriptionService.UpdateStatus(ctx, snowflake.ParseInt64(req.Uid), vobj.GlobalStatus(req.Status)); err != nil {
		return nil, err
	}
	return &apiv1.UpdateSubscriptionStatusReply{}, nil
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, req *apiv1.DeleteSubscriptionRequest) (*apiv1.DeleteSubscriptionReply, error) {
	if err := s.subscriptionService.Delete(ctx, snowflake.ParseInt64(req.Uid)); err != nil {
		return nil, err
	}
	return &apiv1.DeleteSubscriptionReply{}, nil
}

func (s *SubscriptionService) GetSubscription(ctx context.Context, req *apiv1.GetSubscriptionRequest) (*apiv1.SubscriptionItem, error) {
	subscriptionEntity, err := s.subscriptionService.Get(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return toAPIV1SubscriptionItem(subscriptionEntity), nil
}

func (s *SubscriptionService) ListSubscription(ctx context.Context, req *apiv1.ListSubscriptionRequest) (*apiv1.ListSubscriptionReply, error) {
	query := &subscription.ListQuery{
		PageRequest:  shared.NewPageRequest(req.Page, req.PageSize),
		UserID:       snowflake.ParseInt64(req.UserId),
		NamespaceUID: snowflake.ParseInt64(req.NamespaceUid),
		Type:         subscription.SubscriptionType(req.Type),
		Keyword:      req.Keyword,
		Status:       vobj.GlobalStatus(req.Status),
	}
	page, err := s.subscriptionService.List(ctx, query)
	if err != nil {
		return nil, err
	}
	return toAPIV1ListSubscriptionReply(page), nil
}

func (s *SubscriptionService) SelectSubscription(ctx context.Context, req *apiv1.SelectSubscriptionRequest) (*apiv1.SelectSubscriptionReply, error) {
	var nextUID snowflake.ID
	if req.NextUid > 0 {
		nextUID = snowflake.ParseInt64(req.NextUid)
	}
	query := &subscription.SelectQuery{
		UserID:       snowflake.ParseInt64(req.UserId),
		NamespaceUID: snowflake.ParseInt64(req.NamespaceUid),
		Keyword:      req.Keyword,
		Limit:        req.Limit,
		NextUID:      nextUID,
		Status:       vobj.GlobalStatus(req.Status),
	}
	result, err := s.subscriptionService.Select(ctx, query)
	if err != nil {
		return nil, err
	}
	return toAPIV1SelectSubscriptionReply(result), nil
}

func (s *SubscriptionService) UpdateSubscriptionStrategyGroups(ctx context.Context, req *apiv1.UpdateSubscriptionStrategyGroupsRequest) (*apiv1.UpdateSubscriptionStrategyGroupsReply, error) {
	uids := make(map[snowflake.ID]bool)
	for _, uid := range req.StrategyGroupUids {
		uids[snowflake.ParseInt64(uid)] = true
	}
	if err := s.subscriptionService.UpdateStrategyGroupUIDs(ctx, snowflake.ParseInt64(req.Uid), uids); err != nil {
		return nil, err
	}
	return &apiv1.UpdateSubscriptionStrategyGroupsReply{}, nil
}

func (s *SubscriptionService) UpdateSubscriptionDataSources(ctx context.Context, req *apiv1.UpdateSubscriptionDataSourcesRequest) (*apiv1.UpdateSubscriptionDataSourcesReply, error) {
	uids := make(map[snowflake.ID]bool)
	for _, uid := range req.DatasourceUids {
		uids[snowflake.ParseInt64(uid)] = true
	}
	if err := s.subscriptionService.UpdateDataSourceUIDs(ctx, snowflake.ParseInt64(req.Uid), uids); err != nil {
		return nil, err
	}
	return &apiv1.UpdateSubscriptionDataSourcesReply{}, nil
}

func (s *SubscriptionService) UpdateSubscriptionAlertLevels(ctx context.Context, req *apiv1.UpdateSubscriptionAlertLevelsRequest) (*apiv1.UpdateSubscriptionAlertLevelsReply, error) {
	levels := make([]vobj.AlertLevel, 0, len(req.AlertLevels))
	for _, level := range req.AlertLevels {
		levels = append(levels, vobj.AlertLevel(level))
	}
	if err := s.subscriptionService.UpdateAlertLevels(ctx, snowflake.ParseInt64(req.Uid), levels); err != nil {
		return nil, err
	}
	return &apiv1.UpdateSubscriptionAlertLevelsReply{}, nil
}

func (s *SubscriptionService) UpdateSubscriptionNotifyTypes(ctx context.Context, req *apiv1.UpdateSubscriptionNotifyTypesRequest) (*apiv1.UpdateSubscriptionNotifyTypesReply, error) {
	types := make([]vobj.NotifyType, 0, len(req.NotifyTypes))
	for _, t := range req.NotifyTypes {
		types = append(types, vobj.NotifyType(t))
	}
	if err := s.subscriptionService.UpdateNotifyTypes(ctx, snowflake.ParseInt64(req.Uid), types); err != nil {
		return nil, err
	}
	return &apiv1.UpdateSubscriptionNotifyTypesReply{}, nil
}

// toAPIV1SubscriptionItem converts subscription entity to API response
func toAPIV1SubscriptionItem(s *subscription.Subscription) *apiv1.SubscriptionItem {
	strategyGroupUIDs := make([]int64, 0, len(s.StrategyGroupUIDs()))
	for uid := range s.StrategyGroupUIDs() {
		strategyGroupUIDs = append(strategyGroupUIDs, uid.Int64())
	}

	datasourceUIDs := make([]int64, 0, len(s.DataSourceUIDs()))
	for uid := range s.DataSourceUIDs() {
		datasourceUIDs = append(datasourceUIDs, uid.Int64())
	}

	alertLevels := make([]int32, 0, len(s.AlertLevels()))
	for _, level := range s.AlertLevels() {
		alertLevels = append(alertLevels, int32(level))
	}

	notifyTypes := make([]string, 0, len(s.NotifyTypes()))
	for _, t := range s.NotifyTypes() {
		notifyTypes = append(notifyTypes, string(t))
	}

	return &apiv1.SubscriptionItem{
		Uid:              s.UID().Int64(),
		UserId:           s.UserID().Int64(),
		NamespaceUid:    s.NamespaceUID().Int64(),
		Type:             string(s.Type()),
		Name:             s.Name(),
		Description:      s.Description(),
		Status:           enum.GlobalStatus(s.Status()),
		StrategyGroupUids: strategyGroupUIDs,
		DatasourceUids:   datasourceUIDs,
		AlertLevels:      alertLevels,
		NotifyTypes:      notifyTypes,
		CreatedAt:        s.CreatedAt().Unix(),
		UpdatedAt:        s.UpdatedAt().Unix(),
	}
}

// toAPIV1ListSubscriptionReply converts subscription page to API response
func toAPIV1ListSubscriptionReply(page *shared.Page[*subscription.Subscription]) *apiv1.ListSubscriptionReply {
	items := make([]*apiv1.SubscriptionItem, 0, len(page.Items))
	for _, s := range page.Items {
		items = append(items, toAPIV1SubscriptionItem(s))
	}
	return &apiv1.ListSubscriptionReply{
		Subscriptions: items,
		Total:         page.Total,
	}
}

// toAPIV1SelectSubscriptionReply converts subscription select result to API response
func toAPIV1SelectSubscriptionReply(result *subscription.SelectResult) *apiv1.SelectSubscriptionReply {
	selectItems := make([]*apiv1.SelectItem, 0, len(result.Items))
	for _, item := range result.Items {
		selectItems = append(selectItems, &apiv1.SelectItem{
			Value:    item.UID.String(),
			Label:    item.Name,
			Disabled: item.Disabled,
			Tooltip:  item.Tooltip,
		})
	}
	return &apiv1.SelectSubscriptionReply{
		Items:   selectItems,
		Total:   result.Total,
		NextUid: result.NextUID.Int64(),
		HasMore: result.HasMore,
	}
}

