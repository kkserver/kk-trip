package line

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type LineCountTaskResult struct {
	app.Result
	Count int `json:"count,omitempty"`
}

type LineCountTask struct {
	app.Task
	Id        int64  `json:"id,string"`
	Status    string `json:"status"`
	Keyword   string `json:"q"`
	Direction string `json:"direction"` //方向
	RouteId   int64  `json:"routeId"`
	Result    LineCountTaskResult
}

func (T *LineCountTask) GetResult() interface{} {
	return &T.Result
}

func (T *LineCountTask) GetInhertType() string {
	return "line"
}

func (T *LineCountTask) GetClientName() string {
	return "Line.Count"
}
