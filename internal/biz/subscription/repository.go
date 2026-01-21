package subscription

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/shared"
)

type Repository interface {
	Save(ctx context.Context, s *Subscription) error
	FindByID(ctx context.Context, uid snowflake.ID) (*Subscription, error)
	Delete(ctx context.Context, uid snowflake.ID) error
	List(ctx context.Context, query *ListQuery) (*shared.Page[*Subscription], error)
	Select(ctx context.Context, query *SelectQuery) (*SelectResult, error)
}

