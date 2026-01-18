// Package fileimpl is the implementation of the file repository for datasource service.
package fileimpl

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
	"github.com/fsnotify/fsnotify"
	klog "github.com/go-kratos/kratos/v2/log"
	"go.yaml.in/yaml/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/aide-family/sovereign/pkg/config"
	"github.com/aide-family/sovereign/pkg/enum"
	"github.com/aide-family/sovereign/pkg/merr"
	"github.com/aide-family/sovereign/pkg/repo"
	datasourcev1 "github.com/aide-family/sovereign/pkg/repo/datasource/v1"
	"github.com/aide-family/sovereign/pkg/repo/datasource/v1/fileimpl/model"
)

func init() {
	repo.RegisterDataSourceV1Factory(config.DataSourceConfig_FILE, NewFileRepository)
}

func NewFileRepository(c *config.DataSourceConfig) (datasourcev1.Repository, func() error, error) {
	fileConfig := &config.FileConfig{}
	if c != nil && c.GetOptions() != nil {
		if err := anypb.UnmarshalTo(c.GetOptions(), fileConfig, proto.UnmarshalOptions{Merge: true}); err != nil {
			return nil, nil, merr.ErrorInternalServer("unmarshal file config failed: %v", err)
		}
	}

	// ensure directory exists
	if err := os.MkdirAll(fileConfig.Path, 0755); err != nil {
		return nil, nil, merr.ErrorInternalServer("create directory failed: %v", err)
	}

	tmpFilepath := filepath.Join(fileConfig.Path, fmt.Sprintf("%s.tmp", fileConfig.Filename))
	path := filepath.Join(fileConfig.Path, fileConfig.Filename)
	node, err := snowflake.NewNode(hello.NodeID())
	if err != nil {
		return nil, nil, err
	}
	f := &fileRepository{
		fileConfig:      fileConfig,
		tmpFilepath:     tmpFilepath,
		filepath:        path,
		stopChan:        make(chan struct{}),
		storageInterval: fileConfig.StorageInterval.AsDuration(),
		node:            node,
		dataSources:     make([]*model.DataSourceModel, 0),
	}
	if err := f.load(); err != nil {
		return nil, nil, err
	}
	f.watch()
	return f, func() error {
		close(f.stopChan)
		return f.save()
	}, nil
}

type fileRepository struct {
	fileConfig      *config.FileConfig
	tmpFilepath     string
	filepath        string
	mu              sync.RWMutex
	dataSources     []*model.DataSourceModel
	nextID          uint32
	stopChan        chan struct{}
	storageInterval time.Duration
	changed         bool
	node            *snowflake.Node
}

func (f *fileRepository) load() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, err := os.Stat(f.filepath); os.IsNotExist(err) {
		f.dataSources = make([]*model.DataSourceModel, 0)
		f.nextID = 0
		return nil
	}

	file, err := os.Open(f.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	var items []*model.DataSourceModel
	if err := yaml.NewDecoder(file).Decode(&items); err != nil {
		if err == io.EOF {
			f.dataSources = make([]*model.DataSourceModel, 0)
			f.nextID = 0
			return nil
		}
		return err
	}

	if len(items) == 0 {
		f.dataSources = make([]*model.DataSourceModel, 0)
		f.nextID = 0
		return nil
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	f.nextID = items[len(items)-1].ID
	for _, item := range items {
		if item.ID == 0 {
			f.nextID++
			item.ID = f.nextID
		}
		if item.UID == 0 {
			item.UID = f.node.Generate().Int64()
		}
	}

	f.dataSources = items
	return nil
}

func (f *fileRepository) save() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.changed = false
	file, err := os.Create(f.tmpFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := yaml.NewEncoder(file).Encode(f.dataSources); err != nil {
		return err
	}
	if err := os.Rename(f.tmpFilepath, f.filepath); err != nil {
		return err
	}
	klog.Debugw("msg", "save data sources to file", "filepath", f.filepath)
	return nil
}

func (f *fileRepository) watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		klog.Errorw("msg", "create watcher failed", "error", err)
		return
	}
	defer watcher.Close()
	watcher.Add(f.filepath)
	go func() {
		ticker := time.NewTicker(f.storageInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if f.changed {
					f.save()
				}
			case err := <-watcher.Errors:
				if err != nil {
					klog.Warnw("msg", "watch file failed", "error", err)
				}
			case <-f.stopChan:
				klog.Debugw("msg", "stop watch data sources")
				return
			}
		}
	}()
}

// CreateDataSource implements [datasourcev1.Repository].
func (f *fileRepository) CreateDataSource(ctx context.Context, req *datasourcev1.CreateDataSourceRequest) (*datasourcev1.DataSourceModel, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.changed = true
	f.nextID++
	now := time.Now().Unix()
	item := &model.DataSourceModel{
		ID:           f.nextID,
		UID:          f.node.Generate().Int64(),
		NamespaceUID: req.NamespaceUid,
		Type:         req.Type,
		Engine:       req.Engine,
		Name:         req.Name,
		Status:       int8(req.Status),
		Endpoint:     req.Endpoint,
		Description:  req.Description,
		Config:       req.Config,
		Metadata:     req.Metadata,
		CreatedAt:    now,
		UpdatedAt:    now,
		DeletedAt:    0,
		Creator:      f.node.Generate().Int64(),
	}
	f.dataSources = append(f.dataSources, item)
	return convertDataSourceModel(item), nil
}

// GetDataSource implements [datasourcev1.Repository].
func (f *fileRepository) GetDataSource(ctx context.Context, req *datasourcev1.GetDataSourceRequest) (*datasourcev1.DataSourceModel, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, item := range f.dataSources {
		if item.UID == req.Uid {
			return convertDataSourceModel(item), nil
		}
	}
	return nil, merr.ErrorNotFound("data source %d not found", req.Uid)
}

// UpdateDataSource implements [datasourcev1.Repository].
func (f *fileRepository) UpdateDataSource(ctx context.Context, req *datasourcev1.UpdateDataSourceRequest) (*datasourcev1.ResultInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, item := range f.dataSources {
		if item.UID != req.Uid {
			continue
		}
		f.changed = true
		item.Name = req.Name
		item.Status = int8(req.Status)
		item.Endpoint = req.Endpoint
		item.Description = req.Description
		item.Config = req.Config
		item.Metadata = req.Metadata
		item.UpdatedAt = time.Now().Unix()
		return &datasourcev1.ResultInfo{RowsAffected: 1}, nil
	}
	return &datasourcev1.ResultInfo{RowsAffected: 0}, nil
}

// UpdateDataSourceStatus implements [datasourcev1.Repository].
func (f *fileRepository) UpdateDataSourceStatus(ctx context.Context, req *datasourcev1.UpdateDataSourceStatusRequest) (*datasourcev1.ResultInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, item := range f.dataSources {
		if item.UID != req.Uid {
			continue
		}
		f.changed = true
		item.Status = int8(req.Status)
		item.UpdatedAt = time.Now().Unix()
		return &datasourcev1.ResultInfo{RowsAffected: 1}, nil
	}
	return &datasourcev1.ResultInfo{RowsAffected: 0}, nil
}

// DeleteDataSource implements [datasourcev1.Repository].
func (f *fileRepository) DeleteDataSource(ctx context.Context, req *datasourcev1.DeleteDataSourceRequest) (*datasourcev1.ResultInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	index := -1
	for i, item := range f.dataSources {
		if item.UID == req.Uid {
			index = i
			break
		}
	}
	if index < 0 {
		return &datasourcev1.ResultInfo{RowsAffected: 0}, nil
	}
	f.changed = true
	f.dataSources = append(f.dataSources[:index], f.dataSources[index+1:]...)
	return &datasourcev1.ResultInfo{RowsAffected: 1}, nil
}

// ListDataSource implements [datasourcev1.Repository].
func (f *fileRepository) ListDataSource(ctx context.Context, req *datasourcev1.ListDataSourceRequest) (*datasourcev1.ListDataSourceResponse, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	items := make([]*model.DataSourceModel, 0, len(f.dataSources))
	for _, item := range f.dataSources {
		if req.NamespaceUid > 0 && item.NamespaceUID != req.NamespaceUid {
			continue
		}
		if req.Status > enum.GlobalStatus_GlobalStatus_UNKNOWN && item.Status != int8(req.Status) {
			continue
		}
		if req.Type != "" && item.Type != req.Type {
			continue
		}
		if req.Engine != "" && item.Engine != req.Engine {
			continue
		}
		if req.Keyword != "" {
			if !strings.Contains(strings.ToLower(item.Name), strings.ToLower(req.Keyword)) {
				continue
			}
		}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		less := func(a, b *model.DataSourceModel) bool {
			switch req.OrderBy {
			case datasourcev1.Field_ID:
				return a.ID < b.ID
			case datasourcev1.Field_UID:
				return a.UID < b.UID
			case datasourcev1.Field_NAME:
				return a.Name < b.Name
			case datasourcev1.Field_CREATED_AT:
				return a.CreatedAt < b.CreatedAt
			default:
				return a.UID < b.UID
			}
		}
		if req.Order == datasourcev1.Order_DESC {
			return !less(items[i], items[j])
		}
		return less(items[i], items[j])
	})

	page := req.Page
	pageSize := req.PageSize
	if page <= 0 || pageSize <= 0 {
		page = 1
		pageSize = int32(len(items))
	}
	start := int((page - 1) * pageSize)
	if start > len(items) {
		start = len(items)
	}
	end := int(page * pageSize)
	if end > len(items) {
		end = len(items)
	}

	responseItems := make([]*datasourcev1.DataSourceModel, 0, end-start)
	for _, item := range items[start:end] {
		responseItems = append(responseItems, convertDataSourceModel(item))
	}

	return &datasourcev1.ListDataSourceResponse{
		Items:    responseItems,
		Total:    int64(len(items)),
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// SelectDataSource implements [datasourcev1.Repository].
func (f *fileRepository) SelectDataSource(ctx context.Context, req *datasourcev1.SelectDataSourceRequest) (*datasourcev1.SelectDataSourceResponse, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	items := make([]*model.DataSourceModel, 0, len(f.dataSources))
	for _, item := range f.dataSources {
		if req.NamespaceUid > 0 && item.NamespaceUID != req.NamespaceUid {
			continue
		}
		if req.Status > enum.GlobalStatus_GlobalStatus_UNKNOWN && item.Status != int8(req.Status) {
			continue
		}
		if req.Type != "" && item.Type != req.Type {
			continue
		}
		if req.Engine != "" && item.Engine != req.Engine {
			continue
		}
		if req.Keyword != "" {
			if !strings.Contains(strings.ToLower(item.Name), strings.ToLower(req.Keyword)) {
				continue
			}
		}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return cmp.Compare(items[i].UID, items[j].UID) < 0
	})
	if req.Order == datasourcev1.Order_DESC {
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].UID > items[j].UID
		})
	}

	filtered := make([]*model.DataSourceModel, 0, len(items))
	for _, item := range items {
		if req.LastUID > 0 {
			if req.Order == datasourcev1.Order_DESC && item.UID >= req.LastUID {
				continue
			}
			if req.Order == datasourcev1.Order_ASC && item.UID <= req.LastUID {
				continue
			}
		}
		filtered = append(filtered, item)
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 20
	}
	if limit > len(filtered) {
		limit = len(filtered)
	}
	selected := filtered[:limit]

	result := make([]*datasourcev1.DataSourceItemSelect, 0, len(selected))
	for _, item := range selected {
		result = append(result, convertDataSourceItemSelect(item))
	}

	lastUID := int64(0)
	if len(selected) > 0 {
		lastUID = selected[len(selected)-1].UID
	}

	return &datasourcev1.SelectDataSourceResponse{
		Items:   result,
		Total:   int64(len(filtered)),
		LastUID: lastUID,
		HasMore: len(filtered) > len(selected),
	}, nil
}
