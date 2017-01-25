package schedule

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ScheduleOffTaskResult struct {
	app.Result
	Schedule *Schedule `json:"schedule,omitempty"`
}

/**
 * 下线
 */
type ScheduleOffTask struct {
	app.Task
	Id     int64 `json:"id,string"`
	Result ScheduleOffTaskResult
}

func (task *ScheduleOffTask) GetResult() interface{} {
	return &task.Result
}

func (task *ScheduleOffTask) GetInhertType() string {
	return "schedule"
}

func (task *ScheduleOffTask) GetClientName() string {
	return "Schedule.Off"
}
