package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RouteCountTaskResult struct {
	app.Result
	Count int `json:"count,omitempty"`
}

type RouteCountTask struct {
	app.Task
	Id      int64  `json:"id,string"`
	Status  string `json:"status"`
	Keyword string `json:"q"`
	Result  RouteCountTaskResult
}

func (T *RouteCountTask) GetResult() interface{} {
	return &T.Result
}

func (T *RouteCountTask) GetInhertType() string {
	return "route"
}

func (T *RouteCountTask) GetClientName() string {
	return "Route.Count"
}
