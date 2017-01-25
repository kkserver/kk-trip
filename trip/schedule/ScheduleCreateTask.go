package schedule

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ScheduleCreateTaskResult struct {
	app.Result
	Schedule *Schedule `json:"schedule,omitempty"`
}

type ScheduleCreateTask struct {
	app.Task
	LineId    int64 `json:"lineId"`    //车次ID
	Date      int64 `json:"date"`      //日期
	MaxCount  int   `json:"maxCount"`  //最大售票数
	UMaxCount int   `json:"umaxCount"` //用户最大购票数
	CarId     int64 `json:"carId"`     //车辆ID
	DriverId  int64 `json:"driverId"`  //司机ID
	InTime    int64 `json:"inTime"`    //上线时间
	Result    ScheduleCreateTaskResult
}

func (task *ScheduleCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *ScheduleCreateTask) GetInhertType() string {
	return "schedule"
}

func (task *ScheduleCreateTask) GetClientName() string {
	return "Schedule.Create"
}
