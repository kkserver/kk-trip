package car

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CarQueryCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
}

type CarQueryTaskResult struct {
	app.Result
	Counter *CarQueryCounter `json:"counter,omitempty"`
	Cars    []Car            `json:"cars,omitempty"`
}

type CarQueryTask struct {
	app.Task
	Id          int64  `json:"id"`
	Keyword     string `json:"q"`           //姓名
	LicenceCode string `json:"licenceCode"` //行驶证
	PlateNo     string `json:"plateNo"`     //车牌号 唯一
	Capacity    string `json:"capacity"`    //排量
	OrderBy     string `json:"orderBy"`     // desc, asc , date
	PageIndex   int    `json:"p"`
	PageSize    int    `json:"size"`
	Counter     bool   `json:"counter"`
	Result      CarQueryTaskResult
}

func (task *CarQueryTask) GetResult() interface{} {
	return &task.Result
}

func (task *CarQueryTask) GetInhertType() string {
	return "car"
}

func (task *CarQueryTask) GetClientName() string {
	return "Car.Query"
}
