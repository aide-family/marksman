package namespace

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/pkg/merr"
)

type Validator struct {
	repo Repository
}

func NewValidator(repo Repository) *Validator {
	return &Validator{repo: repo}
}

func (v *Validator) ValidateUnique(ctx context.Context, name string) error {
	_, err := v.repo.FindByName(ctx, name)
	if err == nil {
		return merr.ErrorParams("namespace %s already exists", name)
	}
	if !merr.IsNotFound(err) {
		return merr.ErrorInternal("check namespace exists failed").WithCause(err)
	}
	return nil
}

func (v *Validator) ValidateUniqueForUpdate(ctx context.Context, name string, excludeUID snowflake.ID) error {
	existing, err := v.repo.FindByName(ctx, name)
	if err == nil && existing.UID() != excludeUID {
		return merr.ErrorParams("namespace %s already exists", name)
	}
	if err != nil && !merr.IsNotFound(err) {
		return merr.ErrorInternal("check namespace exists failed").WithCause(err)
	}
	return nil
}

