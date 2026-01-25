package strategy

import (
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/internal/biz/vobj"
)

type ListQuery struct {
	*shared.PageRequest
	NamespaceUID snowflake.ID
	GroupUID     snowflake.ID
	Type         vobj.StrategyType
	Keyword      string
	Status       vobj.GlobalStatus
}

type SelectQuery struct {
	NamespaceUID snowflake.ID
	Keyword      string
	Limit        int32
	NextUID      snowflake.ID
	Status       vobj.GlobalStatus
}

type SelectResult struct {
	Items   []*SelectItem
	Total   int64
	NextUID snowflake.ID
	HasMore bool
}

type SelectItem struct {
	UID      snowflake.ID
	Name     string
	Disabled bool
	Tooltip  string
}

type GroupListQuery struct {
	*shared.PageRequest
	NamespaceUID snowflake.ID
	Keyword      string
	Status       vobj.GlobalStatus
}

type ReceiverListQuery struct {
	*shared.PageRequest
	NamespaceUID snowflake.ID
	Type         vobj.ReceiverType
	Keyword      string
}

