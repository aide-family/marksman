// Package fileimpl is the file repository for datasource service.
package fileimpl

import (
	datasourcev1 "github.com/aide-family/sovereign/pkg/repo/datasource/v1"
	"github.com/aide-family/sovereign/pkg/repo/datasource/v1/fileimpl/model"
)

func convertDataSourceModel(item *model.DataSourceModel) *datasourcev1.DataSourceModel {
	return &datasourcev1.DataSourceModel{
		Id:           item.ID,
		Uid:          item.UID,
		NamespaceUid: item.NamespaceUID,
		Type:         item.Type,
		Engine:       item.Engine,
		Name:         item.Name,
		Status:       datasourcev1.StatusFromInt8(item.Status),
		Endpoint:     item.Endpoint,
		Description:  item.Description,
		Config:       item.Config,
		Metadata:     item.Metadata,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
		DeletedAt:    item.DeletedAt,
		Creator:      item.Creator,
	}
}

func convertDataSourceItemSelect(item *model.DataSourceModel) *datasourcev1.DataSourceItemSelect {
	return &datasourcev1.DataSourceItemSelect{
		Value:    item.UID,
		Label:    item.Name,
		Disabled: item.Status != int8(datasourcev1.StatusEnabled()),
		Tooltip:  item.Description,
	}
}
