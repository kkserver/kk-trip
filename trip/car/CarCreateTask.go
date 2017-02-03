package car

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CarCreateTaskResult struct {
	app.Result
	Car *Car `json:"car,omitempty"`
}

type CarCreateTask struct {
	app.Task

	LicenceCode string `json:"licenceCode"` //行驶证

	Brand    string `json:"brand"`    //品牌名
	Name     string `json:"name"`     //车辆明
	PlateNo  string `json:"plateNo"`  //车牌号
	Capacity string `json:"capacity"` //排量

	Result CarCreateTaskResult
}

func (task *CarCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *CarCreateTask) GetInhertType() string {
	return "car"
}

func (task *CarCreateTask) GetClientName() string {
	return "Car.Create"
}
