package bo

import (
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/sovereign/internal/biz/vobj"
	apiv1 "github.com/aide-family/sovereign/pkg/api/v1"
	"github.com/aide-family/sovereign/pkg/enum"
)

type CreateDataSourceBo struct {
	NamespaceUID snowflake.ID
	Type         string
	Engine       string
	Name         string
	Status       vobj.GlobalStatus
	Endpoint     string
	Description  string
	Config       map[string]string
	Metadata     map[string]string
}

func NewCreateDataSourceBo(req *apiv1.CreateDataSourceRequest) *CreateDataSourceBo {
	return &CreateDataSourceBo{
		NamespaceUID: snowflake.ParseInt64(req.NamespaceUid),
		Type:         req.Type,
		Engine:       req.Engine,
		Name:         req.Name,
		Status:       vobj.GlobalStatusEnabled,
		Endpoint:     req.Endpoint,
		Description:  req.Description,
		Config:       req.Config,
		Metadata:     req.Metadata,
	}
}

type UpdateDataSourceBo struct {
	UID         snowflake.ID
	Name        string
	Status      vobj.GlobalStatus
	Endpoint    string
	Description string
	Config      map[string]string
	Metadata    map[string]string
}

func NewUpdateDataSourceBo(req *apiv1.UpdateDataSourceRequest) *UpdateDataSourceBo {
	return &UpdateDataSourceBo{
		UID:         snowflake.ParseInt64(req.Uid),
		Name:        req.Name,
		Status:      vobj.GlobalStatus(req.Status),
		Endpoint:    req.Endpoint,
		Description: req.Description,
		Config:      req.Config,
		Metadata:    req.Metadata,
	}
}

type UpdateDataSourceStatusBo struct {
	UID    snowflake.ID
	Status vobj.GlobalStatus
}

func NewUpdateDataSourceStatusBo(req *apiv1.UpdateDataSourceStatusRequest) *UpdateDataSourceStatusBo {
	return &UpdateDataSourceStatusBo{
		UID:    snowflake.ParseInt64(req.Uid),
		Status: vobj.GlobalStatus(req.Status),
	}
}

type ListDataSourceBo struct {
	*PageRequestBo
	Keyword      string
	Status       vobj.GlobalStatus
	Type         string
	Engine       string
	NamespaceUID snowflake.ID
	OrderBy      int32
	Order        int32
}

func NewListDataSourceBo(req *apiv1.ListDataSourceRequest) *ListDataSourceBo {
	return &ListDataSourceBo{
		PageRequestBo: NewPageRequestBo(req.Page, req.PageSize),
		Keyword:       req.Keyword,
		Status:        vobj.GlobalStatus(req.Status),
		Type:          req.Type,
		Engine:        req.Engine,
		NamespaceUID:  snowflake.ParseInt64(req.NamespaceUid),
	}
}

type DataSourceItemBo struct {
	UID          snowflake.ID
	NamespaceUID snowflake.ID
	Type         string
	Engine       string
	Name         string
	Status       vobj.GlobalStatus
	Endpoint     string
	Description  string
	Config       map[string]string
	Metadata     map[string]string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (b *DataSourceItemBo) ToAPIV1DataSourceItem() *apiv1.DataSourceItem {
	return &apiv1.DataSourceItem{
		Uid:          b.UID.Int64(),
		NamespaceUid: b.NamespaceUID.Int64(),
		Type:         b.Type,
		Engine:       b.Engine,
		Name:         b.Name,
		Status:       enum.GlobalStatus(b.Status),
		Endpoint:     b.Endpoint,
		Description:  b.Description,
		Config:       b.Config,
		Metadata:     b.Metadata,
		CreatedAt:    b.CreatedAt.Format(time.DateTime),
		UpdatedAt:    b.UpdatedAt.Format(time.DateTime),
	}
}

func ToAPIV1ListDataSourceReply(pageResponseBo *PageResponseBo[*DataSourceItemBo]) *apiv1.ListDataSourceReply {
	items := make([]*apiv1.DataSourceItem, 0, len(pageResponseBo.GetItems()))
	for _, item := range pageResponseBo.GetItems() {
		items = append(items, item.ToAPIV1DataSourceItem())
	}
	return &apiv1.ListDataSourceReply{
		Items:    items,
		Total:    pageResponseBo.GetTotal(),
		Page:     pageResponseBo.GetPage(),
		PageSize: pageResponseBo.GetPageSize(),
	}
}

type SelectDataSourceBo struct {
	Keyword      string
	Limit        int32
	LastUID      snowflake.ID
	Status       vobj.GlobalStatus
	Type         string
	Engine       string
	NamespaceUID snowflake.ID
	Order        int32
}

func NewSelectDataSourceBo(req *apiv1.SelectDataSourceRequest) *SelectDataSourceBo {
	var lastUID snowflake.ID
	if req.LastUID > 0 {
		lastUID = snowflake.ParseInt64(req.LastUID)
	}
	return &SelectDataSourceBo{
		Keyword:      req.Keyword,
		Limit:        req.Limit,
		LastUID:      lastUID,
		Status:       vobj.GlobalStatus(req.Status),
		Type:         req.Type,
		Engine:       req.Engine,
		NamespaceUID: snowflake.ParseInt64(req.NamespaceUid),
	}
}

type DataSourceItemSelectBo struct {
	UID      snowflake.ID
	Name     string
	Disabled bool
	Tooltip  string
}

func (b *DataSourceItemSelectBo) ToAPIV1DataSourceItemSelect() *apiv1.DataSourceItemSelect {
	return &apiv1.DataSourceItemSelect{
		Value:    b.UID.Int64(),
		Label:    b.Name,
		Disabled: b.Disabled,
		Tooltip:  b.Tooltip,
	}
}

type SelectDataSourceBoResult struct {
	Items   []*DataSourceItemSelectBo
	Total   int64
	LastUID snowflake.ID
	HasMore bool
}

func ToAPIV1SelectDataSourceReply(result *SelectDataSourceBoResult) *apiv1.SelectDataSourceReply {
	selectItems := make([]*apiv1.DataSourceItemSelect, 0, len(result.Items))
	for _, item := range result.Items {
		selectItems = append(selectItems, item.ToAPIV1DataSourceItemSelect())
	}
	return &apiv1.SelectDataSourceReply{
		Items:   selectItems,
		Total:   result.Total,
		LastUID: result.LastUID.Int64(),
		HasMore: result.HasMore,
	}
}
