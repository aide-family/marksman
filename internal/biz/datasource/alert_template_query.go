package datasource

import (
	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/internal/biz/vobj"
	"github.com/bwmarrin/snowflake"
)

// AlertTemplateListQuery 告警模板列表查询
type AlertTemplateListQuery struct {
	*shared.PageRequest
	DataSourceUID snowflake.ID
	Keyword       string
	Status        vobj.GlobalStatus
}

func NewAlertTemplateListQuery(datasourceUID snowflake.ID, page, pageSize int32, keyword string, status vobj.GlobalStatus) *AlertTemplateListQuery {
	return &AlertTemplateListQuery{
		PageRequest:   shared.NewPageRequest(page, pageSize),
		DataSourceUID: datasourceUID,
		Keyword:       keyword,
		Status:        status,
	}
}

