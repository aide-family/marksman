// Package vobj is the value object package for the Sovereign service.
package vobj

type NotifyType string

const (
	NotifyTypeSMS   NotifyType = "sms"   // 短信
	NotifyTypePhone NotifyType = "phone" // 电话
	NotifyTypeEmail NotifyType = "email" // 邮件
)

func (t NotifyType) String() string {
	return string(t)
}

