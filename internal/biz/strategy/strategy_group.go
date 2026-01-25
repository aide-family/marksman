package strategy

import (
	"time"
	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/vobj"
	"github.com/aide-family/sovereign/pkg/merr"
)

type StrategyGroup struct {
	uid          snowflake.ID
	namespaceUID snowflake.ID
	name         string
	description  string
	status       vobj.GlobalStatus
	upgradeMode  vobj.UpgradeMode
	upgradeConfig string
	createdAt    time.Time
	updatedAt    time.Time
}

func NewGroup(namespaceUID snowflake.ID, name string) *StrategyGroup {
	var uid snowflake.ID
	node, err := snowflake.NewNode(hello.NodeID())
	if err == nil {
		uid = node.Generate()
	}
	return &StrategyGroup{
		uid:          uid,
		namespaceUID: namespaceUID,
		name:         name,
		status:       vobj.GlobalStatusEnabled,
		upgradeMode:  vobj.UpgradeModeNone,
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
	}
}

func (g *StrategyGroup) Enable() error {
	if g.status == vobj.GlobalStatusEnabled {
		return merr.ErrorParams("strategy group already enabled")
	}
	g.status = vobj.GlobalStatusEnabled
	g.updatedAt = time.Now()
	return nil
}

func (g *StrategyGroup) Disable() error {
	g.status = vobj.GlobalStatusDisabled
	g.updatedAt = time.Now()
	return nil
}

func (g *StrategyGroup) IsEnabled() bool {
	return g.status == vobj.GlobalStatusEnabled
}

func (g *StrategyGroup) UpdateName(name string) error {
	if name == "" {
		return merr.ErrorParams("name cannot be empty")
	}
	g.name = name
	g.updatedAt = time.Now()
	return nil
}

func (g *StrategyGroup) UpdateDescription(description string) {
	g.description = description
	g.updatedAt = time.Now()
}

func (g *StrategyGroup) UpdateUpgradeMode(mode vobj.UpgradeMode) {
	g.upgradeMode = mode
	g.updatedAt = time.Now()
}

func (g *StrategyGroup) UpdateUpgradeConfig(config string) {
	g.upgradeConfig = config
	g.updatedAt = time.Now()
}

// FromModel creates a StrategyGroup entity from repository model
func GroupFromModel(uid, namespaceUID snowflake.ID, name, description string, status vobj.GlobalStatus, upgradeMode vobj.UpgradeMode, upgradeConfig string, createdAt, updatedAt time.Time) *StrategyGroup {
	return &StrategyGroup{
		uid:          uid,
		namespaceUID: namespaceUID,
		name:         name,
		description:  description,
		status:       status,
		upgradeMode:  upgradeMode,
		upgradeConfig: upgradeConfig,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Getters
func (g *StrategyGroup) UID() snowflake.ID           { return g.uid }
func (g *StrategyGroup) NamespaceUID() snowflake.ID  { return g.namespaceUID }
func (g *StrategyGroup) Name() string                 { return g.name }
func (g *StrategyGroup) Description() string         { return g.description }
func (g *StrategyGroup) Status() vobj.GlobalStatus   { return g.status }
func (g *StrategyGroup) UpgradeMode() vobj.UpgradeMode { return g.upgradeMode }
func (g *StrategyGroup) UpgradeConfig() string        { return g.upgradeConfig }
func (g *StrategyGroup) CreatedAt() time.Time         { return g.createdAt }
func (g *StrategyGroup) UpdatedAt() time.Time         { return g.updatedAt }

