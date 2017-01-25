package schedule

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ScheduleStartTaskResult struct {
	app.Result
	Schedule *Schedule `json:"schedule,omitempty"`
}

/**
 * 已发车
 */
type ScheduleStartTask struct {
	app.Task
	Id     int64 `json:"id,string"`
	Result ScheduleStartTaskResult
}

func (task *ScheduleStartTask) GetResult() interface{} {
	return &task.Result
}

func (task *ScheduleStartTask) GetInhertType() string {
	return "schedule"
}

func (task *ScheduleStartTask) GetClientName() string {
	return "Schedule.Start"
}
