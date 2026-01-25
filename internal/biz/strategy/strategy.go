package strategy

import (
	"time"
	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/vobj"
	"github.com/aide-family/sovereign/pkg/merr"
)

type Strategy struct {
	uid            snowflake.ID
	namespaceUID   snowflake.ID
	groupUID       snowflake.ID
	typ            vobj.StrategyType
	name           string
	description    string
	status         vobj.GlobalStatus
	dataSourceUIDs map[snowflake.ID]bool
	query          string
	dataSourceType string
	dialTestType   vobj.DialTestType
	dialTestTargets map[string]string
	suppressType   string
	suppressConfig string
	alertTitle     string
	alertContent   string
	alertLevel     vobj.AlertLevel
	alertPages     []string
	rules          string
	labels         map[string]string
	metadata       map[string]string
	createdAt      time.Time
	updatedAt      time.Time
}

func New(namespaceUID snowflake.ID, typ vobj.StrategyType, name string) *Strategy {
	var uid snowflake.ID
	node, err := snowflake.NewNode(hello.NodeID())
	if err == nil {
		uid = node.Generate()
	}
	return &Strategy{
		uid:            uid,
		namespaceUID:   namespaceUID,
		typ:            typ,
		name:           name,
		status:         vobj.GlobalStatusEnabled,
		dataSourceUIDs: make(map[snowflake.ID]bool),
		dialTestTargets: make(map[string]string),
		alertPages:     make([]string, 0),
		labels:         make(map[string]string),
		metadata:       make(map[string]string),
		createdAt:      time.Now(),
		updatedAt:      time.Now(),
	}
}

func (s *Strategy) Enable() error {
	if s.status == vobj.GlobalStatusEnabled {
		return merr.ErrorParams("strategy already enabled")
	}
	s.status = vobj.GlobalStatusEnabled
	s.updatedAt = time.Now()
	return nil
}

func (s *Strategy) Disable() error {
	s.status = vobj.GlobalStatusDisabled
	s.updatedAt = time.Now()
	return nil
}

func (s *Strategy) IsEnabled() bool {
	return s.status == vobj.GlobalStatusEnabled
}

func (s *Strategy) UpdateName(name string) error {
	if name == "" {
		return merr.ErrorParams("name cannot be empty")
	}
	s.name = name
	s.updatedAt = time.Now()
	return nil
}

func (s *Strategy) UpdateDescription(description string) {
	s.description = description
	s.updatedAt = time.Now()
}

func (s *Strategy) UpdateGroupUID(groupUID snowflake.ID) {
	s.groupUID = groupUID
	s.updatedAt = time.Now()
}

func (s *Strategy) UpdateAlertLevel(level vobj.AlertLevel) {
	s.alertLevel = level
	s.updatedAt = time.Now()
}

func (s *Strategy) UpdateMetadata(metadata map[string]string) {
	s.metadata = metadata
	s.updatedAt = time.Now()
}

func (s *Strategy) UpdateDataSourceUIDs(uids map[snowflake.ID]bool) {
	s.dataSourceUIDs = uids
	s.updatedAt = time.Now()
}

func (s *Strategy) UpdateQuery(query string) {
	s.query = query
	s.updatedAt = time.Now()
}

func (s *Strategy) UpdateDataSourceType(dataSourceType string) {
	s.dataSourceType = dataSourceType
	s.updatedAt = time.Now()
}

func (s *Strategy) UpdateAlertTitle(title string) {
	s.alertTitle = title
	s.updatedAt = time.Now()
}

func (s *Strategy) UpdateAlertContent(content string) {
	s.alertContent = content
	s.updatedAt = time.Now()
}

func (s *Strategy) UpdateAlertPages(pages []string) {
	s.alertPages = pages
	s.updatedAt = time.Now()
}

func (s *Strategy) UpdateDialTestType(dialTestType vobj.DialTestType) {
	s.dialTestType = dialTestType
	s.updatedAt = time.Now()
}

func (s *Strategy) UpdateDialTestTargets(targets map[string]string) {
	s.dialTestTargets = targets
	s.updatedAt = time.Now()
}

func (s *Strategy) UpdateSuppressType(suppressType string) {
	s.suppressType = suppressType
	s.updatedAt = time.Now()
}

func (s *Strategy) UpdateSuppressConfig(config string) {
	s.suppressConfig = config
	s.updatedAt = time.Now()
}

// FromModel creates a Strategy entity from repository model
func FromModel(uid, namespaceUID, groupUID snowflake.ID, typ vobj.StrategyType, name, description string, status vobj.GlobalStatus, dataSourceUIDs map[snowflake.ID]bool, query, dataSourceType string, dialTestType vobj.DialTestType, dialTestTargets map[string]string, suppressType, suppressConfig, alertTitle, alertContent string, alertLevel vobj.AlertLevel, alertPages []string, rules string, labels, metadata map[string]string, createdAt, updatedAt time.Time) *Strategy {
	return &Strategy{
		uid:             uid,
		namespaceUID:    namespaceUID,
		groupUID:        groupUID,
		typ:             typ,
		name:            name,
		description:     description,
		status:          status,
		dataSourceUIDs:  dataSourceUIDs,
		query:           query,
		dataSourceType:  dataSourceType,
		dialTestType:    dialTestType,
		dialTestTargets: dialTestTargets,
		suppressType:    suppressType,
		suppressConfig:  suppressConfig,
		alertTitle:      alertTitle,
		alertContent:    alertContent,
		alertLevel:      alertLevel,
		alertPages:      alertPages,
		rules:           rules,
		labels:          labels,
		metadata:        metadata,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}

// Getters
func (s *Strategy) UID() snowflake.ID                    { return s.uid }
func (s *Strategy) NamespaceUID() snowflake.ID           { return s.namespaceUID }
func (s *Strategy) GroupUID() snowflake.ID               { return s.groupUID }
func (s *Strategy) Type() vobj.StrategyType              { return s.typ }
func (s *Strategy) Name() string                          { return s.name }
func (s *Strategy) Description() string                   { return s.description }
func (s *Strategy) Status() vobj.GlobalStatus            { return s.status }
func (s *Strategy) DataSourceUIDs() map[snowflake.ID]bool { return s.dataSourceUIDs }
func (s *Strategy) Query() string                         { return s.query }
func (s *Strategy) DataSourceType() string                { return s.dataSourceType }
func (s *Strategy) DialTestType() vobj.DialTestType      { return s.dialTestType }
func (s *Strategy) DialTestTargets() map[string]string    { return s.dialTestTargets }
func (s *Strategy) SuppressType() string                   { return s.suppressType }
func (s *Strategy) SuppressConfig() string                { return s.suppressConfig }
func (s *Strategy) AlertTitle() string                     { return s.alertTitle }
func (s *Strategy) AlertContent() string                  { return s.alertContent }
func (s *Strategy) AlertLevel() vobj.AlertLevel           { return s.alertLevel }
func (s *Strategy) AlertPages() []string                  { return s.alertPages }
func (s *Strategy) Rules() string                          { return s.rules }
func (s *Strategy) Labels() map[string]string             { return s.labels }
func (s *Strategy) Metadata() map[string]string            { return s.metadata }
func (s *Strategy) CreatedAt() time.Time                  { return s.createdAt }
func (s *Strategy) UpdatedAt() time.Time                  { return s.updatedAt }

