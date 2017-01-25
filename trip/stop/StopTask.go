package stop

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type StopTaskResult struct {
	app.Result
	Stop *Stop `json:"stop,omitempty"`
}

type StopTask struct {
	app.Task
	Id     int64 `json:"id,string"`
	Result StopTaskResult
}

func (task *StopTask) GetResult() interface{} {
	return &task.Result
}

func (task *StopTask) GetInhertType() string {
	return "stop"
}

func (task *StopTask) GetClientName() string {
	return "Stop.Get"
}
