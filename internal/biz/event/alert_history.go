package event

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/vobj"
)

// AlertHistory 告警历史实体
type AlertHistory struct {
	uid               snowflake.ID
	alertLevel        vobj.AlertLevel
	alertTime         time.Time
	duration          int64 // 持续时长（秒）
	isRecovered       bool
	recoveredAt       time.Time
	alertTitle        string
	alertContent      string
	processingDuration int64 // 处理耗时（从介入到恢复，秒）
	strategyUID       snowflake.ID
	strategyGroupUID  snowflake.ID
	labels            map[string]string
	createdAt         time.Time
	updatedAt         time.Time
}

func NewAlertHistory(alertLevel vobj.AlertLevel, alertTitle, alertContent string, strategyUID, strategyGroupUID snowflake.ID, labels map[string]string) *AlertHistory {
	return &AlertHistory{
		alertLevel:       alertLevel,
		alertTime:        time.Now(),
		alertTitle:       alertTitle,
		alertContent:     alertContent,
		strategyUID:      strategyUID,
		strategyGroupUID: strategyGroupUID,
		labels:           labels,
		createdAt:        time.Now(),
		updatedAt:        time.Now(),
	}
}

func (h *AlertHistory) Recover() {
	if h.isRecovered {
		return
	}
	h.isRecovered = true
	h.recoveredAt = time.Now()
	h.duration = int64(h.recoveredAt.Sub(h.alertTime).Seconds())
	h.updatedAt = time.Now()
}

func (h *AlertHistory) UpdateProcessingDuration(duration int64) {
	h.processingDuration = duration
	h.updatedAt = time.Now()
}

// FromModel creates an AlertHistory entity from repository model
func AlertHistoryFromModel(uid snowflake.ID, alertLevel vobj.AlertLevel, alertTime time.Time, duration int64, isRecovered bool, recoveredAt time.Time, alertTitle, alertContent string, processingDuration int64, strategyUID, strategyGroupUID snowflake.ID, labels map[string]string, createdAt, updatedAt time.Time) *AlertHistory {
	return &AlertHistory{
		uid:                uid,
		alertLevel:         alertLevel,
		alertTime:          alertTime,
		duration:           duration,
		isRecovered:        isRecovered,
		recoveredAt:        recoveredAt,
		alertTitle:         alertTitle,
		alertContent:       alertContent,
		processingDuration: processingDuration,
		strategyUID:        strategyUID,
		strategyGroupUID:   strategyGroupUID,
		labels:             labels,
		createdAt:          createdAt,
		updatedAt:          updatedAt,
	}
}

// Getters
func (h *AlertHistory) UID() snowflake.ID              { return h.uid }
func (h *AlertHistory) AlertLevel() vobj.AlertLevel    { return h.alertLevel }
func (h *AlertHistory) AlertTime() time.Time            { return h.alertTime }
func (h *AlertHistory) Duration() int64                 { return h.duration }
func (h *AlertHistory) IsRecovered() bool               { return h.isRecovered }
func (h *AlertHistory) RecoveredAt() time.Time          { return h.recoveredAt }
func (h *AlertHistory) AlertTitle() string              { return h.alertTitle }
func (h *AlertHistory) AlertContent() string            { return h.alertContent }
func (h *AlertHistory) ProcessingDuration() int64      { return h.processingDuration }
func (h *AlertHistory) StrategyUID() snowflake.ID       { return h.strategyUID }
func (h *AlertHistory) StrategyGroupUID() snowflake.ID { return h.strategyGroupUID }
func (h *AlertHistory) Labels() map[string]string       { return h.labels }
func (h *AlertHistory) CreatedAt() time.Time            { return h.createdAt }
func (h *AlertHistory) UpdatedAt() time.Time            { return h.updatedAt }

