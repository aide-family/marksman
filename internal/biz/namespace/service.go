package namespace

import (
	"context"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/aide-family/sovereign/internal/biz/shared"
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
		helper:    klog.NewHelper(klog.With(helper.Logger(), "biz", "namespace")),
	}
}

func (s *Service) Create(ctx context.Context, name string, metadata map[string]string) error {
	if err := s.validator.ValidateUnique(ctx, name); err != nil {
		return err
	}
	
	ns := New(name, metadata)
	if err := s.repo.Save(ctx, ns); err != nil {
		s.helper.Errorw("msg", "create namespace failed", "error", err, "name", name)
		return merr.ErrorInternal("create namespace %s failed", name).WithCause(err)
	}
	return nil
}

func (s *Service) Update(ctx context.Context, uid snowflake.ID, name string, metadata map[string]string) error {
	if err := s.validator.ValidateUniqueForUpdate(ctx, name, uid); err != nil {
		return err
	}
	
	ns, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("namespace %s not found", uid)
		}
		return merr.ErrorInternal("get namespace failed").WithCause(err)
	}
	
	if err := ns.UpdateName(name); err != nil {
		return err
	}
	ns.UpdateMetadata(metadata)
	
	if err := s.repo.Save(ctx, ns); err != nil {
		s.helper.Errorw("msg", "update namespace failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update namespace %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, uid snowflake.ID, status Status) error {
	ns, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("namespace %s not found", uid)
		}
		return merr.ErrorInternal("get namespace failed").WithCause(err)
	}
	
	if status == StatusEnabled {
		if err := ns.Enable(); err != nil {
			return err
		}
	} else {
		if err := ns.Disable(); err != nil {
			return err
		}
	}
	
	if err := s.repo.Save(ctx, ns); err != nil {
		s.helper.Errorw("msg", "update namespace status failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update namespace status %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, uid snowflake.ID) error {
	if err := s.repo.Delete(ctx, uid); err != nil {
		s.helper.Errorw("msg", "delete namespace failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete namespace %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, uid snowflake.ID) (*Namespace, error) {
	ns, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("namespace %s not found", uid)
		}
		s.helper.Errorw("msg", "get namespace failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get namespace %s failed", uid).WithCause(err)
	}
	return ns, nil
}

func (s *Service) GetByName(ctx context.Context, name string) (*Namespace, error) {
	ns, err := s.repo.FindByName(ctx, name)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("namespace %s not found", name)
		}
		s.helper.Errorw("msg", "get namespace failed", "error", err, "name", name)
		return nil, merr.ErrorInternal("get namespace %s failed", name).WithCause(err)
	}
	return ns, nil
}

func (s *Service) List(ctx context.Context, query *ListQuery) (*shared.Page[*Namespace], error) {
	page, err := s.repo.List(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "list namespace failed", "error", err, "query", query)
		return nil, merr.ErrorInternal("list namespace failed").WithCause(err)
	}
	return page, nil
}

func (s *Service) Select(ctx context.Context, query *SelectQuery) (*SelectResult, error) {
	result, err := s.repo.Select(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "select namespace failed", "error", err, "query", query)
		return nil, merr.ErrorInternal("select namespace failed").WithCause(err)
	}
	return result, nil
}

