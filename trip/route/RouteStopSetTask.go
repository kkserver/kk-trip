package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RouteStopSetTaskResult struct {
	app.Result
	Stop *RouteStop `json:"stop,omitempty"`
}

type RouteStopSetTask struct {
	app.Task
	Id        int64       `json:"id"`
	Title     interface{} `json:"title"`     //站点备注名
	Longitude interface{} `json:"longitude"` //经度
	Latitude  interface{} `json:"latitude"`  //纬度
	Body      interface{} `json:"body"`      //站点描述
	Images    interface{} `json:"images"`    //站点图片
	Result    RouteStopSetTaskResult
}

func (task *RouteStopSetTask) GetResult() interface{} {
	return &task.Result
}

func (task *RouteStopSetTask) GetInhertType() string {
	return "route"
}

func (task *RouteStopSetTask) GetClientName() string {
	return "RouteStop.Set"
}
