package schedule

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ScheduleFailTaskResult struct {
	app.Result
	Schedule *Schedule `json:"schedule,omitempty"`
}

/**
 * 上线
 */
type ScheduleFailTask struct {
	app.Task
	Id     int64 `json:"id,string"`
	Result ScheduleFailTaskResult
}

func (task *ScheduleFailTask) GetResult() interface{} {
	return &task.Result
}

func (task *ScheduleFailTask) GetInhertType() string {
	return "schedule"
}

func (task *ScheduleFailTask) GetClientName() string {
	return "Schedule.Fail"
}
