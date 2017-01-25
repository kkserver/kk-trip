package stop

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type StopCreateTaskResult struct {
	app.Result
	Stop *Stop `json:"stop,omitempty"`
}

type StopCreateTask struct {
	app.Task
	Title     string  `json:"string"`    //站点名
	Longitude float64 `json:"longitude"` //经度
	Latitude  float64 `json:"latitude"`  //纬度
	Result    StopCreateTaskResult
}

func (task *StopCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *StopCreateTask) GetInhertType() string {
	return "stop"
}

func (task *StopCreateTask) GetClientName() string {
	return "Stop.Create"
}
