// Package vobj is the value object package for the Sovereign service.
package vobj

//go:generate stringer -type=UpgradeMode -linecomment -output=upgrade_mode__string.go
type UpgradeMode int8

const (
	UpgradeModeNone   UpgradeMode = iota // 不升级
	UpgradeModeAuto                      // 自动升级
	UpgradeModeManual                    // 手动升级
)

