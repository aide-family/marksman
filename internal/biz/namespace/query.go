package namespace

import (
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/shared"
)

type ListQuery struct {
	*shared.PageRequest
	Keyword string
	Status  Status
}

type SelectQuery struct {
	Keyword string
	Limit   int32
	NextUID snowflake.ID
	Status  Status
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

