package schedule

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ScheduleQueryCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
}

type ScheduleQueryTaskResult struct {
	app.Result
	Counter   *ScheduleQueryCounter `json:"counter,omitempty"`
	Schedules []Schedule            `json:"schedules,omitempty"`
}

type ScheduleQueryTask struct {
	app.Task
	Id        int64       `json:"id,string"`
	LineId    int64       `json:"lineId"`    //车次ID
	StartDate interface{} `json:"startDate"` //日期
	EndDate   interface{} `json:"endDate"`   //日期
	Status    string      `json:"status"`    //状态
	CarId     interface{} `json:"carId"`     //车辆ID
	DriverId  interface{} `json:"driverId"`  //司机ID
	OrderBy   string      `json:"orderBy"`   // desc, asc
	PageIndex int         `json:"p"`
	PageSize  int         `json:"size"`
	Counter   bool        `json:"counter"`
	Result    ScheduleQueryTaskResult
}

func (T *ScheduleQueryTask) GetResult() interface{} {
	return &T.Result
}

func (T *ScheduleQueryTask) GetInhertType() string {
	return "schedule"
}

func (T *ScheduleQueryTask) GetClientName() string {
	return "Schedule.Query"
}
