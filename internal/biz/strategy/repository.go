package strategy

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/internal/biz/vobj"
)

type Repository interface {
	Save(ctx context.Context, s *Strategy) error
	FindByID(ctx context.Context, uid snowflake.ID) (*Strategy, error)
	Delete(ctx context.Context, uid snowflake.ID) error
	List(ctx context.Context, query *ListQuery) (*shared.Page[*Strategy], error)
	Select(ctx context.Context, query *SelectQuery) (*SelectResult, error)
}

type GroupRepository interface {
	Save(ctx context.Context, g *StrategyGroup) error
	FindByID(ctx context.Context, uid snowflake.ID) (*StrategyGroup, error)
	Delete(ctx context.Context, uid snowflake.ID) error
	List(ctx context.Context, query *GroupListQuery) (*shared.Page[*StrategyGroup], error)
}

type ReceiverRepository interface {
	Save(ctx context.Context, r *Receiver) error
	FindByID(ctx context.Context, uid snowflake.ID) (*Receiver, error)
	Delete(ctx context.Context, uid snowflake.ID) error
	List(ctx context.Context, query *ReceiverListQuery) (*shared.Page[*Receiver], error)
}

type StrategyRuleRepository interface {
	Save(ctx context.Context, rule *StrategyRule) error
	FindByID(ctx context.Context, uid snowflake.ID) (*StrategyRule, error)
	Delete(ctx context.Context, uid snowflake.ID) error
	List(ctx context.Context, strategyUID snowflake.ID, status vobj.GlobalStatus) ([]*StrategyRule, error)
}

// StrategyReceiverRepository 管理策略与接收对象的关联
type StrategyReceiverRepository interface {
	SaveStrategyReceivers(ctx context.Context, strategyUID snowflake.ID, receiverUIDs map[snowflake.ID]bool) error
	GetStrategyReceivers(ctx context.Context, strategyUID snowflake.ID) (map[snowflake.ID]bool, error)
}

