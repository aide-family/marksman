package strategy

import (
	"time"
	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
	"github.com/aide-family/sovereign/internal/biz/vobj"
	"github.com/aide-family/sovereign/pkg/merr"
)

type Receiver struct {
	uid          snowflake.ID
	namespaceUID snowflake.ID
	typ          vobj.ReceiverType
	name         string
	description  string
	userIDs      map[snowflake.ID]bool
	labelMatch   map[string]string
	notifyTypes  []vobj.NotifyType
	createdAt    time.Time
	updatedAt    time.Time
}

func NewReceiver(namespaceUID snowflake.ID, typ vobj.ReceiverType, name string) *Receiver {
	var uid snowflake.ID
	node, err := snowflake.NewNode(hello.NodeID())
	if err == nil {
		uid = node.Generate()
	}
	return &Receiver{
		uid:          uid,
		namespaceUID: namespaceUID,
		typ:          typ,
		name:         name,
		userIDs:      make(map[snowflake.ID]bool),
		labelMatch:   make(map[string]string),
		notifyTypes:  make([]vobj.NotifyType, 0),
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
	}
}

func (r *Receiver) UpdateName(name string) error {
	if name == "" {
		return merr.ErrorParams("name cannot be empty")
	}
	r.name = name
	r.updatedAt = time.Now()
	return nil
}

func (r *Receiver) UpdateDescription(description string) {
	r.description = description
	r.updatedAt = time.Now()
}

func (r *Receiver) UpdateUserIDs(userIDs map[snowflake.ID]bool) {
	r.userIDs = userIDs
	r.updatedAt = time.Now()
}

func (r *Receiver) UpdateLabelMatch(labelMatch map[string]string) {
	r.labelMatch = labelMatch
	r.updatedAt = time.Now()
}

func (r *Receiver) UpdateNotifyTypes(notifyTypes []vobj.NotifyType) {
	r.notifyTypes = notifyTypes
	r.updatedAt = time.Now()
}

// FromModel creates a Receiver entity from repository model
func ReceiverFromModel(uid, namespaceUID snowflake.ID, typ vobj.ReceiverType, name, description string, userIDs map[snowflake.ID]bool, labelMatch map[string]string, notifyTypes []vobj.NotifyType, createdAt, updatedAt time.Time) *Receiver {
	return &Receiver{
		uid:          uid,
		namespaceUID: namespaceUID,
		typ:          typ,
		name:         name,
		description:  description,
		userIDs:      userIDs,
		labelMatch:   labelMatch,
		notifyTypes:  notifyTypes,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Getters
func (r *Receiver) UID() snowflake.ID                { return r.uid }
func (r *Receiver) NamespaceUID() snowflake.ID       { return r.namespaceUID }
func (r *Receiver) Type() vobj.ReceiverType          { return r.typ }
func (r *Receiver) Name() string                     { return r.name }
func (r *Receiver) Description() string              { return r.description }
func (r *Receiver) UserIDs() map[snowflake.ID]bool   { return r.userIDs }
func (r *Receiver) LabelMatch() map[string]string    { return r.labelMatch }
func (r *Receiver) NotifyTypes() []vobj.NotifyType   { return r.notifyTypes }
func (r *Receiver) CreatedAt() time.Time             { return r.createdAt }
func (r *Receiver) UpdatedAt() time.Time             { return r.updatedAt }

