package datasource

import (
	"time"
	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/pkg/merr"
)

type DataSource struct {
	uid          snowflake.ID
	namespaceUID snowflake.ID
	typ          Type
	engine       Engine
	name         string
	status       Status
	endpoint     string
	description  string
	config       map[string]string
	metadata     map[string]string
	createdAt    time.Time
	updatedAt    time.Time
}

func New(namespaceUID snowflake.ID, typ Type, engine Engine, name, endpoint, description string, config, metadata map[string]string) *DataSource {
	var uid snowflake.ID
	node, err := snowflake.NewNode(hello.NodeID())
	if err == nil {
		uid = node.Generate()
	}
	return &DataSource{
		uid:          uid,
		namespaceUID: namespaceUID,
		typ:          typ,
		engine:       engine,
		name:         name,
		status:       StatusEnabled,
		endpoint:     endpoint,
		description:  description,
		config:       config,
		metadata:     metadata,
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
	}
}

func (d *DataSource) Enable() error {
	if d.status == StatusEnabled {
		return merr.ErrorParams("datasource already enabled")
	}
	d.status = StatusEnabled
	d.updatedAt = time.Now()
	return nil
}

func (d *DataSource) Disable() error {
	d.status = StatusDisabled
	d.updatedAt = time.Now()
	return nil
}

func (d *DataSource) IsEnabled() bool {
	return d.status == StatusEnabled
}

func (d *DataSource) UpdateEndpoint(endpoint string) error {
	if endpoint == "" {
		return merr.ErrorParams("endpoint cannot be empty")
	}
	d.endpoint = endpoint
	d.updatedAt = time.Now()
	return nil
}

func (d *DataSource) UpdateConfig(config map[string]string) {
	d.config = config
	d.updatedAt = time.Now()
}

func (d *DataSource) UpdateName(name string) error {
	if name == "" {
		return merr.ErrorParams("name cannot be empty")
	}
	d.name = name
	d.updatedAt = time.Now()
	return nil
}

func (d *DataSource) UpdateDescription(description string) {
	d.description = description
	d.updatedAt = time.Now()
}

func (d *DataSource) UpdateMetadata(metadata map[string]string) {
	d.metadata = metadata
	d.updatedAt = time.Now()
}

// FromModel creates a DataSource entity from repository model
func FromModel(uid, namespaceUID snowflake.ID, typ Type, engine Engine, name string, status Status, endpoint, description string, config, metadata map[string]string, createdAt, updatedAt time.Time) *DataSource {
	return &DataSource{
		uid:          uid,
		namespaceUID: namespaceUID,
		typ:          typ,
		engine:       engine,
		name:         name,
		status:       status,
		endpoint:     endpoint,
		description:  description,
		config:       config,
		metadata:     metadata,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Getters
func (d *DataSource) UID() snowflake.ID          { return d.uid }
func (d *DataSource) NamespaceUID() snowflake.ID { return d.namespaceUID }
func (d *DataSource) Type() Type                 { return d.typ }
func (d *DataSource) Engine() Engine              { return d.engine }
func (d *DataSource) Name() string                { return d.name }
func (d *DataSource) Status() Status              { return d.status }
func (d *DataSource) Endpoint() string            { return d.endpoint }
func (d *DataSource) Description() string         { return d.description }
func (d *DataSource) Config() map[string]string   { return d.config }
func (d *DataSource) Metadata() map[string]string { return d.metadata }
func (d *DataSource) CreatedAt() time.Time        { return d.createdAt }
func (d *DataSource) UpdatedAt() time.Time       { return d.updatedAt }

