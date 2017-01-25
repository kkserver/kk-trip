package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RouteCreateTaskResult struct {
	app.Result
	Route *Route `json:"route,omitempty"`
}

type RouteCreateTask struct {
	app.Task
	Start    string  `json:"start"`    //开始名
	End      string  `json:"end"`      //结束名
	Alias    string  `json:"alias"`    //别名
	Tags     string  `json:"tags"`     //搜索标签
	Distance float64 `json:"distance"` //路面距离 km
	Result   RouteCreateTaskResult
}

func (task *RouteCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *RouteCreateTask) GetInhertType() string {
	return "route"
}

func (task *RouteCreateTask) GetClientName() string {
	return "Route.Create"
}
