package driver

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type DriverRemoveTaskResult struct {
	app.Result
	Driver *Driver `json:"driver,omitempty"`
}

type DriverRemoveTask struct {
	app.Task
	Id     int64 `json:"id"`
	Result DriverRemoveTaskResult
}

func (task *DriverRemoveTask) GetResult() interface{} {
	return &task.Result
}

func (task *DriverRemoveTask) GetInhertType() string {
	return "driver"
}

func (task *DriverRemoveTask) GetClientName() string {
	return "Driver.Remove"
}
