package service

import (
	"context"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/sovereign/internal/biz/event"
	"github.com/aide-family/sovereign/internal/biz/vobj"
	apiv1 "github.com/aide-family/sovereign/pkg/api/v1"
)

func NewEventCenterService(eventService *event.Service) *EventCenterService {
	return &EventCenterService{
		eventService: eventService,
	}
}

type EventCenterService struct {
	// apiv1.UnimplementedEventCenterServer // 注释掉，proto生成后取消注释
	eventService *event.Service
}

// 实时告警相关接口实现

func (s *EventCenterService) ListRealtimeAlerts(ctx context.Context, req *apiv1.ListRealtimeAlertsRequest) (*apiv1.ListRealtimeAlertsReply, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	query := event.NewRealtimeAlertListQuery(
		page,
		pageSize,
		vobj.AlertLevel(req.AlertLevel),
		snowflake.ParseInt64(req.StrategyUid),
		snowflake.ParseInt64(req.StrategyGroupUid),
		req.Keyword,
		req.OnlyUnhandled,
	)

	pageResult, err := s.eventService.ListRealtimeAlerts(ctx, query)
	if err != nil {
		return nil, err
	}

	alerts := make([]*apiv1.RealtimeAlertItem, 0, len(pageResult.Items))
	for _, a := range pageResult.Items {
		alerts = append(alerts, toAPIV1RealtimeAlertItem(a))
	}

	return &apiv1.ListRealtimeAlertsReply{
		Alerts: alerts,
		Total:  pageResult.Total,
	}, nil
}

func (s *EventCenterService) GetRealtimeAlert(ctx context.Context, req *apiv1.GetRealtimeAlertRequest) (*apiv1.RealtimeAlertItem, error) {
	alert, err := s.eventService.GetRealtimeAlert(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return toAPIV1RealtimeAlertItem(alert), nil
}

func (s *EventCenterService) InterveneAlert(ctx context.Context, req *apiv1.InterveneAlertRequest) (*apiv1.InterveneAlertReply, error) {
	if err := s.eventService.InterveneAlert(ctx, snowflake.ParseInt64(req.Uid), "", req.Log); err != nil {
		return nil, err
	}
	return &apiv1.InterveneAlertReply{}, nil
}

func (s *EventCenterService) SuppressAlert(ctx context.Context, req *apiv1.SuppressAlertRequest) (*apiv1.SuppressAlertReply, error) {
	if err := s.eventService.SuppressAlert(ctx, snowflake.ParseInt64(req.Uid), req.SuppressDuration, req.Reason); err != nil {
		return nil, err
	}
	return &apiv1.SuppressAlertReply{}, nil
}

func (s *EventCenterService) UpgradeAlert(ctx context.Context, req *apiv1.UpgradeAlertRequest) (*apiv1.UpgradeAlertReply, error) {
	receiverUIDs := make([]snowflake.ID, 0, len(req.ReceiverUids))
	for _, uid := range req.ReceiverUids {
		receiverUIDs = append(receiverUIDs, snowflake.ParseInt64(uid))
	}
	if err := s.eventService.UpgradeAlert(ctx, snowflake.ParseInt64(req.Uid), receiverUIDs, req.Log); err != nil {
		return nil, err
	}
	return &apiv1.UpgradeAlertReply{}, nil
}

func (s *EventCenterService) GetAlertEventChart(ctx context.Context, req *apiv1.GetAlertEventChartRequest) (*apiv1.GetAlertEventChartReply, error) {
	chartData, err := s.eventService.GetAlertEventChart(ctx, snowflake.ParseInt64(req.Uid), req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	return &apiv1.GetAlertEventChartReply{
		ChartData: chartData,
	}, nil
}

// 告警历史相关接口实现

func (s *EventCenterService) ListAlertHistory(ctx context.Context, req *apiv1.ListAlertHistoryRequest) (*apiv1.ListAlertHistoryReply, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	query := event.NewAlertHistoryListQuery(
		page,
		pageSize,
		vobj.AlertLevel(req.AlertLevel),
		snowflake.ParseInt64(req.StrategyUid),
		snowflake.ParseInt64(req.StrategyGroupUid),
		req.Keyword,
		req.OnlyRecovered,
		req.StartTime,
		req.EndTime,
	)

	pageResult, err := s.eventService.ListAlertHistory(ctx, query)
	if err != nil {
		return nil, err
	}

	alerts := make([]*apiv1.AlertHistoryItem, 0, len(pageResult.Items))
	for _, h := range pageResult.Items {
		alerts = append(alerts, toAPIV1AlertHistoryItem(h))
	}

	return &apiv1.ListAlertHistoryReply{
		Alerts: alerts,
		Total:  pageResult.Total,
	}, nil
}

func (s *EventCenterService) GetAlertHistory(ctx context.Context, req *apiv1.GetAlertHistoryRequest) (*apiv1.AlertHistoryItem, error) {
	history, err := s.eventService.GetAlertHistory(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return toAPIV1AlertHistoryItem(history), nil
}

func (s *EventCenterService) GetAlertRawData(ctx context.Context, req *apiv1.GetAlertRawDataRequest) (*apiv1.GetAlertRawDataReply, error) {
	rawData, recoveredRawData, err := s.eventService.GetAlertRawData(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return &apiv1.GetAlertRawDataReply{
		RawData:         rawData,
		RecoveredRawData: recoveredRawData,
	}, nil
}

func (s *EventCenterService) GetAlertEventReport(ctx context.Context, req *apiv1.GetAlertEventReportRequest) (*apiv1.GetAlertEventReportReply, error) {
	originalReport, aiReport, err := s.eventService.GetAlertEventReport(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return &apiv1.GetAlertEventReportReply{
		OriginalReport:    originalReport,
		AiAnalysisReport: aiReport,
	}, nil
}

// 通知历史相关接口实现

func (s *EventCenterService) ListNotificationHistory(ctx context.Context, req *apiv1.ListNotificationHistoryRequest) (*apiv1.ListNotificationHistoryReply, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	query := event.NewNotificationHistoryListQuery(
		page,
		pageSize,
		snowflake.ParseInt64(req.AlertUid),
		req.Status,
		vobj.NotifyType(req.NotifyType),
		req.StartTime,
		req.EndTime,
	)

	pageResult, err := s.eventService.ListNotificationHistory(ctx, query)
	if err != nil {
		return nil, err
	}

	notifications := make([]*apiv1.NotificationHistoryItem, 0, len(pageResult.Items))
	for _, n := range pageResult.Items {
		notifications = append(notifications, toAPIV1NotificationHistoryItem(n))
	}

	return &apiv1.ListNotificationHistoryReply{
		Notifications: notifications,
		Total:         pageResult.Total,
	}, nil
}

func (s *EventCenterService) RetryNotification(ctx context.Context, req *apiv1.RetryNotificationRequest) (*apiv1.RetryNotificationReply, error) {
	if err := s.eventService.RetryNotification(ctx, snowflake.ParseInt64(req.Uid)); err != nil {
		return nil, err
	}
	return &apiv1.RetryNotificationReply{}, nil
}

// 事件总线相关接口实现

func (s *EventCenterService) PushExternalAlert(ctx context.Context, req *apiv1.PushExternalAlertRequest) (*apiv1.PushExternalAlertReply, error) {
	alertUID, err := s.eventService.PushExternalAlert(ctx, req.Source, req.AlertData)
	if err != nil {
		return nil, err
	}
	return &apiv1.PushExternalAlertReply{
		AlertUid: alertUID.Int64(),
	}, nil
}

func (s *EventCenterService) PushInternalAlert(ctx context.Context, req *apiv1.PushInternalAlertRequest) (*apiv1.PushInternalAlertReply, error) {
	alertUID, err := s.eventService.PushInternalAlert(ctx, req.EventType, req.EventData)
	if err != nil {
		return nil, err
	}
	return &apiv1.PushInternalAlertReply{
		AlertUid: alertUID.Int64(),
	}, nil
}

// 定时任务相关接口实现

func (s *EventCenterService) ListScheduledTasks(ctx context.Context, req *apiv1.ListScheduledTasksRequest) (*apiv1.ListScheduledTasksReply, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	query := event.NewScheduledTaskListQuery(
		page,
		pageSize,
		snowflake.ParseInt64(req.StrategyUid),
		req.ExecutorType,
		req.Status,
	)

	pageResult, err := s.eventService.ListScheduledTasks(ctx, query)
	if err != nil {
		return nil, err
	}

	tasks := make([]*apiv1.ScheduledTaskItem, 0, len(pageResult.Items))
	for _, t := range pageResult.Items {
		tasks = append(tasks, toAPIV1ScheduledTaskItem(t))
	}

	return &apiv1.ListScheduledTasksReply{
		Tasks: tasks,
		Total: pageResult.Total,
	}, nil
}

func (s *EventCenterService) GetScheduledTask(ctx context.Context, req *apiv1.GetScheduledTaskRequest) (*apiv1.ScheduledTaskItem, error) {
	task, err := s.eventService.GetScheduledTask(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return toAPIV1ScheduledTaskItem(task), nil
}

func (s *EventCenterService) UpdateScheduledTaskStatus(ctx context.Context, req *apiv1.UpdateScheduledTaskStatusRequest) (*apiv1.UpdateScheduledTaskStatusReply, error) {
	if err := s.eventService.UpdateScheduledTaskStatus(ctx, snowflake.ParseInt64(req.Uid), req.Status); err != nil {
		return nil, err
	}
	return &apiv1.UpdateScheduledTaskStatusReply{}, nil
}

// 辅助函数

func toAPIV1RealtimeAlertItem(alert *event.RealtimeAlert) *apiv1.RealtimeAlertItem {
	return &apiv1.RealtimeAlertItem{
		Uid:             alert.UID().Int64(),
		AlertLevel:      int32(alert.AlertLevel()),
		AlertTime:       alert.AlertTime().Unix(),
		Duration:        alert.Duration(),
		AlertTitle:      alert.AlertTitle(),
		AlertContent:    alert.AlertContent(),
		Intervener:      alert.Intervener(),
		StrategyUid:     alert.StrategyUID().Int64(),
		StrategyGroupUid: alert.StrategyGroupUID().Int64(),
		IsSuppressed:    alert.IsSuppressed(),
		IsUpgraded:      alert.IsUpgraded(),
		Labels:          alert.Labels(),
		CreatedAt:       alert.CreatedAt().Unix(),
		UpdatedAt:       alert.UpdatedAt().Unix(),
	}
}

func toAPIV1AlertHistoryItem(history *event.AlertHistory) *apiv1.AlertHistoryItem {
	recoveredAt := int64(0)
	if history.IsRecovered() {
		recoveredAt = history.RecoveredAt().Unix()
	}
	return &apiv1.AlertHistoryItem{
		Uid:                history.UID().Int64(),
		AlertLevel:         int32(history.AlertLevel()),
		AlertTime:          history.AlertTime().Unix(),
		Duration:           history.Duration(),
		IsRecovered:        history.IsRecovered(),
		RecoveredAt:        recoveredAt,
		AlertTitle:         history.AlertTitle(),
		AlertContent:       history.AlertContent(),
		ProcessingDuration: history.ProcessingDuration(),
		StrategyUid:        history.StrategyUID().Int64(),
		StrategyGroupUid:   history.StrategyGroupUID().Int64(),
		Labels:             history.Labels(),
		CreatedAt:          history.CreatedAt().Unix(),
		UpdatedAt:          history.UpdatedAt().Unix(),
	}
}

func toAPIV1NotificationHistoryItem(notification *event.NotificationHistory) *apiv1.NotificationHistoryItem {
	return &apiv1.NotificationHistoryItem{
		Uid:              notification.UID().Int64(),
		NotificationTime: notification.NotificationTime().Unix(),
		Status:           notification.Status(),
		NotifyType:       notification.NotifyType().String(),
		Receiver:         notification.Receiver(),
		ReceiverType:     notification.ReceiverType().String(),
		Content:          notification.Content(),
		AlertUid:         notification.AlertUID().Int64(),
		RetryCount:       notification.RetryCount(),
		CreatedAt:        notification.CreatedAt().Unix(),
		UpdatedAt:        notification.UpdatedAt().Unix(),
	}
}

func toAPIV1ScheduledTaskItem(task *event.ScheduledTask) *apiv1.ScheduledTaskItem {
	return &apiv1.ScheduledTaskItem{
		Uid:             task.UID().Int64(),
		StrategyUid:     task.StrategyUID().Int64(),
		ExecutorType:    task.ExecutorType(),
		Status:          task.Status(),
		LastExecuteTime: task.LastExecuteTime().Unix(),
		NextExecuteTime: task.NextExecuteTime().Unix(),
		CreatedAt:       task.CreatedAt().Unix(),
		UpdatedAt:       task.UpdatedAt().Unix(),
	}
}

