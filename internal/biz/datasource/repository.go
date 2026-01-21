package datasource

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/shared"
)

type Repository interface {
	Save(ctx context.Context, ds *DataSource) error
	FindByID(ctx context.Context, uid snowflake.ID) (*DataSource, error)
	Delete(ctx context.Context, uid snowflake.ID) error
	List(ctx context.Context, query *ListQuery) (*shared.Page[*DataSource], error)
	Select(ctx context.Context, query *SelectQuery) (*SelectResult, error)
}

