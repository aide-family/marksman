// Package repo is the repository service implementation.
package repo

import (
	"sync"

	datasourcemodel "github.com/aide-family/sovereign/pkg/repo/datasource/v1/gormimpl/model"
	"github.com/aide-family/sovereign/pkg/repo/gormmodel"
	namespacemodel "github.com/aide-family/sovereign/pkg/repo/namespace/v1/gormimpl/model"
)

var registerModelsOnce sync.Once

// RegisterAllModels registers all known models in one place.
func RegisterAllModels() {
	registerModelsOnce.Do(func() {
		gormmodel.RegisterModels(
			&namespacemodel.Namespace{},
			&namespacemodel.NamespaceModel{},
			&datasourcemodel.DataSource{},
			&datasourcemodel.DataSourceMetadata{},
			&datasourcemodel.DataSourceProxy{},
			&datasourcemodel.DataSourceTest{},
			&datasourcemodel.DataSourceQueryHistory{},
		)
	})
}
