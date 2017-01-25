package schedule

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ScheduleTaskResult struct {
	app.Result
	Schedule *Schedule `json:"schedule,omitempty"`
}

type ScheduleTask struct {
	app.Task
	Id     int64 `json:"id,string"`
	Result ScheduleTaskResult
}

func (task *ScheduleTask) GetResult() interface{} {
	return &task.Result
}

func (task *ScheduleTask) GetInhertType() string {
	return "schedule"
}

func (task *ScheduleTask) GetClientName() string {
	return "Schedule.Get"
}
