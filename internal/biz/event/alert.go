package event

import (
	"time"

	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/vobj"
)

// RealtimeAlert 实时告警实体
type RealtimeAlert struct {
	uid              snowflake.ID
	alertLevel       vobj.AlertLevel
	alertTime        time.Time
	alertTitle      string
	alertContent     string
	intervener       string
	strategyUID      snowflake.ID
	strategyGroupUID snowflake.ID
	isSuppressed     bool
	isUpgraded       bool
	labels           map[string]string
	createdAt        time.Time
	updatedAt        time.Time
}

func NewRealtimeAlert(alertLevel vobj.AlertLevel, alertTitle, alertContent string, strategyUID, strategyGroupUID snowflake.ID, labels map[string]string) *RealtimeAlert {
	var uid snowflake.ID
	node, err := snowflake.NewNode(hello.NodeID())
	if err == nil {
		uid = node.Generate()
	}
	return &RealtimeAlert{
		uid:              uid,
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

func (a *RealtimeAlert) Intervene(intervener, log string) {
	a.intervener = intervener
	a.updatedAt = time.Now()
	// TODO: 记录介入日志
}

func (a *RealtimeAlert) Suppress(duration int64, reason string) {
	a.isSuppressed = true
	a.updatedAt = time.Now()
	// TODO: 记录抑制信息
}

func (a *RealtimeAlert) Upgrade(receiverUIDs []snowflake.ID, log string) {
	a.isUpgraded = true
	a.updatedAt = time.Now()
	// TODO: 记录升级信息
}

// FromModel creates a RealtimeAlert entity from repository model
func RealtimeAlertFromModel(uid snowflake.ID, alertLevel vobj.AlertLevel, alertTime time.Time, alertTitle, alertContent, intervener string, strategyUID, strategyGroupUID snowflake.ID, isSuppressed, isUpgraded bool, labels map[string]string, createdAt, updatedAt time.Time) *RealtimeAlert {
	return &RealtimeAlert{
		uid:              uid,
		alertLevel:       alertLevel,
		alertTime:        alertTime,
		alertTitle:       alertTitle,
		alertContent:     alertContent,
		intervener:       intervener,
		strategyUID:      strategyUID,
		strategyGroupUID: strategyGroupUID,
		isSuppressed:     isSuppressed,
		isUpgraded:       isUpgraded,
		labels:           labels,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}

// Getters
func (a *RealtimeAlert) UID() snowflake.ID              { return a.uid }
func (a *RealtimeAlert) AlertLevel() vobj.AlertLevel    { return a.alertLevel }
func (a *RealtimeAlert) AlertTime() time.Time            { return a.alertTime }
func (a *RealtimeAlert) AlertTitle() string              { return a.alertTitle }
func (a *RealtimeAlert) AlertContent() string            { return a.alertContent }
func (a *RealtimeAlert) Intervener() string              { return a.intervener }
func (a *RealtimeAlert) StrategyUID() snowflake.ID      { return a.strategyUID }
func (a *RealtimeAlert) StrategyGroupUID() snowflake.ID { return a.strategyGroupUID }
func (a *RealtimeAlert) IsSuppressed() bool              { return a.isSuppressed }
func (a *RealtimeAlert) IsUpgraded() bool                { return a.isUpgraded }
func (a *RealtimeAlert) Labels() map[string]string       { return a.labels }
func (a *RealtimeAlert) CreatedAt() time.Time            { return a.createdAt }
func (a *RealtimeAlert) UpdatedAt() time.Time            { return a.updatedAt }
func (a *RealtimeAlert) Duration() int64                 { return int64(time.Since(a.alertTime).Seconds()) }

