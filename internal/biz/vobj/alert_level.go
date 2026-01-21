// Package vobj is the value object package for the Sovereign service.
package vobj

//go:generate stringer -type=AlertLevel -linecomment -output=alert_level__string.go
type AlertLevel int8

const (
	AlertLevelUnknown AlertLevel = iota // 未知
	AlertLevelP0                        // P0-紧急
	AlertLevelP1                        // P1-严重
	AlertLevelP2                        // P2-警告
	AlertLevelP3                        // P3-信息
)

