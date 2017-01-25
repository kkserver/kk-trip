package stop

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type StopSetTaskResult struct {
	app.Result
	Stop *Stop `json:"stop,omitempty"`
}

type StopSetTask struct {
	app.Task
	Id        int64       `json:"id,string"`
	Title     interface{} `json:"string"`    //站点名
	Longitude interface{} `json:"longitude"` //经度
	Latitude  interface{} `json:"latitude"`  //纬度
	Result    StopRemoveTaskResult
}

func (task *StopSetTask) GetResult() interface{} {
	return &task.Result
}

func (task *StopSetTask) GetInhertType() string {
	return "stop"
}

func (task *StopSetTask) GetClientName() string {
	return "Stop.Set"
}
