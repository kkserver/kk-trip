package ticket

import (
	"github.com/kkserver/kk-lib/kk/app"
	"time"
)

type CalendarDay struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
	Week  int `json:"week"`
}

func (C *CalendarDay) SetTime(v time.Time) {
	C.Year = v.Year()
	C.Month = int(v.Month())
	C.Day = v.Day()
	C.Week = int(v.Weekday())
}

func (C CalendarDay) Compare(day CalendarDay) int {
	if C.Year < day.Year {
		return -1
	} else if C.Year > day.Year {
		return 1
	} else if C.Month < day.Month {
		return -1
	} else if C.Month > day.Month {
		return 1
	} else if C.Day < day.Day {
		return -1
	} else if C.Day > day.Day {
		return 1
	}
	return 0
}

func (C CalendarDay) IsZero() bool {
	return C.Year == 0 && C.Month == 0 && C.Day == 0 && C.Week == 0
}

type CalendarCell struct {
	Id        int64       `json:"id"`
	LineId    int64       `json:"lineId"`    //车次ID
	Date      int64       `json:"date"`      //日期
	Status    int         `json:"status"`    //状态
	MaxCount  int         `json:"maxCount"`  //最大售票数
	UMaxCount int         `json:"umaxCount"` //用户最大购票数
	Count     int         `json:"count"`     //已售票数
	Day       CalendarDay `json:"day"`
	Ucount    int         `json:"ucount"`
}

type CalendarRow []CalendarCell
type Calendar []CalendarRow

type TicketCalendarTaskResult struct {
	app.Result
	Calendar Calendar `json:"calendar,omitempty"`
}

type TicketCalendarTask struct {
	app.Task
	LineId int64 `json:"lineId"` //车次ID
	Uid    int64 `json:"uid"`    //用户ID
	Result TicketCalendarTaskResult
}

func (T *TicketCalendarTask) GetResult() interface{} {
	return &T.Result
}

func (T *TicketCalendarTask) GetInhertType() string {
	return "ticket"
}

func (T *TicketCalendarTask) GetClientName() string {
	return "Ticket.Calendar"
}
