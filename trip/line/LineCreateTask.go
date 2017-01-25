package line

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type LineCreateTaskResult struct {
	app.Result
	Line *Line `json:"line,omitempty"`
}

type LineCreateTask struct {
	app.Task
	RouteId   int64  `json:"routeId,string"` //路线ID
	Price     int64  `json:"price"`          //原价
	Alias     string `json:"alias"`          //别名
	Direction int    `json:"direction"`      //方向
	Times     string `json:"times"`          //站点时间 09:00,09:20
	Result    LineCreateTaskResult
}

func (task *LineCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *LineCreateTask) GetInhertType() string {
	return "line"
}

func (task *LineCreateTask) GetClientName() string {
	return "Line.Create"
}
