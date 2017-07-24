package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RouteSetTaskResult struct {
	app.Result
	Route *Route `json:"route,omitempty"`
}

type RouteSetTask struct {
	app.Task
	Id            int64       `json:"id,string"`
	Start         interface{} `json:"start"`    //开始名
	End           interface{} `json:"end"`      //结束名
	Alias         interface{} `json:"alias"`    //别名
	Tags          interface{} `json:"tags"`     //搜索标签
	Distance      interface{} `json:"distance"` //路面距离 km
	Status        interface{} `json:"status"`
	StartCityId   interface{} `json:"startCityId"`   //开始城市ID
	StartCityPath interface{} `json:"startCityPath"` //开始城市
	EndCityId     interface{} `json:"endCityId"`     //结束城市ID
	EndCityPath   interface{} `json:"endCityPath"`   //结束城市
	WhiteList     interface{} `json:"whiteList"`     //白名单
	Result        RouteSetTaskResult
}

func (task *RouteSetTask) GetResult() interface{} {
	return &task.Result
}

func (task *RouteSetTask) GetInhertType() string {
	return "route"
}

func (task *RouteSetTask) GetClientName() string {
	return "Route.Set"
}
