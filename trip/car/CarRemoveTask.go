package car

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CarRemoveTaskResult struct {
	app.Result
	Car *Car `json:"car,omitempty"`
}

type CarRemoveTask struct {
	app.Task
	Id     int64 `json:"id"`
	Result CarRemoveTaskResult
}

func (task *CarRemoveTask) GetResult() interface{} {
	return &task.Result
}

func (task *CarRemoveTask) GetInhertType() string {
	return "car"
}

func (task *CarRemoveTask) GetClientName() string {
	return "Car.Remove"
}
