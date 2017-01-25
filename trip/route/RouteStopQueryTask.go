package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RouteStopQueryResult struct {
	app.Result
	Stops []RouteStop `json:"stops,omitempty"`
}

type RouteStopQueryTask struct {
	app.Task
	Id        int64       `json:"id,string"`
	RouteId   int64       `json:"routeId,string"`
	Type      interface{} `json:"type"`      //站点类型
	Direction interface{} `json:"direction"` //方向
	Result    RouteStopQueryResult
}

func (T *RouteStopQueryTask) GetResult() interface{} {
	return &T.Result
}

func (T *RouteStopQueryTask) GetInhertType() string {
	return "route"
}

func (T *RouteStopQueryTask) GetClientName() string {
	return "RouteStop.Query"
}
