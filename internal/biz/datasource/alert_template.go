package datasource

import (
	"time"

	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/vobj"
	"github.com/aide-family/sovereign/pkg/merr"
)

// AlertTemplate 告警模板实体
type AlertTemplate struct {
	uid            snowflake.ID
	datasourceUID  snowflake.ID
	name           string
	titleTemplate  string
	contentTemplate string
	status         vobj.GlobalStatus
	createdAt      time.Time
	updatedAt      time.Time
}

func NewAlertTemplate(datasourceUID snowflake.ID, name, titleTemplate, contentTemplate string) *AlertTemplate {
	var uid snowflake.ID
	node, err := snowflake.NewNode(hello.NodeID())
	if err == nil {
		uid = node.Generate()
	}
	return &AlertTemplate{
		uid:             uid,
		datasourceUID:   datasourceUID,
		name:            name,
		titleTemplate:   titleTemplate,
		contentTemplate: contentTemplate,
		status:          vobj.GlobalStatusEnabled,
		createdAt:       time.Now(),
		updatedAt:       time.Now(),
	}
}

func (a *AlertTemplate) UpdateName(name string) error {
	if name == "" {
		return merr.ErrorParams("name cannot be empty")
	}
	a.name = name
	a.updatedAt = time.Now()
	return nil
}

func (a *AlertTemplate) UpdateTemplates(titleTemplate, contentTemplate string) {
	a.titleTemplate = titleTemplate
	a.contentTemplate = contentTemplate
	a.updatedAt = time.Now()
}

func (a *AlertTemplate) Enable() error {
	if a.status == vobj.GlobalStatusEnabled {
		return merr.ErrorParams("alert template already enabled")
	}
	a.status = vobj.GlobalStatusEnabled
	a.updatedAt = time.Now()
	return nil
}

func (a *AlertTemplate) Disable() error {
	a.status = vobj.GlobalStatusDisabled
	a.updatedAt = time.Now()
	return nil
}

func (a *AlertTemplate) UpdateStatus(status vobj.GlobalStatus) error {
	if status == vobj.GlobalStatusEnabled {
		return a.Enable()
	}
	return a.Disable()
}

// FromModel creates an AlertTemplate entity from repository model
func AlertTemplateFromModel(uid, datasourceUID snowflake.ID, name, titleTemplate, contentTemplate string, status vobj.GlobalStatus, createdAt, updatedAt time.Time) *AlertTemplate {
	return &AlertTemplate{
		uid:             uid,
		datasourceUID:   datasourceUID,
		name:            name,
		titleTemplate:   titleTemplate,
		contentTemplate: contentTemplate,
		status:          status,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}

// Getters
func (a *AlertTemplate) UID() snowflake.ID              { return a.uid }
func (a *AlertTemplate) DataSourceUID() snowflake.ID    { return a.datasourceUID }
func (a *AlertTemplate) Name() string                    { return a.name }
func (a *AlertTemplate) TitleTemplate() string           { return a.titleTemplate }
func (a *AlertTemplate) ContentTemplate() string         { return a.contentTemplate }
func (a *AlertTemplate) Status() vobj.GlobalStatus      { return a.status }
func (a *AlertTemplate) CreatedAt() time.Time            { return a.createdAt }
func (a *AlertTemplate) UpdatedAt() time.Time            { return a.updatedAt }

