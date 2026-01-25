package datasource

import (
	"context"

	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/internal/biz/vobj"
	"github.com/aide-family/sovereign/pkg/merr"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
)

type Service struct {
	repo                  Repository
	alertTemplateRepo     AlertTemplateRepository
	queryHistoryRepo      QueryHistoryRepository
	queryFavoriteRepo     QueryFavoriteRepository
	metadataRepo          DataSourceMetadataRepository
	proxyRepo             DataSourceProxyRepository
	validator             *Validator
	helper                *klog.Helper
}

func NewService(repo Repository, alertTemplateRepo AlertTemplateRepository, queryHistoryRepo QueryHistoryRepository, queryFavoriteRepo QueryFavoriteRepository, metadataRepo DataSourceMetadataRepository, proxyRepo DataSourceProxyRepository, validator *Validator, helper *klog.Helper) *Service {
	return &Service{
		repo:              repo,
		alertTemplateRepo: alertTemplateRepo,
		queryHistoryRepo:  queryHistoryRepo,
		queryFavoriteRepo: queryFavoriteRepo,
		metadataRepo:      metadataRepo,
		proxyRepo:         proxyRepo,
		validator:         validator,
		helper:            klog.NewHelper(klog.With(helper.Logger(), "biz", "datasource")),
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

// AlertTemplate 相关方法

func (s *Service) CreateAlertTemplate(ctx context.Context, datasourceUID snowflake.ID, name, titleTemplate, contentTemplate string) error {
	// 验证数据源存在
	_, err := s.repo.FindByID(ctx, datasourceUID)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("datasource %s not found", datasourceUID)
		}
		return merr.ErrorInternal("get datasource failed").WithCause(err)
	}

	template := NewAlertTemplate(datasourceUID, name, titleTemplate, contentTemplate)
	if err := s.alertTemplateRepo.Save(ctx, template); err != nil {
		s.helper.Errorw("msg", "create alert template failed", "error", err, "name", name)
		return merr.ErrorInternal("create alert template %s failed", name).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateAlertTemplate(ctx context.Context, uid snowflake.ID, name, titleTemplate, contentTemplate string) error {
	template, err := s.alertTemplateRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("alert template %s not found", uid)
		}
		s.helper.Errorw("msg", "get alert template failed", "error", err, "uid", uid)
		return merr.ErrorInternal("get alert template %s failed", uid).WithCause(err)
	}

	if name != "" {
		if err := template.UpdateName(name); err != nil {
			return err
		}
	}
	if titleTemplate != "" || contentTemplate != "" {
		template.UpdateTemplates(titleTemplate, contentTemplate)
	}

	if err := s.alertTemplateRepo.Save(ctx, template); err != nil {
		s.helper.Errorw("msg", "update alert template failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update alert template %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) DeleteAlertTemplate(ctx context.Context, uid snowflake.ID) error {
	if err := s.alertTemplateRepo.Delete(ctx, uid); err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("alert template %s not found", uid)
		}
		s.helper.Errorw("msg", "delete alert template failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete alert template %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) GetAlertTemplate(ctx context.Context, uid snowflake.ID) (*AlertTemplate, error) {
	template, err := s.alertTemplateRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("alert template %s not found", uid)
		}
		s.helper.Errorw("msg", "get alert template failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get alert template %s failed", uid).WithCause(err)
	}
	return template, nil
}

func (s *Service) ListAlertTemplate(ctx context.Context, query *AlertTemplateListQuery) (*shared.Page[*AlertTemplate], error) {
	page, err := s.alertTemplateRepo.List(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "list alert template failed", "error", err, "query", query)
		return nil, merr.ErrorInternal("list alert template failed").WithCause(err)
	}
	return page, nil
}

func (s *Service) UpdateAlertTemplateStatus(ctx context.Context, uid snowflake.ID, status vobj.GlobalStatus) error {
	template, err := s.alertTemplateRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("alert template %s not found", uid)
		}
		s.helper.Errorw("msg", "get alert template failed", "error", err, "uid", uid)
		return merr.ErrorInternal("get alert template %s failed", uid).WithCause(err)
	}

	if err := template.UpdateStatus(status); err != nil {
		return err
	}

	if err := s.alertTemplateRepo.Save(ctx, template); err != nil {
		s.helper.Errorw("msg", "update alert template status failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update alert template status %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) ApplyAlertTemplate(ctx context.Context, templateUID, strategyUID snowflake.ID) error {
	// 验证模板存在
	template, err := s.alertTemplateRepo.FindByID(ctx, templateUID)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("alert template %s not found", templateUID)
		}
		return merr.ErrorInternal("get alert template failed").WithCause(err)
	}

	if template.Status() != vobj.GlobalStatusEnabled {
		return merr.ErrorParams("alert template is disabled")
	}

	// TODO: 这里需要调用策略服务来应用模板
	// 暂时只验证模板存在和启用状态
	s.helper.Infow("msg", "apply alert template", "template_uid", templateUID, "strategy_uid", strategyUID)
	return nil
}

// Query 相关方法

func (s *Service) QueryDataSource(ctx context.Context, datasourceUID snowflake.ID, query, format string) (string, error) {
	// 验证数据源存在
	ds, err := s.repo.FindByID(ctx, datasourceUID)
	if err != nil {
		if merr.IsNotFound(err) {
			return "", merr.ErrorNotFound("datasource %s not found", datasourceUID)
		}
		return "", merr.ErrorInternal("get datasource failed").WithCause(err)
	}

	if !ds.IsEnabled() {
		return "", merr.ErrorParams("datasource is disabled")
	}

	// TODO: 实际执行查询逻辑，根据数据源类型和引擎执行查询
	// 这里返回空结果，实际应该调用对应的查询引擎
	s.helper.Infow("msg", "query datasource", "datasource_uid", datasourceUID, "query", query, "format", format)

	// 保存查询历史
	history := NewQueryHistory(datasourceUID, query, format)
	if err := s.queryHistoryRepo.Save(ctx, history); err != nil {
		s.helper.Errorw("msg", "save query history failed", "error", err)
		// 不返回错误，查询历史保存失败不影响查询结果
	}

	return "{}", nil
}

func (s *Service) ListQueryHistory(ctx context.Context, datasourceUID snowflake.ID, page, pageSize int32) ([]*QueryHistory, int64, error) {
	histories, total, err := s.queryHistoryRepo.List(ctx, datasourceUID, page, pageSize)
	if err != nil {
		s.helper.Errorw("msg", "list query history failed", "error", err, "datasource_uid", datasourceUID)
		return nil, 0, merr.ErrorInternal("list query history failed").WithCause(err)
	}
	return histories, total, nil
}

func (s *Service) SaveQueryFavorite(ctx context.Context, datasourceUID snowflake.ID, name, query string) (uint32, error) {
	// 验证数据源存在
	_, err := s.repo.FindByID(ctx, datasourceUID)
	if err != nil {
		if merr.IsNotFound(err) {
			return 0, merr.ErrorNotFound("datasource %s not found", datasourceUID)
		}
		return 0, merr.ErrorInternal("get datasource failed").WithCause(err)
	}

	favorite := NewQueryFavorite(datasourceUID, name, query)
	if err := s.queryFavoriteRepo.Save(ctx, favorite); err != nil {
		s.helper.Errorw("msg", "save query favorite failed", "error", err, "name", name)
		return 0, merr.ErrorInternal("save query favorite %s failed", name).WithCause(err)
	}
	return favorite.ID(), nil
}

func (s *Service) ListQueryFavorites(ctx context.Context, datasourceUID snowflake.ID) ([]*QueryFavorite, error) {
	favorites, err := s.queryFavoriteRepo.List(ctx, datasourceUID)
	if err != nil {
		s.helper.Errorw("msg", "list query favorites failed", "error", err, "datasource_uid", datasourceUID)
		return nil, merr.ErrorInternal("list query favorites failed").WithCause(err)
	}
	return favorites, nil
}

func (s *Service) DeleteQueryFavorite(ctx context.Context, id uint32) error {
	if err := s.queryFavoriteRepo.Delete(ctx, id); err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("query favorite %d not found", id)
		}
		s.helper.Errorw("msg", "delete query favorite failed", "error", err, "id", id)
		return merr.ErrorInternal("delete query favorite %d failed", id).WithCause(err)
	}
	return nil
}

// Metadata 相关方法

func (s *Service) ListDataSourceMetadata(ctx context.Context, datasourceUID snowflake.ID, metadataType string) ([]*DataSourceMetadata, error) {
	// 验证数据源存在
	_, err := s.repo.FindByID(ctx, datasourceUID)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("datasource %s not found", datasourceUID)
		}
		return nil, merr.ErrorInternal("get datasource failed").WithCause(err)
	}

	metadataList, err := s.metadataRepo.List(ctx, datasourceUID, metadataType)
	if err != nil {
		s.helper.Errorw("msg", "list datasource metadata failed", "error", err, "datasource_uid", datasourceUID)
		return nil, merr.ErrorInternal("list datasource metadata failed").WithCause(err)
	}
	return metadataList, nil
}

func (s *Service) GetDataSourceMetadata(ctx context.Context, datasourceUID snowflake.ID, key string) (*DataSourceMetadata, error) {
	// 验证数据源存在
	_, err := s.repo.FindByID(ctx, datasourceUID)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("datasource %s not found", datasourceUID)
		}
		return nil, merr.ErrorInternal("get datasource failed").WithCause(err)
	}

	metadata, err := s.metadataRepo.FindByKey(ctx, datasourceUID, key)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("metadata %s not found for datasource %s", key, datasourceUID)
		}
		s.helper.Errorw("msg", "get datasource metadata failed", "error", err, "datasource_uid", datasourceUID, "key", key)
		return nil, merr.ErrorInternal("get datasource metadata failed").WithCause(err)
	}
	return metadata, nil
}

func (s *Service) RefreshDataSourceMetadata(ctx context.Context, datasourceUID snowflake.ID) error {
	// 验证数据源存在
	ds, err := s.repo.FindByID(ctx, datasourceUID)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("datasource %s not found", datasourceUID)
		}
		return merr.ErrorInternal("get datasource failed").WithCause(err)
	}

	if !ds.IsEnabled() {
		return merr.ErrorParams("datasource is disabled")
	}

	// TODO: 实际刷新元数据逻辑，根据数据源类型和引擎获取元数据
	// 这里先删除旧的元数据，然后获取新的元数据
	if err := s.metadataRepo.DeleteByDataSourceUID(ctx, datasourceUID); err != nil {
		s.helper.Errorw("msg", "delete old metadata failed", "error", err, "datasource_uid", datasourceUID)
		// 不返回错误，继续刷新
	}

	// TODO: 根据数据源类型和引擎获取元数据并保存
	// 这里暂时只记录日志
	s.helper.Infow("msg", "refresh datasource metadata", "datasource_uid", datasourceUID, "type", ds.Type(), "engine", ds.Engine())
	return nil
}

// Proxy 相关方法

func (s *Service) CreateDataSourceProxy(ctx context.Context, namespaceUID, datasourceUID snowflake.ID, typ, name string, config map[string]string) error {
	// 验证数据源存在
	_, err := s.repo.FindByID(ctx, datasourceUID)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("datasource %s not found", datasourceUID)
		}
		return merr.ErrorInternal("get datasource failed").WithCause(err)
	}

	proxy := NewDataSourceProxy(namespaceUID, datasourceUID, typ, name, config)
	if err := s.proxyRepo.Save(ctx, proxy); err != nil {
		s.helper.Errorw("msg", "create datasource proxy failed", "error", err, "name", name)
		return merr.ErrorInternal("create datasource proxy %s failed", name).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateDataSourceProxy(ctx context.Context, uid snowflake.ID, name string, config map[string]string) error {
	proxy, err := s.proxyRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("datasource proxy %s not found", uid)
		}
		s.helper.Errorw("msg", "get datasource proxy failed", "error", err, "uid", uid)
		return merr.ErrorInternal("get datasource proxy %s failed", uid).WithCause(err)
	}

	if name != "" {
		if err := proxy.UpdateName(name); err != nil {
			return err
		}
	}
	if config != nil {
		proxy.UpdateConfig(config)
	}

	if err := s.proxyRepo.Save(ctx, proxy); err != nil {
		s.helper.Errorw("msg", "update datasource proxy failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update datasource proxy %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) DeleteDataSourceProxy(ctx context.Context, uid snowflake.ID) error {
	if err := s.proxyRepo.Delete(ctx, uid); err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("datasource proxy %s not found", uid)
		}
		s.helper.Errorw("msg", "delete datasource proxy failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete datasource proxy %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) GetDataSourceProxy(ctx context.Context, uid snowflake.ID) (*DataSourceProxy, error) {
	proxy, err := s.proxyRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("datasource proxy %s not found", uid)
		}
		s.helper.Errorw("msg", "get datasource proxy failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get datasource proxy %s failed", uid).WithCause(err)
	}
	return proxy, nil
}

func (s *Service) ListDataSourceProxy(ctx context.Context, namespaceUID snowflake.ID, page, pageSize int32) ([]*DataSourceProxy, int64, error) {
	proxies, total, err := s.proxyRepo.List(ctx, namespaceUID, page, pageSize)
	if err != nil {
		s.helper.Errorw("msg", "list datasource proxy failed", "error", err, "namespace_uid", namespaceUID)
		return nil, 0, merr.ErrorInternal("list datasource proxy failed").WithCause(err)
	}
	return proxies, total, nil
}
