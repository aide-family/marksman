package strategy

import (
	"context"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/pkg/merr"
)

type Validator struct {
	repo      Repository
	groupRepo GroupRepository
	helper    *klog.Helper
}

func NewValidator(repo Repository, groupRepo GroupRepository, helper *klog.Helper) *Validator {
	return &Validator{
		repo:      repo,
		groupRepo: groupRepo,
		helper:    klog.NewHelper(klog.With(helper.Logger(), "biz", "strategy", "validator")),
	}
}

func (v *Validator) ValidateUnique(ctx context.Context, namespaceUID snowflake.ID, name string) error {
	query := &ListQuery{
		PageRequest: &shared.PageRequest{Page: 1, PageSize: 1},
		NamespaceUID: namespaceUID,
		Keyword: name,
	}
	page, err := v.repo.List(ctx, query)
	if err != nil {
		return merr.ErrorInternal("validate unique failed").WithCause(err)
	}
	if len(page.Items) > 0 {
		return merr.ErrorParams("strategy name %s already exists", name)
	}
	return nil
}

func (v *Validator) ValidateUniqueForUpdate(ctx context.Context, namespaceUID snowflake.ID, name string, uid snowflake.ID) error {
	query := &ListQuery{
		PageRequest: &shared.PageRequest{Page: 1, PageSize: 1},
		NamespaceUID: namespaceUID,
		Keyword: name,
	}
	page, err := v.repo.List(ctx, query)
	if err != nil {
		return merr.ErrorInternal("validate unique failed").WithCause(err)
	}
	if len(page.Items) > 0 && page.Items[0].UID() != uid {
		return merr.ErrorParams("strategy name %s already exists", name)
	}
	return nil
}

func (v *Validator) ValidateGroupExists(ctx context.Context, groupUID snowflake.ID) error {
	_, err := v.groupRepo.FindByID(ctx, groupUID)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("strategy group %s not found", groupUID)
		}
		return merr.ErrorInternal("get strategy group failed").WithCause(err)
	}
	return nil
}

