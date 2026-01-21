package subscription

import (
	"time"
	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/vobj"
	"github.com/aide-family/sovereign/pkg/merr"
)

type SubscriptionType string

const (
	SubscriptionTypeStrategyGroup SubscriptionType = "strategy_group" // 按策略组订阅
	SubscriptionTypeDataSource    SubscriptionType = "datasource"     // 按数据源订阅
)

func (t SubscriptionType) String() string {
	return string(t)
}

type Subscription struct {
	uid              snowflake.ID
	userID           snowflake.ID
	namespaceUID     snowflake.ID
	typ              SubscriptionType
	name             string
	description      string
	status           vobj.GlobalStatus
	strategyGroupUIDs map[snowflake.ID]bool
	dataSourceUIDs    map[snowflake.ID]bool
	alertLevels       []vobj.AlertLevel
	notifyTypes       []vobj.NotifyType
	metadata          map[string]string
	createdAt         time.Time
	updatedAt         time.Time
}

func New(userID, namespaceUID snowflake.ID, typ SubscriptionType, name string) *Subscription {
	var uid snowflake.ID
	node, err := snowflake.NewNode(hello.NodeID())
	if err == nil {
		uid = node.Generate()
	}
	return &Subscription{
		uid:               uid,
		userID:            userID,
		namespaceUID:      namespaceUID,
		typ:               typ,
		name:              name,
		status:            vobj.GlobalStatusEnabled,
		strategyGroupUIDs: make(map[snowflake.ID]bool),
		dataSourceUIDs:    make(map[snowflake.ID]bool),
		alertLevels:       make([]vobj.AlertLevel, 0),
		notifyTypes:       make([]vobj.NotifyType, 0),
		metadata:          make(map[string]string),
		createdAt:         time.Now(),
		updatedAt:         time.Now(),
	}
}

func (s *Subscription) Enable() error {
	if s.status == vobj.GlobalStatusEnabled {
		return merr.ErrorParams("subscription already enabled")
	}
	s.status = vobj.GlobalStatusEnabled
	s.updatedAt = time.Now()
	return nil
}

func (s *Subscription) Disable() error {
	s.status = vobj.GlobalStatusDisabled
	s.updatedAt = time.Now()
	return nil
}

func (s *Subscription) IsEnabled() bool {
	return s.status == vobj.GlobalStatusEnabled
}

func (s *Subscription) UpdateName(name string) error {
	if name == "" {
		return merr.ErrorParams("name cannot be empty")
	}
	s.name = name
	s.updatedAt = time.Now()
	return nil
}

func (s *Subscription) UpdateDescription(description string) {
	s.description = description
	s.updatedAt = time.Now()
}

func (s *Subscription) UpdateStrategyGroupUIDs(uids map[snowflake.ID]bool) {
	s.strategyGroupUIDs = uids
	s.updatedAt = time.Now()
}

func (s *Subscription) UpdateDataSourceUIDs(uids map[snowflake.ID]bool) {
	s.dataSourceUIDs = uids
	s.updatedAt = time.Now()
}

func (s *Subscription) UpdateAlertLevels(levels []vobj.AlertLevel) {
	s.alertLevels = levels
	s.updatedAt = time.Now()
}

func (s *Subscription) UpdateNotifyTypes(types []vobj.NotifyType) {
	s.notifyTypes = types
	s.updatedAt = time.Now()
}

func (s *Subscription) UpdateMetadata(metadata map[string]string) {
	s.metadata = metadata
	s.updatedAt = time.Now()
}

// FromModel creates a Subscription entity from repository model
func FromModel(uid, userID, namespaceUID snowflake.ID, typ SubscriptionType, name, description string, status vobj.GlobalStatus, strategyGroupUIDs, dataSourceUIDs map[snowflake.ID]bool, alertLevels []vobj.AlertLevel, notifyTypes []vobj.NotifyType, metadata map[string]string, createdAt, updatedAt time.Time) *Subscription {
	return &Subscription{
		uid:               uid,
		userID:            userID,
		namespaceUID:      namespaceUID,
		typ:               typ,
		name:              name,
		description:      description,
		status:            status,
		strategyGroupUIDs: strategyGroupUIDs,
		dataSourceUIDs:    dataSourceUIDs,
		alertLevels:       alertLevels,
		notifyTypes:       notifyTypes,
		metadata:          metadata,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}
}

// Getters
func (s *Subscription) UID() snowflake.ID                    { return s.uid }
func (s *Subscription) UserID() snowflake.ID                 { return s.userID }
func (s *Subscription) NamespaceUID() snowflake.ID           { return s.namespaceUID }
func (s *Subscription) Type() SubscriptionType                { return s.typ }
func (s *Subscription) Name() string                          { return s.name }
func (s *Subscription) Description() string                   { return s.description }
func (s *Subscription) Status() vobj.GlobalStatus            { return s.status }
func (s *Subscription) StrategyGroupUIDs() map[snowflake.ID]bool { return s.strategyGroupUIDs }
func (s *Subscription) DataSourceUIDs() map[snowflake.ID]bool    { return s.dataSourceUIDs }
func (s *Subscription) AlertLevels() []vobj.AlertLevel          { return s.alertLevels }
func (s *Subscription) NotifyTypes() []vobj.NotifyType          { return s.notifyTypes }
func (s *Subscription) Metadata() map[string]string             { return s.metadata }
func (s *Subscription) CreatedAt() time.Time                    { return s.createdAt }
func (s *Subscription) UpdatedAt() time.Time                   { return s.updatedAt }

