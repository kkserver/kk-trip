package stop

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type StopNearbyCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
}

type NearbyStop struct {
	Stop
	Distance float64 `json:"distance"`
}

type NearbyStopSlice []NearbyStop

type StopNearbyTaskResult struct {
	app.Result
	Counter *StopNearbyCounter `json:"counter,omitempty"`
	Stops   NearbyStopSlice    `json:"stops,omitempty"`
}

type StopNearbyTask struct {
	app.Task
	Longitude float64 `json:"longitude"` //经度
	Latitude  float64 `json:"latitude"`  //纬度
	Distance  float64 `json:"distance"`  //路面距离 km
	PageIndex int     `json:"p"`
	PageSize  int     `json:"size"`
	Counter   bool    `json:"counter"`
	Result    StopNearbyTaskResult
}

func (T *StopNearbyTask) GetResult() interface{} {
	return &T.Result
}

func (T *StopNearbyTask) GetInhertType() string {
	return "stop"
}

func (T *StopNearbyTask) GetClientName() string {
	return "Stop.Nearby"
}

func (S NearbyStopSlice) Len() int {
	return len(S)
}

func (S NearbyStopSlice) Less(i, j int) bool {
	return S[i].Distance < S[j].Distance
}

func (S NearbyStopSlice) Swap(i, j int) {
	v := S[i]
	S[i] = S[j]
	S[j] = v
}
