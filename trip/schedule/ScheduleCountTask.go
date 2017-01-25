package schedule

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ScheduleCountDate struct {
	Date      int64 `json:"date"`
	LineCount int   `json:"lineCount"`
}

type ScheduleCountTaskResult struct {
	app.Result
	Dates []ScheduleCountDate `json:"dates,omitempty"`
}

type ScheduleCountTask struct {
	app.Task
	StartDate interface{} `json:"startDate"` //日期
	EndDate   interface{} `json:"endDate"`   //日期
	Status    string      `json:"status"`    //状态
	CarId     interface{} `json:"carId"`     //车辆ID
	DriverId  interface{} `json:"driverId"`  //司机ID
	Result    ScheduleCountTaskResult
}

func (T *ScheduleCountTask) GetResult() interface{} {
	return &T.Result
}

func (T *ScheduleCountTask) GetInhertType() string {
	return "schedule"
}

func (T *ScheduleCountTask) GetClientName() string {
	return "Schedule.Count"
}
