package datasource

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

// DataSourceMetadata 数据源元数据实体
type DataSourceMetadata struct {
	id            uint32
	datasourceUID snowflake.ID
	key           string
	value         string
	metadataType  string // table/column/index等
	description   string
	createdAt     time.Time
	updatedAt     time.Time
}

func NewDataSourceMetadata(datasourceUID snowflake.ID, key, value, metadataType, description string) *DataSourceMetadata {
	return &DataSourceMetadata{
		datasourceUID: datasourceUID,
		key:           key,
		value:         value,
		metadataType:  metadataType,
		description:   description,
		createdAt:     time.Now(),
		updatedAt:     time.Now(),
	}
}

func (m *DataSourceMetadata) UpdateValue(value string) {
	m.value = value
	m.updatedAt = time.Now()
}

func (m *DataSourceMetadata) UpdateDescription(description string) {
	m.description = description
	m.updatedAt = time.Now()
}

// FromModel creates a DataSourceMetadata entity from repository model
func DataSourceMetadataFromModel(id uint32, datasourceUID snowflake.ID, key, value, metadataType, description string, createdAt, updatedAt time.Time) *DataSourceMetadata {
	return &DataSourceMetadata{
		id:            id,
		datasourceUID: datasourceUID,
		key:           key,
		value:         value,
		metadataType:  metadataType,
		description:   description,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

// Getters
func (m *DataSourceMetadata) ID() uint32              { return m.id }
func (m *DataSourceMetadata) DataSourceUID() snowflake.ID { return m.datasourceUID }
func (m *DataSourceMetadata) Key() string            { return m.key }
func (m *DataSourceMetadata) Value() string          { return m.value }
func (m *DataSourceMetadata) MetadataType() string   { return m.metadataType }
func (m *DataSourceMetadata) Description() string    { return m.description }
func (m *DataSourceMetadata) CreatedAt() time.Time   { return m.createdAt }
func (m *DataSourceMetadata) UpdatedAt() time.Time   { return m.updatedAt }

