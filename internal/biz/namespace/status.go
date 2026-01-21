// Package namespace is the namespace aggregate.
package namespace

//go:generate stringer -type=Status -linecomment -output=status_string.go
type Status int8

const (
	StatusUnknown  Status = iota // 未知
	StatusEnabled                // 启用
	StatusDisabled               // 禁用
)

func (s Status) IsEnabled() bool {
	return s == StatusEnabled
}

func (s Status) IsDisabled() bool {
	return s == StatusDisabled
}

