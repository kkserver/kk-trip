package driver

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type DriverCreateTaskResult struct {
	app.Result
	Driver *Driver `json:"driver,omitempty"`
}

type DriverCreateTask struct {
	app.Task
	Name        string `json:"name"`        //姓名
	Phone       string `json:"phone"`       //手机号	唯一
	Code        string `json:"code"`        //身份证	唯一
	LicenceCode string `json:"licenceCode"` //驾驶证 唯一
	Result      DriverCreateTaskResult
}

func (task *DriverCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *DriverCreateTask) GetInhertType() string {
	return "driver"
}

func (task *DriverCreateTask) GetClientName() string {
	return "Driver.Create"
}
