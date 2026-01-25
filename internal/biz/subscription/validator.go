package subscription

import (
	"context"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/pkg/merr"
)

type Validator struct {
	repo   Repository
	helper *klog.Helper
}

func NewValidator(repo Repository, helper *klog.Helper) *Validator {
	return &Validator{
		repo:   repo,
		helper: klog.NewHelper(klog.With(helper.Logger(), "biz", "subscription", "validator")),
	}
}

func (v *Validator) ValidateUnique(ctx context.Context, userID, namespaceUID snowflake.ID, name string) error {
	query := &ListQuery{
		PageRequest:  &shared.PageRequest{Page: 1, PageSize: 1},
		UserID:       userID,
		NamespaceUID: namespaceUID,
		Keyword:      name,
	}
	page, err := v.repo.List(ctx, query)
	if err != nil {
		return merr.ErrorInternal("validate unique failed").WithCause(err)
	}
	if len(page.Items) > 0 {
		return merr.ErrorParams("subscription name %s already exists", name)
	}
	return nil
}

func (v *Validator) ValidateUniqueForUpdate(ctx context.Context, userID, namespaceUID snowflake.ID, name string, uid snowflake.ID) error {
	query := &ListQuery{
		PageRequest:  &shared.PageRequest{Page: 1, PageSize: 1},
		UserID:       userID,
		NamespaceUID: namespaceUID,
		Keyword:      name,
	}
	page, err := v.repo.List(ctx, query)
	if err != nil {
		return merr.ErrorInternal("validate unique failed").WithCause(err)
	}
	if len(page.Items) > 0 && page.Items[0].UID() != uid {
		return merr.ErrorParams("subscription name %s already exists", name)
	}
	return nil
}

