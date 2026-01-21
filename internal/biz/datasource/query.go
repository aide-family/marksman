package datasource

import (
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/shared"
)

type ListQuery struct {
	*shared.PageRequest
	Keyword      string
	Status       Status
	Type         Type
	Engine       Engine
	NamespaceUID snowflake.ID
	OrderBy      int32
	Order        int32
}

type SelectQuery struct {
	Keyword      string
	Limit        int32
	NextUID      snowflake.ID
	Status       Status
	Type         Type
	Engine       Engine
	NamespaceUID snowflake.ID
	Order        int32
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

