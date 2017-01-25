package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RouteStopExchangeTaskResult struct {
	app.Result
}

type RouteStopExchangeTask struct {
	app.Task
	FromId int64 `json:"fromId"`
	ToId   int64 `json:"toId"`
	Result RouteStopExchangeTaskResult
}

func (task *RouteStopExchangeTask) GetResult() interface{} {
	return &task.Result
}

func (task *RouteStopExchangeTask) GetInhertType() string {
	return "route"
}

func (task *RouteStopExchangeTask) GetClientName() string {
	return "RouteStop.Exchange"
}
