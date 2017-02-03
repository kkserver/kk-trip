package car

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CarLocationSetTaskResult struct {
	app.Result
}

type CarLocationSetTask struct {
	app.Task

	Id int64 `json:"id"`

	Latitude  float64 `json:"latitude"`  //经度
	Longitude float64 `json:"longitude"` //纬度
	Ip        string  `json:"ip"`        //IP地址

	Result CarLocationSetTaskResult
}

func (task *CarLocationSetTask) GetResult() interface{} {
	return &task.Result
}

func (task *CarLocationSetTask) GetInhertType() string {
	return "car"
}

func (task *CarLocationSetTask) GetClientName() string {
	return "Car.LocationSet"
}
