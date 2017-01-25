package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RouteTaskResult struct {
	app.Result
	Route *Route `json:"route,omitempty"`
}

type RouteTask struct {
	app.Task
	Id     int64 `json:"id,string"`
	Result RouteTaskResult
}

func (task *RouteTask) GetResult() interface{} {
	return &task.Result
}

func (task *RouteTask) GetInhertType() string {
	return "route"
}

func (task *RouteTask) GetClientName() string {
	return "Route.Get"
}
