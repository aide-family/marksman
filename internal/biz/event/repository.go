package event

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/internal/biz/vobj"
)

// RealtimeAlertRepository 实时告警仓库接口
type RealtimeAlertRepository interface {
	Save(ctx context.Context, alert *RealtimeAlert) error
	FindByID(ctx context.Context, uid snowflake.ID) (*RealtimeAlert, error)
	List(ctx context.Context, query *RealtimeAlertListQuery) (*shared.Page[*RealtimeAlert], error)
	Delete(ctx context.Context, uid snowflake.ID) error
}

// AlertHistoryRepository 告警历史仓库接口
type AlertHistoryRepository interface {
	Save(ctx context.Context, history *AlertHistory) error
	FindByID(ctx context.Context, uid snowflake.ID) (*AlertHistory, error)
	List(ctx context.Context, query *AlertHistoryListQuery) (*shared.Page[*AlertHistory], error)
}

// NotificationHistoryRepository 通知历史仓库接口
type NotificationHistoryRepository interface {
	Save(ctx context.Context, notification *NotificationHistory) error
	FindByID(ctx context.Context, uid snowflake.ID) (*NotificationHistory, error)
	List(ctx context.Context, query *NotificationHistoryListQuery) (*shared.Page[*NotificationHistory], error)
}

// ScheduledTaskRepository 定时任务仓库接口
type ScheduledTaskRepository interface {
	Save(ctx context.Context, task *ScheduledTask) error
	FindByID(ctx context.Context, uid snowflake.ID) (*ScheduledTask, error)
	List(ctx context.Context, query *ScheduledTaskListQuery) (*shared.Page[*ScheduledTask], error)
}

// RealtimeAlertListQuery 实时告警列表查询
type RealtimeAlertListQuery struct {
	*shared.PageRequest
	AlertLevel       vobj.AlertLevel
	StrategyUID      snowflake.ID
	StrategyGroupUID snowflake.ID
	Keyword          string
	OnlyUnhandled    bool
}

func NewRealtimeAlertListQuery(page, pageSize int32, alertLevel vobj.AlertLevel, strategyUID, strategyGroupUID snowflake.ID, keyword string, onlyUnhandled bool) *RealtimeAlertListQuery {
	return &RealtimeAlertListQuery{
		PageRequest:      shared.NewPageRequest(page, pageSize),
		AlertLevel:       alertLevel,
		StrategyUID:      strategyUID,
		StrategyGroupUID: strategyGroupUID,
		Keyword:          keyword,
		OnlyUnhandled:    onlyUnhandled,
	}
}

// AlertHistoryListQuery 告警历史列表查询
type AlertHistoryListQuery struct {
	*shared.PageRequest
	AlertLevel       vobj.AlertLevel
	StrategyUID      snowflake.ID
	StrategyGroupUID snowflake.ID
	Keyword          string
	OnlyRecovered    bool
	StartTime        int64
	EndTime          int64
}

func NewAlertHistoryListQuery(page, pageSize int32, alertLevel vobj.AlertLevel, strategyUID, strategyGroupUID snowflake.ID, keyword string, onlyRecovered bool, startTime, endTime int64) *AlertHistoryListQuery {
	return &AlertHistoryListQuery{
		PageRequest:      shared.NewPageRequest(page, pageSize),
		AlertLevel:       alertLevel,
		StrategyUID:      strategyUID,
		StrategyGroupUID: strategyGroupUID,
		Keyword:          keyword,
		OnlyRecovered:    onlyRecovered,
		StartTime:        startTime,
		EndTime:          endTime,
	}
}

// NotificationHistoryListQuery 通知历史列表查询
type NotificationHistoryListQuery struct {
	*shared.PageRequest
	AlertUID   snowflake.ID
	Status     string
	NotifyType vobj.NotifyType
	StartTime  int64
	EndTime    int64
}

func NewNotificationHistoryListQuery(page, pageSize int32, alertUID snowflake.ID, status string, notifyType vobj.NotifyType, startTime, endTime int64) *NotificationHistoryListQuery {
	return &NotificationHistoryListQuery{
		PageRequest: shared.NewPageRequest(page, pageSize),
		AlertUID:    alertUID,
		Status:      status,
		NotifyType:  notifyType,
		StartTime:   startTime,
		EndTime:     endTime,
	}
}

// ScheduledTaskListQuery 定时任务列表查询
type ScheduledTaskListQuery struct {
	*shared.PageRequest
	StrategyUID  snowflake.ID
	ExecutorType string
	Status       string
}

func NewScheduledTaskListQuery(page, pageSize int32, strategyUID snowflake.ID, executorType, status string) *ScheduledTaskListQuery {
	return &ScheduledTaskListQuery{
		PageRequest:  shared.NewPageRequest(page, pageSize),
		StrategyUID:  strategyUID,
		ExecutorType: executorType,
		Status:       status,
	}
}

