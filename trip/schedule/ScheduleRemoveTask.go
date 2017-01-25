package schedule

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ScheduleRemoveTaskResult struct {
	app.Result
	Schedule *Schedule `json:"schedule,omitempty"`
}

type ScheduleRemoveTask struct {
	app.Task
	Id     int64 `json:"id,string"`
	Result ScheduleRemoveTaskResult
}

func (task *ScheduleRemoveTask) GetResult() interface{} {
	return &task.Result
}

func (task *ScheduleRemoveTask) GetInhertType() string {
	return "schedule"
}

func (task *ScheduleRemoveTask) GetClientName() string {
	return "Schedule.Remove"
}
