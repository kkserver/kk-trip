package line

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type LineSetTaskResult struct {
	app.Result
	Line *Line `json:"line,omitempty"`
}

type LineSetTask struct {
	app.Task
	Id        int64       `json:"id,string"`
	Price     interface{} `json:"price"`     //原价
	Alias     interface{} `json:"alias"`     //别名
	Direction interface{} `json:"direction"` //方向
	Times     interface{} `json:"times"`     //站点时间 09:00,09:20
	Status    interface{} `json:"status"`

	CarId    interface{} `json:"carId"`    //车辆ID
	DriverId interface{} `json:"driverId"` //司机ID

	Result LineSetTaskResult
}

func (task *LineSetTask) GetResult() interface{} {
	return &task.Result
}

func (task *LineSetTask) GetInhertType() string {
	return "line"
}

func (task *LineSetTask) GetClientName() string {
	return "Line.Set"
}
