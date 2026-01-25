package namespace

import (
	"time"
	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/pkg/merr"
)

type Namespace struct {
	uid       snowflake.ID
	name      string
	metadata  map[string]string
	status    Status
	createdAt time.Time
	updatedAt time.Time
}

func New(name string, metadata map[string]string) *Namespace {
	var uid snowflake.ID
	node, err := snowflake.NewNode(hello.NodeID())
	if err == nil {
		uid = node.Generate()
	}
	return &Namespace{
		uid:       uid,
		name:      name,
		metadata:  metadata,
		status:    StatusEnabled,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
}

func (n *Namespace) Enable() error {
	if n.status == StatusEnabled {
		return merr.ErrorParams("namespace already enabled")
	}
	n.status = StatusEnabled
	n.updatedAt = time.Now()
	return nil
}

func (n *Namespace) Disable() error {
	n.status = StatusDisabled
	n.updatedAt = time.Now()
	return nil
}

func (n *Namespace) IsEnabled() bool {
	return n.status == StatusEnabled
}

func (n *Namespace) UpdateName(name string) error {
	if name == "" {
		return merr.ErrorParams("name cannot be empty")
	}
	n.name = name
	n.updatedAt = time.Now()
	return nil
}

func (n *Namespace) UpdateMetadata(metadata map[string]string) {
	n.metadata = metadata
	n.updatedAt = time.Now()
}

// FromModel creates a Namespace entity from repository model
func FromModel(uid snowflake.ID, name string, metadata map[string]string, status Status, createdAt, updatedAt time.Time) *Namespace {
	return &Namespace{
		uid:       uid,
		name:      name,
		metadata:  metadata,
		status:    status,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// Getters
func (n *Namespace) UID() snowflake.ID        { return n.uid }
func (n *Namespace) Name() string             { return n.name }
func (n *Namespace) Metadata() map[string]string { return n.metadata }
func (n *Namespace) Status() Status           { return n.status }
func (n *Namespace) CreatedAt() time.Time    { return n.createdAt }
func (n *Namespace) UpdatedAt() time.Time     { return n.updatedAt }

