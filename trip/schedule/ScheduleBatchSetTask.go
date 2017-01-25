package schedule

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ScheduleBatchSetTaskResult struct {
	app.Result
}

type ScheduleBatchSetTask struct {
	app.Task
	LineIds     string `json:"lineIds"`     //车次ID, 逗号分隔
	Dates       string `json:"dates"`       //日期, 2016-01-01 逗号分隔
	MaxCount    int    `json:"maxCount"`    //最大售票数
	UMaxCount   int    `json:"umaxCount"`   //用户最大购票数
	AdvanceDays int    `json:"advanceDays"` //提前n天上线
	Result      ScheduleBatchSetTaskResult
}

func (task *ScheduleBatchSetTask) GetResult() interface{} {
	return &task.Result
}

func (task *ScheduleBatchSetTask) GetInhertType() string {
	return "schedule"
}

func (task *ScheduleBatchSetTask) GetClientName() string {
	return "Schedule.BatchSet"
}
