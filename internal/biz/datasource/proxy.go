package datasource

import (
	"time"

	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/pkg/merr"
)

// DataSourceProxy 数据源代理实体
type DataSourceProxy struct {
	uid            snowflake.ID
	namespaceUID   snowflake.ID
	datasourceUID  snowflake.ID
	typ            string
	name           string
	config         map[string]string
	createdAt      time.Time
	updatedAt      time.Time
}

func NewDataSourceProxy(namespaceUID, datasourceUID snowflake.ID, typ, name string, config map[string]string) *DataSourceProxy {
	var uid snowflake.ID
	node, err := snowflake.NewNode(hello.NodeID())
	if err == nil {
		uid = node.Generate()
	}
	return &DataSourceProxy{
		uid:           uid,
		namespaceUID:  namespaceUID,
		datasourceUID: datasourceUID,
		typ:           typ,
		name:          name,
		config:        config,
		createdAt:     time.Now(),
		updatedAt:     time.Now(),
	}
}

func (p *DataSourceProxy) UpdateName(name string) error {
	if name == "" {
		return merr.ErrorParams("name cannot be empty")
	}
	p.name = name
	p.updatedAt = time.Now()
	return nil
}

func (p *DataSourceProxy) UpdateConfig(config map[string]string) {
	p.config = config
	p.updatedAt = time.Now()
}

// FromModel creates a DataSourceProxy entity from repository model
func DataSourceProxyFromModel(uid, namespaceUID, datasourceUID snowflake.ID, typ, name string, config map[string]string, createdAt, updatedAt time.Time) *DataSourceProxy {
	return &DataSourceProxy{
		uid:           uid,
		namespaceUID:  namespaceUID,
		datasourceUID: datasourceUID,
		typ:           typ,
		name:          name,
		config:        config,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

// Getters
func (p *DataSourceProxy) UID() snowflake.ID              { return p.uid }
func (p *DataSourceProxy) NamespaceUID() snowflake.ID      { return p.namespaceUID }
func (p *DataSourceProxy) DataSourceUID() snowflake.ID     { return p.datasourceUID }
func (p *DataSourceProxy) Type() string                    { return p.typ }
func (p *DataSourceProxy) Name() string                    { return p.name }
func (p *DataSourceProxy) Config() map[string]string       { return p.config }
func (p *DataSourceProxy) CreatedAt() time.Time            { return p.createdAt }
func (p *DataSourceProxy) UpdatedAt() time.Time            { return p.updatedAt }

