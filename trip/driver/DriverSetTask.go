package driver

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type DriverSetTaskResult struct {
	app.Result
	Driver *Driver `json:"driver,omitempty"`
}

type DriverSetTask struct {
	app.Task
	Id          int64       `json:"id"`
	Name        interface{} `json:"name"`        //姓名
	Phone       interface{} `json:"phone"`       //手机号 	唯一
	Code        interface{} `json:"code"`        //身份证 	唯一
	LicenceCode interface{} `json:"licenceCode"` //驾驶证		唯一
	Result      DriverSetTaskResult
}

func (task *DriverSetTask) GetResult() interface{} {
	return &task.Result
}

func (task *DriverSetTask) GetInhertType() string {
	return "driver"
}

func (task *DriverSetTask) GetClientName() string {
	return "Driver.Set"
}
