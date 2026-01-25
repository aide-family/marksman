package datasource

import (
	"time"

	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/bwmarrin/snowflake"
)

// ListQuery 数据源列表查询
type ListQuery struct {
	*shared.PageRequest
	NamespaceUID snowflake.ID
	Keyword      string
	Status       Status
	Type         Type
	Engine       Engine
}

// SelectQuery 数据源选择查询
type SelectQuery struct {
	NamespaceUID snowflake.ID
	Keyword      string
	Status       Status
	Type         Type
	Engine       Engine
	NextUID      snowflake.ID
	Limit        int32
}

// SelectResult 数据源选择结果
type SelectResult struct {
	Items  []*SelectItem
	NextUID snowflake.ID
	HasMore bool
}

// SelectItem 数据源选择项
type SelectItem struct {
	UID      snowflake.ID
	Name     string
	Disabled bool
	Tooltip  string
}

// QueryHistory 查询历史实体
type QueryHistory struct {
	id            uint32
	datasourceUID snowflake.ID
	query         string
	format        string
	createdAt     time.Time
}

func NewQueryHistory(datasourceUID snowflake.ID, query, format string) *QueryHistory {
	return &QueryHistory{
		datasourceUID: datasourceUID,
		query:         query,
		format:        format,
		createdAt:     time.Now(),
	}
}

// FromModel creates a QueryHistory entity from repository model
func QueryHistoryFromModel(id uint32, datasourceUID snowflake.ID, query string, createdAt time.Time) *QueryHistory {
	return &QueryHistory{
		id:            id,
		datasourceUID: datasourceUID,
		query:         query,
		format:        "", // 历史记录中没有format字段
		createdAt:     createdAt,
	}
}

func (q *QueryHistory) ID() uint32              { return q.id }
func (q *QueryHistory) DataSourceUID() snowflake.ID { return q.datasourceUID }
func (q *QueryHistory) Query() string           { return q.query }
func (q *QueryHistory) Format() string        { return q.format }
func (q *QueryHistory) CreatedAt() time.Time    { return q.createdAt }

// QueryFavorite 查询收藏实体
type QueryFavorite struct {
	id        uint32
	datasourceUID snowflake.ID
	name      string
	query     string
	createdAt time.Time
}

func NewQueryFavorite(datasourceUID snowflake.ID, name, query string) *QueryFavorite {
	return &QueryFavorite{
		datasourceUID: datasourceUID,
		name:          name,
		query:         query,
		createdAt:     time.Now(),
	}
}

// FromModel creates a QueryFavorite entity from repository model
func QueryFavoriteFromModel(id uint32, datasourceUID snowflake.ID, name, query string, createdAt time.Time) *QueryFavorite {
	return &QueryFavorite{
		id:            id,
		datasourceUID: datasourceUID,
		name:          name,
		query:         query,
		createdAt:     createdAt,
	}
}

func (q *QueryFavorite) UpdateName(name string) {
	q.name = name
}

func (q *QueryFavorite) ID() uint32              { return q.id }
func (q *QueryFavorite) DataSourceUID() snowflake.ID { return q.datasourceUID }
func (q *QueryFavorite) Name() string            { return q.name }
func (q *QueryFavorite) Query() string           { return q.query }
func (q *QueryFavorite) CreatedAt() time.Time    { return q.createdAt }
