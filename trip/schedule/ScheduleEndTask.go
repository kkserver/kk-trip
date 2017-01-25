package schedule

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ScheduleEndTaskResult struct {
	app.Result
	Schedule *Schedule `json:"schedule,omitempty"`
}

/**
 * 已发车
 */
type ScheduleEndTask struct {
	app.Task
	Id     int64 `json:"id,string"`
	Result ScheduleEndTaskResult
}

func (task *ScheduleEndTask) GetResult() interface{} {
	return &task.Result
}

func (task *ScheduleEndTask) GetInhertType() string {
	return "schedule"
}

func (task *ScheduleEndTask) GetClientName() string {
	return "Schedule.End"
}
