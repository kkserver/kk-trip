package suggest

import (
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-trip/trip/line"
)

type SuggestLineTaskResult struct {
	app.Result
	Lines []line.Line `json:"lines,omitempty"`
}

type SuggestLineTask struct {
	app.Task
	Uid       int64   `json:"uid,string"`
	Phone     string  `json:"phone"`
	Status    string  `json:"status"`
	Longitude float64 `json:"longitude"` //经度
	Latitude  float64 `json:"latitude"`  //纬度
	Distance  float64 `json:"distance"`  //路面距离 km
	Direction string  `json:"direction"` //方向
	Limit     int     `json:"limit"`
	Result    SuggestLineTaskResult
}

func (T *SuggestLineTask) GetResult() interface{} {
	return &T.Result
}

func (T *SuggestLineTask) GetInhertType() string {
	return "suggest"
}

func (T *SuggestLineTask) GetClientName() string {
	return "Suggest.Line"
}
