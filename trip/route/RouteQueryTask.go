package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type RouteQueryCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
}

type RouteQueryTaskResult struct {
	app.Result
	Counter *RouteQueryCounter `json:"counter,omitempty"`
	Routes  []Route            `json:"routes,omitempty"`
}

type RouteQueryTask struct {
	app.Task
	Id        int64  `json:"id,string"`
	Status    string `json:"status"`
	Keyword   string `json:"q"`
	OrderBy   string `json:"orderBy"` // desc, asc
	PageIndex int    `json:"p"`
	PageSize  int    `json:"size"`
	Counter   bool   `json:"counter"`
	Result    RouteQueryTaskResult
}

func (T *RouteQueryTask) GetResult() interface{} {
	return &T.Result
}

func (T *RouteQueryTask) GetInhertType() string {
	return "route"
}

func (T *RouteQueryTask) GetClientName() string {
	return "Route.Query"
}
