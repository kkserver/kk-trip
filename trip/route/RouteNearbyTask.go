package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RouteNearbyCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
}

type RouteNearbyTaskResult struct {
	app.Result
	Counter *RouteNearbyCounter `json:"counter,omitempty"`
	Routes  []Route             `json:"routes,omitempty"`
}

type RouteNearbyTask struct {
	app.Task
	Id        int64   `json:"id,string"`
	Status    string  `json:"status"`
	Keyword   string  `json:"q"`
	Longitude float64 `json:"longitude"` //经度
	Latitude  float64 `json:"latitude"`  //纬度
	Distance  float64 `json:"distance"`  //路面距离 km
	Phone     string  `json:"phone"`
	PageIndex int     `json:"p"`
	PageSize  int     `json:"size"`
	Counter   bool    `json:"counter"`
	Result    RouteNearbyTaskResult
}

func (T *RouteNearbyTask) GetResult() interface{} {
	return &T.Result
}

func (T *RouteNearbyTask) GetInhertType() string {
	return "route"
}

func (T *RouteNearbyTask) GetClientName() string {
	return "Route.Nearby"
}
