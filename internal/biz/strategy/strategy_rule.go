package strategy

import (
	"time"
	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/vobj"
	"github.com/aide-family/sovereign/pkg/merr"
)

type StrategyRule struct {
	uid         snowflake.ID
	strategyUID snowflake.ID
	ruleDetail  string
	status      vobj.GlobalStatus
	alertLevel  vobj.AlertLevel
	alertPages  []string
	order       int32
	createdAt   time.Time
	updatedAt   time.Time
}

func NewStrategyRule(strategyUID snowflake.ID, ruleDetail string) *StrategyRule {
	var uid snowflake.ID
	node, err := snowflake.NewNode(hello.NodeID())
	if err == nil {
		uid = node.Generate()
	}
	return &StrategyRule{
		uid:         uid,
		strategyUID: strategyUID,
		ruleDetail:  ruleDetail,
		status:      vobj.GlobalStatusEnabled,
		alertPages:  make([]string, 0),
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}
}

func (r *StrategyRule) Enable() error {
	if r.status == vobj.GlobalStatusEnabled {
		return merr.ErrorParams("strategy rule already enabled")
	}
	r.status = vobj.GlobalStatusEnabled
	r.updatedAt = time.Now()
	return nil
}

func (r *StrategyRule) Disable() error {
	r.status = vobj.GlobalStatusDisabled
	r.updatedAt = time.Now()
	return nil
}

func (r *StrategyRule) IsEnabled() bool {
	return r.status == vobj.GlobalStatusEnabled
}

func (r *StrategyRule) UpdateRuleDetail(ruleDetail string) {
	r.ruleDetail = ruleDetail
	r.updatedAt = time.Now()
}

func (r *StrategyRule) UpdateAlertLevel(level vobj.AlertLevel) {
	r.alertLevel = level
	r.updatedAt = time.Now()
}

func (r *StrategyRule) UpdateAlertPages(pages []string) {
	r.alertPages = pages
	r.updatedAt = time.Now()
}

func (r *StrategyRule) UpdateOrder(order int32) {
	r.order = order
	r.updatedAt = time.Now()
}

// FromModel creates a StrategyRule entity from repository model
func StrategyRuleFromModel(uid, strategyUID snowflake.ID, ruleDetail string, status vobj.GlobalStatus, alertLevel vobj.AlertLevel, alertPages []string, order int32, createdAt, updatedAt time.Time) *StrategyRule {
	return &StrategyRule{
		uid:         uid,
		strategyUID: strategyUID,
		ruleDetail:  ruleDetail,
		status:      status,
		alertLevel:  alertLevel,
		alertPages:  alertPages,
		order:       order,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// Getters
func (r *StrategyRule) UID() snowflake.ID           { return r.uid }
func (r *StrategyRule) StrategyUID() snowflake.ID   { return r.strategyUID }
func (r *StrategyRule) RuleDetail() string         { return r.ruleDetail }
func (r *StrategyRule) Status() vobj.GlobalStatus  { return r.status }
func (r *StrategyRule) AlertLevel() vobj.AlertLevel { return r.alertLevel }
func (r *StrategyRule) AlertPages() []string        { return r.alertPages }
func (r *StrategyRule) Order() int32                { return r.order }
func (r *StrategyRule) CreatedAt() time.Time        { return r.createdAt }
func (r *StrategyRule) UpdatedAt() time.Time        { return r.updatedAt }

