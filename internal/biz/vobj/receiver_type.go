// Package vobj is the value object package for the Sovereign service.
package vobj

type ReceiverType string

const (
	ReceiverTypeCommon    ReceiverType = "common"      // 通用接收对象
	ReceiverTypeLabelMatch ReceiverType = "label_match" // label匹配接收对象
)

func (t ReceiverType) String() string {
	return string(t)
}

