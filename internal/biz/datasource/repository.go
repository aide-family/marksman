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

// AlertTemplateRepository 告警模板仓库接口
type AlertTemplateRepository interface {
	Save(ctx context.Context, template *AlertTemplate) error
	FindByID(ctx context.Context, uid snowflake.ID) (*AlertTemplate, error)
	Delete(ctx context.Context, uid snowflake.ID) error
	List(ctx context.Context, query *AlertTemplateListQuery) (*shared.Page[*AlertTemplate], error)
}

// QueryHistoryRepository 查询历史仓库接口
type QueryHistoryRepository interface {
	Save(ctx context.Context, history *QueryHistory) error
	List(ctx context.Context, datasourceUID snowflake.ID, page, pageSize int32) ([]*QueryHistory, int64, error)
}

// QueryFavoriteRepository 查询收藏仓库接口
type QueryFavoriteRepository interface {
	Save(ctx context.Context, favorite *QueryFavorite) error
	FindByID(ctx context.Context, id uint32) (*QueryFavorite, error)
	Delete(ctx context.Context, id uint32) error
	List(ctx context.Context, datasourceUID snowflake.ID) ([]*QueryFavorite, error)
}

// DataSourceMetadataRepository 数据源元数据仓库接口
type DataSourceMetadataRepository interface {
	Save(ctx context.Context, metadata *DataSourceMetadata) error
	FindByKey(ctx context.Context, datasourceUID snowflake.ID, key string) (*DataSourceMetadata, error)
	List(ctx context.Context, datasourceUID snowflake.ID, metadataType string) ([]*DataSourceMetadata, error)
	DeleteByDataSourceUID(ctx context.Context, datasourceUID snowflake.ID) error
}

// DataSourceProxyRepository 数据源代理仓库接口
type DataSourceProxyRepository interface {
	Save(ctx context.Context, proxy *DataSourceProxy) error
	FindByID(ctx context.Context, uid snowflake.ID) (*DataSourceProxy, error)
	Delete(ctx context.Context, uid snowflake.ID) error
	List(ctx context.Context, namespaceUID snowflake.ID, page, pageSize int32) ([]*DataSourceProxy, int64, error)
}

