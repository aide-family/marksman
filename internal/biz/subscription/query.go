package subscription

import (
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/internal/biz/vobj"
)

type ListQuery struct {
	*shared.PageRequest
	UserID      snowflake.ID
	NamespaceUID snowflake.ID
	Type        SubscriptionType
	Keyword     string
	Status      vobj.GlobalStatus
}

type SelectQuery struct {
	UserID      snowflake.ID
	NamespaceUID snowflake.ID
	Keyword     string
	Limit       int32
	NextUID     snowflake.ID
	Status      vobj.GlobalStatus
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

