// Package vobj is the value object package for the Sovereign service.
package vobj

type StrategyType string

const (
	StrategyTypeBasic    StrategyType = "basic"     // 基础策略
	StrategyTypeDialTest StrategyType = "dial_test" // 拨测策略
	StrategyTypeSuppress StrategyType = "suppress"  // 抑制策略
)

func (t StrategyType) String() string {
	return string(t)
}

