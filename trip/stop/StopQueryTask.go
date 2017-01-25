package stop

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type StopQueryCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
}

type StopQueryTaskResult struct {
	app.Result
	Counter *StopQueryCounter `json:"counter,omitempty"`
	Stops   []Stop            `json:"stops,omitempty"`
}

type StopQueryTask struct {
	app.Task
	Id        int64  `json:"id,string"`
	Keyword   string `json:"q"`
	OrderBy   string `json:"orderBy"` // desc, asc
	PageIndex int    `json:"p"`
	PageSize  int    `json:"size"`
	Counter   bool   `json:"counter"`
	Result    StopQueryTaskResult
}

func (T *StopQueryTask) GetResult() interface{} {
	return &T.Result
}

func (T *StopQueryTask) GetInhertType() string {
	return "stop"
}

func (T *StopQueryTask) GetClientName() string {
	return "Stop.Query"
}
