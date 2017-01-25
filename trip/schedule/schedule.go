package schedule

import (
	"database/sql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
)

const ScheduleStatusNone = 0   //未上线
const ScheduleStatusIn = 1     //已上线
const ScheduleStatusStart = 2  //已发车
const ScheduleStatusEnd = 3    //已收车
const ScheduleStatusFail = 300 //已故障

/**
 * 班次排期
 */
type Schedule struct {
	Id        int64 `json:"id"`
	LineId    int64 `json:"lineId"`    //车次ID
	Date      int64 `json:"date"`      //日期
	Status    int   `json:"status"`    //状态
	MaxCount  int   `json:"maxCount"`  //最大售票数
	UMaxCount int   `json:"umaxCount"` //用户最大购票数
	Count     int   `json:"count"`     //已售票数

	CarId     int64 `json:"carId"`     //车辆ID
	DriverId  int64 `json:"driverId"`  //司机ID
	StartTime int64 `json:"startTime"` //发车时间
	EndTime   int64 `json:"endTime"`   //收车时间

	InTime int64 `json:"inTime"` //上线时间

	Ctime int64 `json:"ctime"`
}

type IScheduleApp interface {
	app.IApp
	GetDB() (*sql.DB, error)
	GetPrefix() string
	GetScheduleTable() *kk.DBTable
	GetTicketTable() *kk.DBTable
}
