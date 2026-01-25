// Package datasourcev1 provides datasource helpers.
package datasourcev1

import "github.com/aide-family/sovereign/pkg/enum"

func StatusFromInt8(value int8) enum.GlobalStatus {
	switch value {
	case int8(enum.GlobalStatus_ENABLED):
		return enum.GlobalStatus_ENABLED
	case int8(enum.GlobalStatus_DISABLED):
		return enum.GlobalStatus_DISABLED
	default:
		return enum.GlobalStatus_GlobalStatus_UNKNOWN
	}
}

func StatusEnabled() enum.GlobalStatus {
	return enum.GlobalStatus_ENABLED
}
