package event

import (
	"time"

	"github.com/aide-family/magicbox/hello"
	"github.com/bwmarrin/snowflake"
)

// ScheduledTask 定时任务实体
type ScheduledTask struct {
	uid             snowflake.ID
	strategyUID     snowflake.ID
	executorType    string // metric/logs/trace/event/dial_test
	status          string // running/stopped
	lastExecuteTime time.Time
	nextExecuteTime time.Time
	createdAt       time.Time
	updatedAt       time.Time
}

func NewScheduledTask(strategyUID snowflake.ID, executorType string) *ScheduledTask {
	var uid snowflake.ID
	node, err := snowflake.NewNode(hello.NodeID())
	if err == nil {
		uid = node.Generate()
	}
	return &ScheduledTask{
		uid:           uid,
		strategyUID:   strategyUID,
		executorType:  executorType,
		status:        "running",
		createdAt:     time.Now(),
		updatedAt:     time.Now(),
	}
}

func (t *ScheduledTask) Start() {
	t.status = "running"
	t.updatedAt = time.Now()
}

func (t *ScheduledTask) Stop() {
	t.status = "stopped"
	t.updatedAt = time.Now()
}

func (t *ScheduledTask) UpdateExecuteTime(lastExecute, nextExecute time.Time) {
	t.lastExecuteTime = lastExecute
	t.nextExecuteTime = nextExecute
	t.updatedAt = time.Now()
}

// FromModel creates a ScheduledTask entity from repository model
func ScheduledTaskFromModel(uid, strategyUID snowflake.ID, executorType, status string, lastExecuteTime, nextExecuteTime time.Time, createdAt, updatedAt time.Time) *ScheduledTask {
	return &ScheduledTask{
		uid:             uid,
		strategyUID:     strategyUID,
		executorType:    executorType,
		status:          status,
		lastExecuteTime: lastExecuteTime,
		nextExecuteTime: nextExecuteTime,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}

// Getters
func (t *ScheduledTask) UID() snowflake.ID        { return t.uid }
func (t *ScheduledTask) StrategyUID() snowflake.ID { return t.strategyUID }
func (t *ScheduledTask) ExecutorType() string     { return t.executorType }
func (t *ScheduledTask) Status() string           { return t.status }
func (t *ScheduledTask) LastExecuteTime() time.Time { return t.lastExecuteTime }
func (t *ScheduledTask) NextExecuteTime() time.Time { return t.nextExecuteTime }
func (t *ScheduledTask) CreatedAt() time.Time       { return t.createdAt }
func (t *ScheduledTask) UpdatedAt() time.Time      { return t.updatedAt }

