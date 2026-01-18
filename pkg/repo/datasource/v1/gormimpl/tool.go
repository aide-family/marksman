// Package gormimpl provides helper conversions for datasource models.
package gormimpl

import (
	datasourcev1 "github.com/aide-family/sovereign/pkg/repo/datasource/v1"
	"github.com/aide-family/sovereign/pkg/repo/datasource/v1/gormimpl/model"
	"gorm.io/gorm"
)

func ConvertDataSourceModel(item *model.DataSource) *datasourcev1.DataSourceModel {
	if item == nil {
		return nil
	}
	var config map[string]string
	if item.Config != nil {
		config = item.Config.Map()
	}
	var metadata map[string]string
	if item.Metadata != nil {
		metadata = item.Metadata.Map()
	}
	return &datasourcev1.DataSourceModel{
		Id:           item.ID,
		Uid:          item.UID.Int64(),
		NamespaceUid: item.NamespaceUID.Int64(),
		Type:         item.Type,
		Engine:       item.Engine,
		Name:         item.Name,
		Status:       datasourcev1.StatusFromInt8(item.Status),
		Endpoint:     item.Endpoint,
		Description:  item.Description,
		Config:       config,
		Metadata:     metadata,
		CreatedAt:    item.CreatedAt.Unix(),
		UpdatedAt:    item.UpdatedAt.Unix(),
		DeletedAt:    item.DeletedAt.Time.Unix(),
		Creator:      item.Creator.Int64(),
	}
}

func ConvertDataSourceItemSelect(item *model.DataSource) *datasourcev1.DataSourceItemSelect {
	if item == nil {
		return nil
	}
	return &datasourcev1.DataSourceItemSelect{
		Value:    item.UID.Int64(),
		Label:    item.Name,
		Disabled: item.Status != int8(datasourcev1.StatusEnabled()),
		Tooltip:  item.Description,
	}
}

func convertResultInfo(result *gorm.DB) *datasourcev1.ResultInfo {
	if result == nil {
		return &datasourcev1.ResultInfo{}
	}
	return &datasourcev1.ResultInfo{
		RowsAffected: result.RowsAffected,
		Error:        errString(result.Error),
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
