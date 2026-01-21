package impl

import (
	_ "github.com/aide-family/sovereign/pkg/repo/namespace/v1/gormimpl"

	"context"
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/sovereign/internal/biz/namespace"
	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/internal/conf"
	"github.com/aide-family/sovereign/internal/data"
	"github.com/aide-family/sovereign/pkg/enum"
	"github.com/aide-family/sovereign/pkg/merr"
	"github.com/aide-family/sovereign/pkg/repo"
	namespacev1 "github.com/aide-family/sovereign/pkg/repo/namespace/v1"
)

func NewNamespaceRepository(c *conf.Bootstrap, d *data.Data) (namespace.Repository, error) {
	repoConfig := c.GetNamespaceConfig()
	version := repoConfig.GetVersion()
	driver := repoConfig.GetDriver()
	switch version {
	default:
		factory, ok := repo.GetNamespaceV1Factory(driver)
		if !ok {
			return nil, merr.ErrorInternalServer("namespace repository factory not found")
		}
		repoImpl, close, err := factory(repoConfig)
		if err != nil {
			return nil, err
		}
		d.AppendClose("namespaceRepo", close)
		return &namespaceRepository{repo: repoImpl}, nil
	}
}

type namespaceRepository struct {
	repo namespacev1.Repository
}

// Save implements namespace.Repository.
func (n *namespaceRepository) Save(ctx context.Context, ns *namespace.Namespace) error {
	uid := ns.UID().Int64()
	
	// 如果 UID 为 0，说明是新实体，直接创建
	if uid == 0 {
		_, err := n.repo.CreateNamespace(ctx, &namespacev1.CreateNamespaceRequest{
			Name:     ns.Name(),
			Metadata: ns.Metadata(),
			Status:   enum.GlobalStatus(ns.Status()),
		})
		return err
	}

	// 检查是否存在
	existing, err := n.repo.GetNamespace(ctx, &namespacev1.GetNamespaceRequest{
		Uid: uid,
	})
	if err != nil {
		if merr.IsNotFound(err) {
			// 不存在，创建
			_, err := n.repo.CreateNamespace(ctx, &namespacev1.CreateNamespaceRequest{
				Name:     ns.Name(),
				Metadata: ns.Metadata(),
				Status:   enum.GlobalStatus(ns.Status()),
			})
			return err
		}
		return err
	}

	// 存在，更新
	if existing != nil {
		_, err := n.repo.UpdateNamespace(ctx, &namespacev1.UpdateNamespaceRequest{
			Uid:      uid,
			Name:     ns.Name(),
			Metadata: ns.Metadata(),
		})
		return err
	}

	return nil
}

// FindByID implements namespace.Repository.
func (n *namespaceRepository) FindByID(ctx context.Context, uid snowflake.ID) (*namespace.Namespace, error) {
	model, err := n.repo.GetNamespace(ctx, &namespacev1.GetNamespaceRequest{
		Uid: uid.Int64(),
	})
	if err != nil {
		return nil, err
	}
	return toNamespaceEntity(model), nil
}

// FindByName implements namespace.Repository.
func (n *namespaceRepository) FindByName(ctx context.Context, name string) (*namespace.Namespace, error) {
	model, err := n.repo.GetNamespaceByName(ctx, &namespacev1.GetNamespaceByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return toNamespaceEntity(model), nil
}

// Delete implements namespace.Repository.
func (n *namespaceRepository) Delete(ctx context.Context, uid snowflake.ID) error {
	_, err := n.repo.DeleteNamespace(ctx, &namespacev1.DeleteNamespaceRequest{
		Uid: uid.Int64(),
	})
	return err
}

// List implements namespace.Repository.
func (n *namespaceRepository) List(ctx context.Context, query *namespace.ListQuery) (*shared.Page[*namespace.Namespace], error) {
	req := &namespacev1.ListNamespaceRequest{
		Page:     query.Page,
		PageSize: query.PageSize,
		Keyword:  query.Keyword,
		Status:   enum.GlobalStatus(query.Status),
	}

	resp, err := n.repo.ListNamespace(ctx, req)
	if err != nil {
		return nil, err
	}

	items := make([]*namespace.Namespace, 0, len(resp.Namespaces))
	for _, model := range resp.Namespaces {
		items = append(items, toNamespaceEntity(model))
	}

	query.WithTotal(resp.Total)
	return shared.NewPage(query.PageRequest, items), nil
}

// Select implements namespace.Repository.
func (n *namespaceRepository) Select(ctx context.Context, query *namespace.SelectQuery) (*namespace.SelectResult, error) {
	req := &namespacev1.SelectNamespaceRequest{
		Keyword: query.Keyword,
		Limit:   query.Limit,
		NextUID: query.NextUID.Int64(),
		Status:  enum.GlobalStatus(query.Status),
	}

	resp, err := n.repo.SelectNamespace(ctx, req)
	if err != nil {
		return nil, err
	}

	items := make([]*namespace.SelectItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &namespace.SelectItem{
			UID:      snowflake.ParseInt64(item.Value),
			Name:     item.Label,
			Disabled: item.Disabled,
			Tooltip:  item.Tooltip,
		})
	}

	return &namespace.SelectResult{
		Items:   items,
		Total:   resp.Total,
		NextUID: snowflake.ParseInt64(resp.NextUID),
		HasMore: resp.HasMore,
	}, nil
}

// toNamespaceEntity converts repo model to domain entity
func toNamespaceEntity(model *namespacev1.NamespaceModel) *namespace.Namespace {
	return namespace.FromModel(
		snowflake.ParseInt64(model.Uid),
		model.Name,
		model.Metadata,
		namespace.Status(model.Status),
		time.Unix(model.CreatedAt, 0),
		time.Unix(model.UpdatedAt, 0),
	)
}
