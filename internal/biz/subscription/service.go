package subscription

import (
	"context"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/internal/biz/vobj"
	"github.com/aide-family/sovereign/pkg/merr"
)

type Service struct {
	repo      Repository
	validator *Validator
	helper    *klog.Helper
}

func NewService(repo Repository, validator *Validator, helper *klog.Helper) *Service {
	return &Service{
		repo:      repo,
		validator: validator,
		helper:    klog.NewHelper(klog.With(helper.Logger(), "biz", "subscription")),
	}
}

func (s *Service) Create(ctx context.Context, userID, namespaceUID snowflake.ID, typ SubscriptionType, name string) error {
	if err := s.validator.ValidateUnique(ctx, userID, namespaceUID, name); err != nil {
		return err
	}

	subscription := New(userID, namespaceUID, typ, name)
	if err := s.repo.Save(ctx, subscription); err != nil {
		s.helper.Errorw("msg", "create subscription failed", "error", err, "name", name)
		return merr.ErrorInternal("create subscription %s failed", name).WithCause(err)
	}
	return nil
}

func (s *Service) Update(ctx context.Context, uid snowflake.ID, name string, description string) error {
	subscription, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("subscription %s not found", uid)
		}
		return merr.ErrorInternal("get subscription failed").WithCause(err)
	}

	if err := s.validator.ValidateUniqueForUpdate(ctx, subscription.UserID(), subscription.NamespaceUID(), name, uid); err != nil {
		return err
	}

	if err := subscription.UpdateName(name); err != nil {
		return err
	}
	subscription.UpdateDescription(description)

	if err := s.repo.Save(ctx, subscription); err != nil {
		s.helper.Errorw("msg", "update subscription failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update subscription %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, uid snowflake.ID, status vobj.GlobalStatus) error {
	subscription, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("subscription %s not found", uid)
		}
		return merr.ErrorInternal("get subscription failed").WithCause(err)
	}

	if status == vobj.GlobalStatusEnabled {
		if err := subscription.Enable(); err != nil {
			return err
		}
	} else {
		if err := subscription.Disable(); err != nil {
			return err
		}
	}

	if err := s.repo.Save(ctx, subscription); err != nil {
		s.helper.Errorw("msg", "update subscription status failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update subscription status %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, uid snowflake.ID) error {
	if err := s.repo.Delete(ctx, uid); err != nil {
		s.helper.Errorw("msg", "delete subscription failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete subscription %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, uid snowflake.ID) (*Subscription, error) {
	subscription, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("subscription %s not found", uid)
		}
		s.helper.Errorw("msg", "get subscription failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get subscription %s failed", uid).WithCause(err)
	}
	return subscription, nil
}

func (s *Service) List(ctx context.Context, query *ListQuery) (*shared.Page[*Subscription], error) {
	page, err := s.repo.List(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "list subscription failed", "error", err, "query", query)
		return nil, merr.ErrorInternal("list subscription failed").WithCause(err)
	}
	return page, nil
}

func (s *Service) Select(ctx context.Context, query *SelectQuery) (*SelectResult, error) {
	result, err := s.repo.Select(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "select subscription failed", "error", err, "query", query)
		return nil, merr.ErrorInternal("select subscription failed").WithCause(err)
	}
	return result, nil
}

func (s *Service) UpdateStrategyGroupUIDs(ctx context.Context, uid snowflake.ID, uids map[snowflake.ID]bool) error {
	subscription, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("subscription %s not found", uid)
		}
		return merr.ErrorInternal("get subscription failed").WithCause(err)
	}

	subscription.UpdateStrategyGroupUIDs(uids)
	if err := s.repo.Save(ctx, subscription); err != nil {
		s.helper.Errorw("msg", "update subscription strategy group uids failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update subscription strategy group uids %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateDataSourceUIDs(ctx context.Context, uid snowflake.ID, uids map[snowflake.ID]bool) error {
	subscription, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("subscription %s not found", uid)
		}
		return merr.ErrorInternal("get subscription failed").WithCause(err)
	}

	subscription.UpdateDataSourceUIDs(uids)
	if err := s.repo.Save(ctx, subscription); err != nil {
		s.helper.Errorw("msg", "update subscription datasource uids failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update subscription datasource uids %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateAlertLevels(ctx context.Context, uid snowflake.ID, levels []vobj.AlertLevel) error {
	subscription, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("subscription %s not found", uid)
		}
		return merr.ErrorInternal("get subscription failed").WithCause(err)
	}

	subscription.UpdateAlertLevels(levels)
	if err := s.repo.Save(ctx, subscription); err != nil {
		s.helper.Errorw("msg", "update subscription alert levels failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update subscription alert levels %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateNotifyTypes(ctx context.Context, uid snowflake.ID, types []vobj.NotifyType) error {
	subscription, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("subscription %s not found", uid)
		}
		return merr.ErrorInternal("get subscription failed").WithCause(err)
	}

	subscription.UpdateNotifyTypes(types)
	if err := s.repo.Save(ctx, subscription); err != nil {
		s.helper.Errorw("msg", "update subscription notify types failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update subscription notify types %s failed", uid).WithCause(err)
	}
	return nil
}

