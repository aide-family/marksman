package strategy

import (
	"context"

	"github.com/aide-family/sovereign/internal/biz/shared"
	"github.com/aide-family/sovereign/internal/biz/vobj"
	"github.com/aide-family/sovereign/pkg/merr"
	"github.com/bwmarrin/snowflake"
	klog "github.com/go-kratos/kratos/v2/log"
)

type Service struct {
	repo         Repository
	groupRepo    GroupRepository
	receiverRepo ReceiverRepository
	validator    *Validator
	helper       *klog.Helper
}

func NewService(repo Repository, groupRepo GroupRepository, receiverRepo ReceiverRepository, validator *Validator, helper *klog.Helper) *Service {
	return &Service{
		repo:         repo,
		groupRepo:    groupRepo,
		receiverRepo: receiverRepo,
		validator:    validator,
		helper:       klog.NewHelper(klog.With(helper.Logger(), "biz", "strategy")),
	}
}

func (s *Service) Create(ctx context.Context, namespaceUID snowflake.ID, typ vobj.StrategyType, name string) error {
	if err := s.validator.ValidateUnique(ctx, namespaceUID, name); err != nil {
		return err
	}

	strategy := New(namespaceUID, typ, name)
	if err := s.repo.Save(ctx, strategy); err != nil {
		s.helper.Errorw("msg", "create strategy failed", "error", err, "name", name)
		return merr.ErrorInternal("create strategy %s failed", name).WithCause(err)
	}
	return nil
}

func (s *Service) Update(ctx context.Context, uid snowflake.ID, name string, description string) error {
	strategy, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("strategy %s not found", uid)
		}
		return merr.ErrorInternal("get strategy failed").WithCause(err)
	}

	if err := s.validator.ValidateUniqueForUpdate(ctx, strategy.NamespaceUID(), name, uid); err != nil {
		return err
	}

	if err := strategy.UpdateName(name); err != nil {
		return err
	}
	strategy.UpdateDescription(description)

	if err := s.repo.Save(ctx, strategy); err != nil {
		s.helper.Errorw("msg", "update strategy failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update strategy %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, uid snowflake.ID, status vobj.GlobalStatus) error {
	strategy, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("strategy %s not found", uid)
		}
		return merr.ErrorInternal("get strategy failed").WithCause(err)
	}

	if status == vobj.GlobalStatusEnabled {
		if err := strategy.Enable(); err != nil {
			return err
		}
	} else {
		if err := strategy.Disable(); err != nil {
			return err
		}
	}

	if err := s.repo.Save(ctx, strategy); err != nil {
		s.helper.Errorw("msg", "update strategy status failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update strategy status %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, uid snowflake.ID) error {
	if err := s.repo.Delete(ctx, uid); err != nil {
		s.helper.Errorw("msg", "delete strategy failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete strategy %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, uid snowflake.ID) (*Strategy, error) {
	strategy, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("strategy %s not found", uid)
		}
		s.helper.Errorw("msg", "get strategy failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get strategy %s failed", uid).WithCause(err)
	}
	return strategy, nil
}

func (s *Service) List(ctx context.Context, query *ListQuery) (*shared.Page[*Strategy], error) {
	page, err := s.repo.List(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "list strategy failed", "error", err, "query", query)
		return nil, merr.ErrorInternal("list strategy failed").WithCause(err)
	}
	return page, nil
}

func (s *Service) Select(ctx context.Context, query *SelectQuery) (*SelectResult, error) {
	result, err := s.repo.Select(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "select strategy failed", "error", err, "query", query)
		return nil, merr.ErrorInternal("select strategy failed").WithCause(err)
	}
	return result, nil
}

// Group methods
func (s *Service) CreateGroup(ctx context.Context, namespaceUID snowflake.ID, name string) error {
	group := NewGroup(namespaceUID, name)
	if err := s.groupRepo.Save(ctx, group); err != nil {
		s.helper.Errorw("msg", "create strategy group failed", "error", err, "name", name)
		return merr.ErrorInternal("create strategy group %s failed", name).WithCause(err)
	}
	return nil
}

func (s *Service) GetGroup(ctx context.Context, uid snowflake.ID) (*StrategyGroup, error) {
	group, err := s.groupRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("strategy group %s not found", uid)
		}
		s.helper.Errorw("msg", "get strategy group failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get strategy group %s failed", uid).WithCause(err)
	}
	return group, nil
}

func (s *Service) ListGroups(ctx context.Context, query *GroupListQuery) (*shared.Page[*StrategyGroup], error) {
	page, err := s.groupRepo.List(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "list strategy group failed", "error", err, "query", query)
		return nil, merr.ErrorInternal("list strategy group failed").WithCause(err)
	}
	return page, nil
}

func (s *Service) UpdateGroup(ctx context.Context, uid snowflake.ID, name string, description string) error {
	group, err := s.groupRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("strategy group %s not found", uid)
		}
		return merr.ErrorInternal("get strategy group failed").WithCause(err)
	}

	if err := group.UpdateName(name); err != nil {
		return err
	}
	group.UpdateDescription(description)

	if err := s.groupRepo.Save(ctx, group); err != nil {
		s.helper.Errorw("msg", "update strategy group failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update strategy group %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateGroupStatus(ctx context.Context, uid snowflake.ID, status vobj.GlobalStatus) error {
	group, err := s.groupRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("strategy group %s not found", uid)
		}
		return merr.ErrorInternal("get strategy group failed").WithCause(err)
	}

	if status == vobj.GlobalStatusEnabled {
		if err := group.Enable(); err != nil {
			return err
		}
	} else {
		if err := group.Disable(); err != nil {
			return err
		}
	}

	if err := s.groupRepo.Save(ctx, group); err != nil {
		s.helper.Errorw("msg", "update strategy group status failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update strategy group status %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) DeleteGroup(ctx context.Context, uid snowflake.ID) error {
	// 检查是否有策略关联到此策略组
	query := &ListQuery{
		PageRequest: &shared.PageRequest{Page: 1, PageSize: 1},
		GroupUID:    uid,
	}
	page, err := s.repo.List(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "check strategies in group failed", "error", err, "uid", uid)
		return merr.ErrorInternal("check strategies in group failed").WithCause(err)
	}
	if len(page.Items) > 0 {
		return merr.ErrorParams("cannot delete strategy group %s, it contains %d strategies", uid, page.Total)
	}

	if err := s.groupRepo.Delete(ctx, uid); err != nil {
		s.helper.Errorw("msg", "delete strategy group failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete strategy group %s failed", uid).WithCause(err)
	}
	return nil
}

// GetGroupStrategies 获取策略组下的所有策略
func (s *Service) GetGroupStrategies(ctx context.Context, groupUID snowflake.ID, query *ListQuery) (*shared.Page[*Strategy], error) {
	// 验证策略组存在
	if _, err := s.groupRepo.FindByID(ctx, groupUID); err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("strategy group %s not found", groupUID)
		}
		return nil, merr.ErrorInternal("get strategy group failed").WithCause(err)
	}

	// 查询该策略组下的策略
	if query == nil {
		query = &ListQuery{
			PageRequest: shared.NewPageRequest(1, 20),
		}
	}
	query.GroupUID = groupUID

	page, err := s.repo.List(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "list strategies in group failed", "error", err, "groupUID", groupUID)
		return nil, merr.ErrorInternal("list strategies in group failed").WithCause(err)
	}
	return page, nil
}

// AddStrategyToGroup 将策略添加到策略组
func (s *Service) AddStrategyToGroup(ctx context.Context, strategyUID, groupUID snowflake.ID) error {
	// 验证策略组存在
	if _, err := s.groupRepo.FindByID(ctx, groupUID); err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("strategy group %s not found", groupUID)
		}
		return merr.ErrorInternal("get strategy group failed").WithCause(err)
	}

	// 获取策略
	strategy, err := s.repo.FindByID(ctx, strategyUID)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("strategy %s not found", strategyUID)
		}
		return merr.ErrorInternal("get strategy failed").WithCause(err)
	}

	// 更新策略的 groupUID
	strategy.UpdateGroupUID(groupUID)

	if err := s.repo.Save(ctx, strategy); err != nil {
		s.helper.Errorw("msg", "add strategy to group failed", "error", err, "strategyUID", strategyUID, "groupUID", groupUID)
		return merr.ErrorInternal("add strategy to group failed").WithCause(err)
	}
	return nil
}

// RemoveStrategyFromGroup 将策略从策略组移除
func (s *Service) RemoveStrategyFromGroup(ctx context.Context, strategyUID snowflake.ID) error {
	// 获取策略
	strategy, err := s.repo.FindByID(ctx, strategyUID)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("strategy %s not found", strategyUID)
		}
		return merr.ErrorInternal("get strategy failed").WithCause(err)
	}

	// 将 groupUID 设置为 0（表示不属于任何策略组）
	strategy.UpdateGroupUID(0)

	if err := s.repo.Save(ctx, strategy); err != nil {
		s.helper.Errorw("msg", "remove strategy from group failed", "error", err, "strategyUID", strategyUID)
		return merr.ErrorInternal("remove strategy from group failed").WithCause(err)
	}
	return nil
}

// Receiver methods
func (s *Service) CreateReceiver(ctx context.Context, namespaceUID snowflake.ID, typ vobj.ReceiverType, name string, description string, userIDs map[snowflake.ID]bool, labelMatch map[string]string, notifyTypes []vobj.NotifyType) error {
	receiver := NewReceiver(namespaceUID, typ, name)
	if description != "" {
		receiver.UpdateDescription(description)
	}
	if userIDs != nil {
		receiver.UpdateUserIDs(userIDs)
	}
	if labelMatch != nil {
		receiver.UpdateLabelMatch(labelMatch)
	}
	if notifyTypes != nil {
		receiver.UpdateNotifyTypes(notifyTypes)
	}
	if err := s.receiverRepo.Save(ctx, receiver); err != nil {
		s.helper.Errorw("msg", "create receiver failed", "error", err, "name", name)
		return merr.ErrorInternal("create receiver %s failed", name).WithCause(err)
	}
	return nil
}

func (s *Service) GetReceiver(ctx context.Context, uid snowflake.ID) (*Receiver, error) {
	receiver, err := s.receiverRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return nil, merr.ErrorNotFound("receiver %s not found", uid)
		}
		s.helper.Errorw("msg", "get receiver failed", "error", err, "uid", uid)
		return nil, merr.ErrorInternal("get receiver %s failed", uid).WithCause(err)
	}
	return receiver, nil
}

func (s *Service) ListReceivers(ctx context.Context, query *ReceiverListQuery) (*shared.Page[*Receiver], error) {
	page, err := s.receiverRepo.List(ctx, query)
	if err != nil {
		s.helper.Errorw("msg", "list receiver failed", "error", err, "query", query)
		return nil, merr.ErrorInternal("list receiver failed").WithCause(err)
	}
	return page, nil
}

func (s *Service) UpdateReceiver(ctx context.Context, uid snowflake.ID, name string, description string, userIDs map[snowflake.ID]bool, labelMatch map[string]string, notifyTypes []vobj.NotifyType) error {
	receiver, err := s.receiverRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("receiver %s not found", uid)
		}
		return merr.ErrorInternal("get receiver failed").WithCause(err)
	}

	if name != "" {
		if err := receiver.UpdateName(name); err != nil {
			return err
		}
	}
	if description != "" {
		receiver.UpdateDescription(description)
	}
	if userIDs != nil {
		receiver.UpdateUserIDs(userIDs)
	}
	if labelMatch != nil {
		receiver.UpdateLabelMatch(labelMatch)
	}
	if notifyTypes != nil {
		receiver.UpdateNotifyTypes(notifyTypes)
	}

	if err := s.receiverRepo.Save(ctx, receiver); err != nil {
		s.helper.Errorw("msg", "update receiver failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update receiver %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) DeleteReceiver(ctx context.Context, uid snowflake.ID) error {
	if err := s.receiverRepo.Delete(ctx, uid); err != nil {
		s.helper.Errorw("msg", "delete receiver failed", "error", err, "uid", uid)
		return merr.ErrorInternal("delete receiver %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateGroupUpgradeMode(ctx context.Context, uid snowflake.ID, mode vobj.UpgradeMode) error {
	group, err := s.groupRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("strategy group %s not found", uid)
		}
		return merr.ErrorInternal("get strategy group failed").WithCause(err)
	}

	group.UpdateUpgradeMode(mode)
	if err := s.groupRepo.Save(ctx, group); err != nil {
		s.helper.Errorw("msg", "update strategy group upgrade mode failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update strategy group upgrade mode %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateGroupUpgradeConfig(ctx context.Context, uid snowflake.ID, config string) error {
	group, err := s.groupRepo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("strategy group %s not found", uid)
		}
		return merr.ErrorInternal("get strategy group failed").WithCause(err)
	}

	group.UpdateUpgradeConfig(config)
	if err := s.groupRepo.Save(ctx, group); err != nil {
		s.helper.Errorw("msg", "update strategy group upgrade config failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update strategy group upgrade config %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateDataSourceConfig(ctx context.Context, uid snowflake.ID, dataSourceUIDs map[snowflake.ID]bool, query string, dataSourceType string) error {
	strategy, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("strategy %s not found", uid)
		}
		return merr.ErrorInternal("get strategy failed").WithCause(err)
	}

	if dataSourceUIDs != nil {
		strategy.UpdateDataSourceUIDs(dataSourceUIDs)
	}
	if query != "" {
		strategy.UpdateQuery(query)
	}
	if dataSourceType != "" {
		strategy.UpdateDataSourceType(dataSourceType)
	}

	if err := s.repo.Save(ctx, strategy); err != nil {
		s.helper.Errorw("msg", "update strategy datasource config failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update strategy datasource config %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateAlertConfig(ctx context.Context, uid snowflake.ID, alertTitle string, alertContent string, alertLevel vobj.AlertLevel, alertPages []string) error {
	strategy, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("strategy %s not found", uid)
		}
		return merr.ErrorInternal("get strategy failed").WithCause(err)
	}

	if alertTitle != "" {
		strategy.UpdateAlertTitle(alertTitle)
	}
	if alertContent != "" {
		strategy.UpdateAlertContent(alertContent)
	}
	if alertLevel != 0 {
		strategy.UpdateAlertLevel(alertLevel)
	}
	if alertPages != nil {
		strategy.UpdateAlertPages(alertPages)
	}

	if err := s.repo.Save(ctx, strategy); err != nil {
		s.helper.Errorw("msg", "update strategy alert config failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update strategy alert config %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateNotifyConfig(ctx context.Context, uid snowflake.ID, receiverUIDs map[snowflake.ID]bool) error {
	// TODO: 需要实现策略与接收对象的关联管理（通过 strategy_receiver 表）
	// 暂时先返回未实现错误
	return merr.ErrorInternal("update notify config not implemented yet")
}

func (s *Service) UpdateDialTestConfig(ctx context.Context, uid snowflake.ID, dialTestType vobj.DialTestType, dialTestTargets map[string]string) error {
	strategy, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("strategy %s not found", uid)
		}
		return merr.ErrorInternal("get strategy failed").WithCause(err)
	}

	if dialTestType != "" {
		strategy.UpdateDialTestType(dialTestType)
	}
	if dialTestTargets != nil {
		strategy.UpdateDialTestTargets(dialTestTargets)
	}

	if err := s.repo.Save(ctx, strategy); err != nil {
		s.helper.Errorw("msg", "update strategy dial test config failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update strategy dial test config %s failed", uid).WithCause(err)
	}
	return nil
}

func (s *Service) UpdateSuppressConfig(ctx context.Context, uid snowflake.ID, suppressType string, suppressConfig string) error {
	strategy, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if merr.IsNotFound(err) {
			return merr.ErrorNotFound("strategy %s not found", uid)
		}
		return merr.ErrorInternal("get strategy failed").WithCause(err)
	}

	if suppressType != "" {
		strategy.UpdateSuppressType(suppressType)
	}
	if suppressConfig != "" {
		strategy.UpdateSuppressConfig(suppressConfig)
	}

	if err := s.repo.Save(ctx, strategy); err != nil {
		s.helper.Errorw("msg", "update strategy suppress config failed", "error", err, "uid", uid)
		return merr.ErrorInternal("update strategy suppress config %s failed", uid).WithCause(err)
	}
	return nil
}
