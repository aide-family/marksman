package event

import (
	"time"

	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/vobj"
)

// NotificationHistory 通知历史实体
type NotificationHistory struct {
	uid             snowflake.ID
	notificationTime time.Time
	status          string // success/failed
	notifyType      vobj.NotifyType
	receiver        string
	receiverType    vobj.ReceiverType
	content         string
	alertUID        snowflake.ID
	retryCount      int32
	createdAt       time.Time
	updatedAt       time.Time
}

func NewNotificationHistory(notifyType vobj.NotifyType, receiver string, receiverType vobj.ReceiverType, content string, alertUID snowflake.ID) *NotificationHistory {
	var uid snowflake.ID
	node, err := snowflake.NewNode(hello.NodeID())
	if err == nil {
		uid = node.Generate()
	}
	return &NotificationHistory{
		uid:              uid,
		notificationTime: time.Now(),
		status:           "pending",
		notifyType:       notifyType,
		receiver:         receiver,
		receiverType:     receiverType,
		content:          content,
		alertUID:         alertUID,
		retryCount:       0,
		createdAt:        time.Now(),
		updatedAt:        time.Now(),
	}
}

func (n *NotificationHistory) MarkSuccess() {
	n.status = "success"
	n.updatedAt = time.Now()
}

func (n *NotificationHistory) MarkFailed() {
	n.status = "failed"
	n.updatedAt = time.Now()
}

func (n *NotificationHistory) Retry() {
	if n.retryCount < 3 {
		n.retryCount++
		n.status = "pending"
		n.updatedAt = time.Now()
	}
}

// FromModel creates a NotificationHistory entity from repository model
func NotificationHistoryFromModel(uid snowflake.ID, notificationTime time.Time, status string, notifyType vobj.NotifyType, receiver string, receiverType vobj.ReceiverType, content string, alertUID snowflake.ID, retryCount int32, createdAt, updatedAt time.Time) *NotificationHistory {
	return &NotificationHistory{
		uid:              uid,
		notificationTime: notificationTime,
		status:           status,
		notifyType:       notifyType,
		receiver:         receiver,
		receiverType:     receiverType,
		content:          content,
		alertUID:         alertUID,
		retryCount:       retryCount,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}

// Getters
func (n *NotificationHistory) UID() snowflake.ID        { return n.uid }
func (n *NotificationHistory) NotificationTime() time.Time { return n.notificationTime }
func (n *NotificationHistory) Status() string          { return n.status }
func (n *NotificationHistory) NotifyType() vobj.NotifyType { return n.notifyType }
func (n *NotificationHistory) Receiver() string        { return n.receiver }
func (n *NotificationHistory) ReceiverType() vobj.ReceiverType { return n.receiverType }
func (n *NotificationHistory) Content() string         { return n.content }
func (n *NotificationHistory) AlertUID() snowflake.ID  { return n.alertUID }
func (n *NotificationHistory) RetryCount() int32       { return n.retryCount }
func (n *NotificationHistory) CreatedAt() time.Time   { return n.createdAt }
func (n *NotificationHistory) UpdatedAt() time.Time    { return n.updatedAt }

