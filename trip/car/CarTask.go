package car

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CarTaskResult struct {
	app.Result
	Car *Car `json:"car,omitempty"`
}

type CarTask struct {
	app.Task
	Id          int64  `json:"id"`
	LicenceCode string `json:"licenceCode"` //行驶证
	PlateNo     string `json:"plateNo"`     //车牌号
	Result      CarTaskResult
}

func (task *CarTask) GetResult() interface{} {
	return &task.Result
}

func (task *CarTask) GetInhertType() string {
	return "car"
}

func (task *CarTask) GetClientName() string {
	return "Car.Get"
}
