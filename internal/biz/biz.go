// Package biz is the business logic for the Sovereign service.
package biz

import (
	"github.com/google/wire"
	"github.com/aide-family/sovereign/internal/biz/datasource"
	"github.com/aide-family/sovereign/internal/biz/namespace"
	"github.com/aide-family/sovereign/internal/biz/strategy"
	"github.com/aide-family/sovereign/internal/biz/subscription"
)

var ProviderSetBiz = wire.NewSet(
	NewHealth,
	namespace.NewService,
	namespace.NewValidator,
	datasource.NewService,
	datasource.NewValidator,
	strategy.NewService,
	strategy.NewValidator,
	subscription.NewService,
	subscription.NewValidator,
)
