package event

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/internal/biz/vobj"
	"github.com/aide-family/sovereign/pkg/merr"
	klog "github.com/go-kratos/kratos/v2/log"
)

type Service struct {
	realtimeAlertRepo      RealtimeAlertRepository
	alertHistoryRepo       AlertHistoryRepository
	notificationHistoryRepo NotificationHistoryRepository
	scheduledTaskRepo      ScheduledTaskRepository
	helper                 *klog.Helper
}

func NewService(realtimeAlertRepo RealtimeAlertRepository, alertHistoryRepo AlertHistoryRepository, notificationHistoryRepo NotificationHistoryRepository, scheduledTaskRepo ScheduledTaskRepository, helper *klog.Helper) *Service {
	return &Service{
		realtimeAlertRepo:       realtimeAlertRepo,
		alertHistoryRepo:        alertHistoryRepo,
		notificationHistoryRepo: notificationHistoryRepo,
		scheduledTaskRepo:       scheduledTaskRepo,
		helper:                  klog.NewHelper(klog.With(helper.Logger(), "biz", "event")),
	}
}

// 实时告警相关方法

func (s *Service) ListRealtimeAlerts(ctx context.Context, query *RealtimeAlertListQuery) (*shared.Page[*RealtimeAlert], error) {
	page, err := s.realtimeAlertRepo.List(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "list realtime alerts failed", "error", err)
		return nil, merr.ErrorInternal("list realtime alerts failed").WithCause(err)
	}
	return page, nil
}

func (s *Service) GetRealtimeAlert(ctx context.Context, uid snowflake.ID) (*RealtimeAlert, error) {
	alert, err := s.realtimeAlertRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("realtime alert %s not found", uid)
		}
		s.helper.Errorw("msg", "get realtime alert failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get realtime alert %s failed", uid).WithCause(err)
	}
	return alert, nil
}

func (s *Service) InterveneAlert(ctx context.Context, uid snowflake.ID, intervener, log string) error {
	alert, err := s.realtimeAlertRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("realtime alert %s not found", uid)
		}
		s.helper.Errorw("msg", "get realtime alert failed", "error", err, "uid", uid)
		return merr.ErrorInternal("get realtime alert %s failed", uid).WithCause(err)
	}

	alert.Intervene(intervener, log)
	if err := s.realtimeAlertRepo.Save(ctx, alert); err != nil {
		s.helper.Errorw("msg", "intervene alert failed", "error", err, "uid", uid)
		return merr.ErrorInternal("intervene alert %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) SuppressAlert(ctx context.Context, uid snowflake.ID, duration int64, reason string) error {
	alert, err := s.realtimeAlertRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("realtime alert %s not found", uid)
		}
		s.helper.Errorw("msg", "get realtime alert failed", "error", err, "uid", uid)
		return merr.ErrorInternal("get realtime alert %s failed", uid).WithCause(err)
	}

	alert.Suppress(duration, reason)
	if err := s.realtimeAlertRepo.Save(ctx, alert); err != nil {
		s.helper.Errorw("msg", "suppress alert failed", "error", err, "uid", uid)
		return merr.ErrorInternal("suppress alert %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpgradeAlert(ctx context.Context, uid snowflake.ID, receiverUIDs []snowflake.ID, log string) error {
	alert, err := s.realtimeAlertRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("realtime alert %s not found", uid)
		}
		s.helper.Errorw("msg", "get realtime alert failed", "error", err, "uid", uid)
		return merr.ErrorInternal("get realtime alert %s failed", uid).WithCause(err)
	}

	alert.Upgrade(receiverUIDs, log)
	if err := s.realtimeAlertRepo.Save(ctx, alert); err != nil {
		s.helper.Errorw("msg", "upgrade alert failed", "error", err, "uid", uid)
		return merr.ErrorInternal("upgrade alert %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) GetAlertEventChart(ctx context.Context, uid snowflake.ID, startTime, endTime int64) (string, error) {
	// TODO: 实现事件图表生成逻辑
	// 这里返回空数据，实际应该根据告警UID查询相关事件数据并生成图表
	s.helper.Infow("msg", "get alert event chart", "uid", uid, "start_time", startTime, "end_time", endTime)
	return "{}", nil
}

// 告警历史相关方法

func (s *Service) ListAlertHistory(ctx context.Context, query *AlertHistoryListQuery) (*shared.Page[*AlertHistory], error) {
	page, err := s.alertHistoryRepo.List(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "list alert history failed", "error", err)
		return nil, merr.ErrorInternal("list alert history failed").WithCause(err)
	}
	return page, nil
}

func (s *Service) GetAlertHistory(ctx context.Context, uid snowflake.ID) (*AlertHistory, error) {
	history, err := s.alertHistoryRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("alert history %s not found", uid)
		}
		s.helper.Errorw("msg", "get alert history failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get alert history %s failed", uid).WithCause(err)
	}
	return history, nil
}

func (s *Service) GetAlertRawData(ctx context.Context, uid snowflake.ID) (string, string, error) {
	// TODO: 实现获取告警原始数据逻辑
	// 这里返回空数据，实际应该从存储中获取原始事件数据
	s.helper.Infow("msg", "get alert raw data", "uid", uid)
	return "{}", "{}", nil
}

func (s *Service) GetAlertEventReport(ctx context.Context, uid snowflake.ID) (string, string, error) {
	// TODO: 实现获取告警事件报告逻辑
	// 这里返回空数据，实际应该生成原始报告和AI分析报告
	s.helper.Infow("msg", "get alert event report", "uid", uid)
	return "{}", "{}", nil
}

// 通知历史相关方法

func (s *Service) ListNotificationHistory(ctx context.Context, query *NotificationHistoryListQuery) (*shared.Page[*NotificationHistory], error) {
	page, err := s.notificationHistoryRepo.List(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "list notification history failed", "error", err)
		return nil, merr.ErrorInternal("list notification history failed").WithCause(err)
	}
	return page, nil
}

func (s *Service) RetryNotification(ctx context.Context, uid snowflake.ID) error {
	notification, err := s.notificationHistoryRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("notification %s not found", uid)
		}
		s.helper.Errorw("msg", "get notification failed", "error", err, "uid", uid)
		return merr.ErrorInternal("get notification %s failed", uid).WithCause(err)
	}

	if notification.RetryCount() >= 3 {
		return merr.ErrorParams("notification retry count exceeded maximum (3)")
	}

	notification.Retry()
	if err := s.notificationHistoryRepo.Save(ctx, notification); err != nil {
		s.helper.Errorw("msg", "retry notification failed", "error", err, "uid", uid)
		return merr.ErrorInternal("retry notification %s failed", uid).WithCause(err)
	}

	// TODO: 实际调用消息服务重新发送通知
	s.helper.Infow("msg", "retry notification", "uid", uid, "retry_count", notification.RetryCount())
	return nil
}

// 事件总线相关方法

func (s *Service) PushExternalAlert(ctx context.Context, source string, alertData map[string]string) (snowflake.ID, error) {
	// TODO: 实现外部告警事件处理逻辑
	// 这里创建实时告警并返回UID
	alert := NewRealtimeAlert(
		vobj.AlertLevelP0, // 默认等级，实际应该从alertData中解析
		alertData["title"],
		alertData["content"],
		0, // 策略UID，外部告警可能没有
		0, // 策略组UID
		alertData,
	)
	if err := s.realtimeAlertRepo.Save(ctx, alert); err != nil {
		s.helper.Errorw("msg", "push external alert failed", "error", err, "source", source)
		return 0, merr.ErrorInternal("push external alert failed").WithCause(err)
	}
	return alert.UID(), nil
}

func (s *Service) PushInternalAlert(ctx context.Context, eventType string, eventData map[string]string) (snowflake.ID, error) {
	// TODO: 实现内部告警事件处理逻辑
	// 这里创建实时告警并返回UID
	alert := NewRealtimeAlert(
		vobj.AlertLevelP0, // 默认等级，实际应该从eventData中解析
		eventData["title"],
		eventData["content"],
		0, // 策略UID，实际应该从eventData中解析
		0, // 策略组UID
		eventData,
	)
	if err := s.realtimeAlertRepo.Save(ctx, alert); err != nil {
		s.helper.Errorw("msg", "push internal alert failed", "error", err, "event_type", eventType)
		return 0, merr.ErrorInternal("push internal alert failed").WithCause(err)
	}
	return alert.UID(), nil
}

// 定时任务相关方法

func (s *Service) ListScheduledTasks(ctx context.Context, query *ScheduledTaskListQuery) (*shared.Page[*ScheduledTask], error) {
	page, err := s.scheduledTaskRepo.List(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "list scheduled tasks failed", "error", err)
		return nil, merr.ErrorInternal("list scheduled tasks failed").WithCause(err)
	}
	return page, nil
}

func (s *Service) GetScheduledTask(ctx context.Context, uid snowflake.ID) (*ScheduledTask, error) {
	task, err := s.scheduledTaskRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("scheduled task %s not found", uid)
		}
		s.helper.Errorw("msg", "get scheduled task failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get scheduled task %s failed", uid).WithCause(err)
	}
	return task, nil
}

func (s *Service) UpdateScheduledTaskStatus(ctx context.Context, uid snowflake.ID, status string) error {
	task, err := s.scheduledTaskRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("scheduled task %s not found", uid)
		}
		s.helper.Errorw("msg", "get scheduled task failed", "error", err, "uid", uid)
		return merr.ErrorInternal("get scheduled task %s failed", uid).WithCause(err)
	}

	if status == "running" {
		task.Start()
	} else if status == "stopped" {
		task.Stop()
	} else {
		return merr.ErrorParams("invalid status: %s", status)
	}

	if err := s.scheduledTaskRepo.Save(ctx, task); err != nil {
		s.helper.Errorw("msg", "update scheduled task status failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update scheduled task status %s failed", uid).WithCause(err)
	}
	return nil
}

