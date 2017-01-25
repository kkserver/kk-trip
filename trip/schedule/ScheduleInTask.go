package schedule

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ScheduleInTaskResult struct {
	app.Result
	Schedule *Schedule `json:"schedule,omitempty"`
}

/**
 * 上线
 */
type ScheduleInTask struct {
	app.Task
	Id     int64 `json:"id,string"`
	Result ScheduleInTaskResult
}

func (task *ScheduleInTask) GetResult() interface{} {
	return &task.Result
}

func (task *ScheduleInTask) GetInhertType() string {
	return "schedule"
}

func (task *ScheduleInTask) GetClientName() string {
	return "Schedule.In"
}
