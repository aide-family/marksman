// Package gormmodel provides shared helpers for GORM models.
package gormmodel

import "sync"

var (
	modelsMu sync.Mutex
	models   []any
)

// RegisterModels registers models for shared migrations or codegen.
func RegisterModels(values ...any) {
	modelsMu.Lock()
	defer modelsMu.Unlock()
	models = append(models, values...)
}

// Models returns a copy of registered models.
func Models() []any {
	modelsMu.Lock()
	defer modelsMu.Unlock()
	out := make([]any, len(models))
	copy(out, models)
	return out
}
