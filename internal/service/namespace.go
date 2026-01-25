package service

import (
	"context"
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/aide-family/magicbox/strutil"
	"github.com/aide-family/magicbox/strutil/cnst"

	"github.com/aide-family/sovereign/internal/biz/namespace"
	"github.com/aide-family/sovereign/internal/biz/shared"
	apiv1 "github.com/aide-family/sovereign/pkg/api/v1"
	"github.com/aide-family/sovereign/pkg/enum"
	"github.com/aide-family/sovereign/pkg/merr"
	"github.com/aide-family/sovereign/pkg/middler"
)

func NewNamespaceService(namespaceService *namespace.Service) *NamespaceService {
	return &NamespaceService{
		namespaceService: namespaceService,
	}
}

type NamespaceService struct {
	apiv1.UnimplementedNamespaceServer

	namespaceService *namespace.Service
}

func (s *NamespaceService) CreateNamespace(ctx context.Context, req *apiv1.CreateNamespaceRequest) (*apiv1.CreateNamespaceReply, error) {
	if err := s.namespaceService.Create(ctx, req.Name, req.Metadata); err != nil {
		return nil, err
	}
	return &apiv1.CreateNamespaceReply{}, nil
}

func (s *NamespaceService) UpdateNamespace(ctx context.Context, req *apiv1.UpdateNamespaceRequest) (*apiv1.UpdateNamespaceReply, error) {
	if err := s.namespaceService.Update(ctx, snowflake.ParseInt64(req.Uid), req.Name, req.Metadata); err != nil {
		return nil, err
	}
	return &apiv1.UpdateNamespaceReply{}, nil
}

func (s *NamespaceService) UpdateNamespaceStatus(ctx context.Context, req *apiv1.UpdateNamespaceStatusRequest) (*apiv1.UpdateNamespaceStatusReply, error) {
	status := namespace.Status(req.Status)
	if err := s.namespaceService.UpdateStatus(ctx, snowflake.ParseInt64(req.Uid), status); err != nil {
		return nil, err
	}
	return &apiv1.UpdateNamespaceStatusReply{}, nil
}

func (s *NamespaceService) DeleteNamespace(ctx context.Context, req *apiv1.DeleteNamespaceRequest) (*apiv1.DeleteNamespaceReply, error) {
	if err := s.namespaceService.Delete(ctx, snowflake.ParseInt64(req.Uid)); err != nil {
		return nil, err
	}
	return &apiv1.DeleteNamespaceReply{}, nil
}

func (s *NamespaceService) GetNamespace(ctx context.Context, req *apiv1.GetNamespaceRequest) (*apiv1.NamespaceItem, error) {
	ns, err := s.namespaceService.Get(ctx, snowflake.ParseInt64(req.Uid))
	if err != nil {
		return nil, err
	}
	return toAPIV1NamespaceItem(ns), nil
}

func (s *NamespaceService) ListNamespace(ctx context.Context, req *apiv1.ListNamespaceRequest) (*apiv1.ListNamespaceReply, error) {
	query := &namespace.ListQuery{
		PageRequest: shared.NewPageRequest(req.Page, req.PageSize),
		Keyword:     req.Keyword,
		Status:      namespace.Status(req.Status),
	}
	page, err := s.namespaceService.List(ctx, query)
	if err != nil {
		return nil, err
	}
	return toAPIV1ListNamespaceReply(page), nil
}

func (s *NamespaceService) SelectNamespace(ctx context.Context, req *apiv1.SelectNamespaceRequest) (*apiv1.SelectNamespaceReply, error) {
	var nextUID snowflake.ID
	if req.NextUID > 0 {
		nextUID = snowflake.ParseInt64(req.NextUID)
	}
	query := &namespace.SelectQuery{
		Keyword: req.Keyword,
		Limit:   req.Limit,
		NextUID: nextUID,
		Status:  namespace.Status(req.Status),
	}
	result, err := s.namespaceService.Select(ctx, query)
	if err != nil {
		return nil, err
	}
	return toAPIV1SelectNamespaceReply(result), nil
}

func (s *NamespaceService) HasNamespace(ctx context.Context) error {
	ns := middler.GetNamespace(ctx)
	if strutil.IsEmpty(ns) {
		return merr.ErrorForbidden("namespace is required, please set the namespace in the request header or metadata, Example: %s: default", cnst.HTTPHeaderXNamespace)
	}
	namespaceEntity, err := s.namespaceService.GetByName(ctx, ns)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorForbidden("namespace %s not found", ns)
		}
		return err
	}
	if !namespaceEntity.IsEnabled() {
		return merr.ErrorForbidden("namespace %s is not enabled", ns)
	}
	return nil
}

// toAPIV1NamespaceItem converts namespace entity to API response
func toAPIV1NamespaceItem(ns *namespace.Namespace) *apiv1.NamespaceItem {
	return &apiv1.NamespaceItem{
		Uid:       ns.UID().Int64(),
		Name:      ns.Name(),
		Metadata:  ns.Metadata(),
		Status:    enum.GlobalStatus(ns.Status()),
		CreatedAt: ns.CreatedAt().Format(time.DateTime),
		UpdatedAt: ns.UpdatedAt().Format(time.DateTime),
	}
}

// toAPIV1ListNamespaceReply converts namespace page to API response
func toAPIV1ListNamespaceReply(page *shared.Page[*namespace.Namespace]) *apiv1.ListNamespaceReply {
	items := make([]*apiv1.NamespaceItem, 0, len(page.Items))
	for _, ns := range page.Items {
		items = append(items, toAPIV1NamespaceItem(ns))
	}
	return &apiv1.ListNamespaceReply{
		Items:    items,
		Total:    page.Total,
		Page:     page.Page,
		PageSize: page.PageSize,
	}
}

// toAPIV1SelectNamespaceReply converts namespace select result to API response
func toAPIV1SelectNamespaceReply(result *namespace.SelectResult) *apiv1.SelectNamespaceReply {
	selectItems := make([]*apiv1.NamespaceItemSelect, 0, len(result.Items))
	for _, item := range result.Items {
		selectItems = append(selectItems, &apiv1.NamespaceItemSelect{
			Value:    item.UID.Int64(),
			Label:    item.Name,
			Disabled: item.Disabled,
			Tooltip:  item.Tooltip,
		})
	}
	return &apiv1.SelectNamespaceReply{
		Items:   selectItems,
		Total:   result.Total,
		NextUID: result.NextUID.Int64(),
		HasMore: result.HasMore,
	}
}
