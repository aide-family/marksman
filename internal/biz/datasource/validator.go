package datasource

import (
	"context"

	"github.com/aide-family/sovereign/pkg/merr"
)

type Validator struct {
	repo Repository
}

func NewValidator(repo Repository) *Validator {
	return &Validator{repo: repo}
}

func (v *Validator) ValidateConnection(ctx context.Context, ds *DataSource) error {
	if ds.Endpoint() == "" {
		return merr.ErrorParams("endpoint is required")
	}
	// TODO: 实际连接测试逻辑，根据类型和引擎实现
	return nil
}

func (v *Validator) ValidateConfig(ctx context.Context, typ Type, config map[string]string) error {
	if typ != TypeMetric && typ != TypeLogs && typ != TypeTrace && typ != TypeEvent {
		return merr.ErrorParams("invalid datasource type: %s", typ)
	}
	return nil
}
