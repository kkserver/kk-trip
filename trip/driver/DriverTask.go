package driver

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type DriverTaskResult struct {
	app.Result
	Driver *Driver `json:"driver,omitempty"`
}

type DriverTask struct {
	app.Task
	Id          int64  `json:"id"`
	Phone       string `json:"string"`      //手机号	唯一
	Code        string `json:"code"`        //身份证	唯一
	LicenceCode string `json:"licenceCode"` //驾驶证 唯一
	Result      DriverTaskResult
}

func (task *DriverTask) GetResult() interface{} {
	return &task.Result
}

func (task *DriverTask) GetInhertType() string {
	return "driver"
}

func (task *DriverTask) GetClientName() string {
	return "Driver.Get"
}
