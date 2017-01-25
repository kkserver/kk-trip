package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RouteStopRemoveTaskResult struct {
	app.Result
	Stop *RouteStop `json:"stop,omitempty"`
}

type RouteStopRemoveTask struct {
	app.Task
	Id     int64 `json:"id,string"`
	Result RouteStopRemoveTaskResult
}

func (task *RouteStopRemoveTask) GetResult() interface{} {
	return &task.Result
}

func (task *RouteStopRemoveTask) GetInhertType() string {
	return "route"
}

func (task *RouteStopRemoveTask) GetClientName() string {
	return "RouteStop.Remove"
}
