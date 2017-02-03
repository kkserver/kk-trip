package car

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type CarSetTaskResult struct {
	app.Result
	Car *Car `json:"car,omitempty"`
}

type CarSetTask struct {
	app.Task

	Id          int64       `json:"id"`
	LicenceCode interface{} `json:"licenceCode"` //行驶证
	Brand       interface{} `json:"brand"`       //品牌名
	Name        interface{} `json:"name"`        //车辆明
	PlateNo     interface{} `json:"plateNo"`     //车牌号
	Capacity    interface{} `json:"capacity"`    //排量

	Result CarSetTaskResult
}

func (task *CarSetTask) GetResult() interface{} {
	return &task.Result
}

func (task *CarSetTask) GetInhertType() string {
	return "car"
}

func (task *CarSetTask) GetClientName() string {
	return "Car.Set"
}
