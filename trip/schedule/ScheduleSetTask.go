package schedule

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ScheduleSetTaskResult struct {
	app.Result
	Schedule *Schedule `json:"schedule,omitempty"`
}

type ScheduleSetTask struct {
	app.Task
	Id        int64       `json:"id,string"`
	MaxCount  interface{} `json:"maxCount"`  //最大售票数
	UMaxCount interface{} `json:"umaxCount"` //用户最大购票数
	CarId     interface{} `json:"carId"`     //车辆ID
	DriverId  interface{} `json:"driverId"`  //司机ID
	InTime    interface{} `json:"inTime"`    //上线时间
	Result    ScheduleSetTaskResult
}

func (task *ScheduleSetTask) GetResult() interface{} {
	return &task.Result
}

func (task *ScheduleSetTask) GetInhertType() string {
	return "schedule"
}

func (task *ScheduleSetTask) GetClientName() string {
	return "Schedule.Set"
}
