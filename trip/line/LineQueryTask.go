package line

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type LineQueryCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
	RowCount  int `json:"rowCount"`
}

type LineQueryTaskResult struct {
	app.Result
	Counter *LineQueryCounter `json:"counter,omitempty"`
	Lines   []Line            `json:"lines,omitempty"`
}

type LineQueryTask struct {
	app.Task
	Id        int64  `json:"id,string"`
	Status    string `json:"status"`
	Keyword   string `json:"q"`
	Direction string `json:"direction"` //方向
	RouteId   int64  `json:"routeId"`
	OrderBy   string `json:"orderBy"` // desc, asc
	PageIndex int    `json:"p"`
	PageSize  int    `json:"size"`
	Counter   bool   `json:"counter"`
	Result    LineQueryTaskResult
}

func (T *LineQueryTask) GetResult() interface{} {
	return &T.Result
}

func (T *LineQueryTask) GetInhertType() string {
	return "route"
}

func (T *LineQueryTask) GetClientName() string {
	return "Line.Query"
}
