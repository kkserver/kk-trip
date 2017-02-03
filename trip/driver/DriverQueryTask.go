package driver

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type DriverQueryCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
}

type DriverQueryTaskResult struct {
	app.Result
	Counter *DriverQueryCounter `json:"counter,omitempty"`
	Drivers []Driver            `json:"drivers,omitempty"`
}

type DriverQueryTask struct {
	app.Task
	Id          int64  `json:"id"`
	Keyword     string `json:"q"`           //姓名
	Name        string `json:"name"`        //姓名
	Phone       string `json:"phone"`       //手机号 	唯一
	Code        string `json:"code"`        //身份证 	唯一
	LicenceCode string `json:"licenceCode"` //驾驶证		唯一
	OrderBy     string `json:"orderBy"`     // desc, asc , date
	PageIndex   int    `json:"p"`
	PageSize    int    `json:"size"`
	Counter     bool   `json:"counter"`
	Result      DriverQueryTaskResult
}

func (task *DriverQueryTask) GetResult() interface{} {
	return &task.Result
}

func (task *DriverQueryTask) GetInhertType() string {
	return "driver"
}

func (task *DriverQueryTask) GetClientName() string {
	return "Driver.Query"
}
