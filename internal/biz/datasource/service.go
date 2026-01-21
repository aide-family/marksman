package datasource

import (
	"context"

	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/pkg/merr"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
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
		helper:    klog.NewHelper(klog.With(helper.Logger(), "biz", "datasource")),
	}
}

func (s *Service) Create(ctx context.Context, namespaceUID snowflake.ID, typ Type, engine Engine, name, endpoint, description string, config, metadata map[string]string) error {
	if err := s.validator.ValidateConfig(ctx, typ, config); err != nil {
		return err
	}

	ds := New(namespaceUID, typ, engine, name, endpoint, description, config, metadata)

	if err := s.validator.ValidateConnection(ctx, ds); err != nil {
		return merr.ErrorParams("datasource connection validation failed").WithCause(err)
	}

	if err := s.repo.Save(ctx, ds); err != nil {
		s.helper.Errorw("msg", "create datasource failed", "error", err, "name", name)
		return merr.ErrorInternal("create datasource %s failed", name).WithCause(err)
	}
	return nil
}

func (s *Service) Update(ctx context.Context, uid snowflake.ID, name, endpoint, description string, config, metadata map[string]string) error {
	ds, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("datasource %s not found", uid)
		}
		return merr.ErrorInternal("get datasource failed").WithCause(err)
	}

	if err := ds.UpdateName(name); err != nil {
		return err
	}
	if err := ds.UpdateEndpoint(endpoint); err != nil {
		return err
	}
	ds.UpdateDescription(description)
	ds.UpdateConfig(config)
	ds.UpdateMetadata(metadata)

	if err := s.repo.Save(ctx, ds); err != nil {
		s.helper.Errorw("msg", "update datasource failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update datasource %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, uid snowflake.ID, status Status) error {
	ds, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("datasource %s not found", uid)
		}
		return merr.ErrorInternal("get datasource failed").WithCause(err)
	}

	if status == StatusEnabled {
		if err := ds.Enable(); err != nil {
			return err
		}
	} else {
		if err := ds.Disable(); err != nil {
			return err
		}
	}

	if err := s.repo.Save(ctx, ds); err != nil {
		s.helper.Errorw("msg", "update datasource status failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update datasource status %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, uid snowflake.ID) error {
	if err := s.repo.Delete(ctx, uid); err != nil {
		s.helper.Errorw("msg", "delete datasource failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete datasource %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, uid snowflake.ID) (*DataSource, error) {
	ds, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("datasource %s not found", uid)
		}
		s.helper.Errorw("msg", "get datasource failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get datasource %s failed", uid).WithCause(err)
	}
	return ds, nil
}

func (s *Service) List(ctx context.Context, query *ListQuery) (*shared.Page[*DataSource], error) {
	page, err := s.repo.List(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "list datasource failed", "error", err, "query", query)
		return nil, merr.ErrorInternal("list datasource failed").WithCause(err)
	}
	return page, nil
}

func (s *Service) Select(ctx context.Context, query *SelectQuery) (*SelectResult, error) {
	result, err := s.repo.Select(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "select datasource failed", "error", err, "query", query)
		return nil, merr.ErrorInternal("select datasource failed").WithCause(err)
	}
	return result, nil
}

func (s *Service) TestConnection(ctx context.Context, uid snowflake.ID) error {
	ds, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("datasource %s not found", uid)
		}
		return merr.ErrorInternal("get datasource failed").WithCause(err)
	}
	return s.validator.ValidateConnection(ctx, ds)
}
