package namespace

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/shared"
)

type Repository interface {
	Save(ctx context.Context, ns *Namespace) error
	FindByID(ctx context.Context, uid snowflake.ID) (*Namespace, error)
	FindByName(ctx context.Context, name string) (*Namespace, error)
	Delete(ctx context.Context, uid snowflake.ID) error
	List(ctx context.Context, query *ListQuery) (*shared.Page[*Namespace], error)
	Select(ctx context.Context, query *SelectQuery) (*SelectResult, error)
}

