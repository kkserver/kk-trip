package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RouteStopCreateTaskResult struct {
	app.Result
	Stop *RouteStop `json:"stop,omitempty"`
}

type RouteStopCreateTask struct {
	app.Task
	RouteId   int64   `json:"routeId"`   //开始名
	Title     string  `json:"title"`     //站点备注名
	Longitude float64 `json:"longitude"` //经度
	Latitude  float64 `json:"latitude"`  //纬度
	Type      int     `json:"type"`      //站点类型
	Direction int     `json:"direction"` //方向
	Result    RouteStopCreateTaskResult
}

func (task *RouteStopCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *RouteStopCreateTask) GetInhertType() string {
	return "route"
}

func (task *RouteStopCreateTask) GetClientName() string {
	return "RouteStop.Create"
}
