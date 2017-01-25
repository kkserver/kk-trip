package stop

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type StopRemoveTaskResult struct {
	app.Result
	Stop *Stop `json:"stop,omitempty"`
}

type StopRemoveTask struct {
	app.Task
	Id     int64 `json:"id,string"`
	Result StopRemoveTaskResult
}

func (task *StopRemoveTask) GetResult() interface{} {
	return &task.Result
}

func (task *StopRemoveTask) GetInhertType() string {
	return "stop"
}

func (task *StopRemoveTask) GetClientName() string {
	return "Stop.Remove"
}
