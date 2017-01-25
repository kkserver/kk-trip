package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RouteRemoveTaskResult struct {
	app.Result
	Route *Route `json:"route,omitempty"`
}

type RouteRemoveTask struct {
	app.Task
	Id     int64 `json:"id,string"`
	Result RouteRemoveTaskResult
}

func (task *RouteRemoveTask) GetResult() interface{} {
	return &task.Result
}

func (task *RouteRemoveTask) GetInhertType() string {
	return "route"
}

func (task *RouteRemoveTask) GetClientName() string {
	return "Route.Remove"
}
