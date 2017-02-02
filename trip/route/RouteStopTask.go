package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RouteStopTaskResult struct {
	app.Result
	Stop *RouteStop `json:"stop,omitempty"`
}

type RouteStopTask struct {
	app.Task
	Id     int64 `json:"id"`
	Result RouteStopTaskResult
}

func (task *RouteStopTask) GetResult() interface{} {
	return &task.Result
}

func (task *RouteStopTask) GetInhertType() string {
	return "route"
}

func (task *RouteStopTask) GetClientName() string {
	return "RouteStop.Get"
}
