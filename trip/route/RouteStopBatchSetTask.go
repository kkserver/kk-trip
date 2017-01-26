package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RouteStopBatch struct {
	Title     string  `json:"title"`     //站点备注名
	Longitude float64 `json:"lng"`       //经度
	Latitude  float64 `json:"lat"`       //纬度
	Type      int     `json:"type"`      //站点类型
	Direction int     `json:"direction"` //方向
}

type RouteStopBatchSetTaskResult struct {
	app.Result
	Stops []RouteStop `json:"stops,omitempty"`
}

type RouteStopBatchSetTask struct {
	app.Task
	RouteId int64            `json:"routeId,string"`
	Stops   []RouteStopBatch `json:"stops"`
	Result  RouteStopBatchSetTaskResult
}

func (task *RouteStopBatchSetTask) GetResult() interface{} {
	return &task.Result
}

func (task *RouteStopBatchSetTask) GetInhertType() string {
	return "route"
}

func (task *RouteStopBatchSetTask) GetClientName() string {
	return "Route.BatchSet"
}
