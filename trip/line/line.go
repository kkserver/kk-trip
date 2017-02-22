package line

import (
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-trip/trip/route"
	"strconv"
	"strings"
)

const LineStatusNone = 0    //未上线
const LineStatusIn = 1      //已上线
const LineStatusTrash = 300 //已回收

/**
 * 车次
 */
type Line struct {
	Id        int64  `json:"id,string"`
	RouteId   int64  `json:"routeId,string"` //路线ID
	Price     int64  `json:"price"`          //原价
	Alias     string `json:"alias"`          //别名
	Time      int64  `json:"time"`           //发车时间
	EndTime   int64  `json:"endTime"`        //收车时间
	Status    int    `json:"status"`         //状态
	Direction int    `json:"direction"`      //方向
	Times     string `json:"times"`          //站点时间

	CarId    int64 `json:"carId"`    //车辆ID
	DriverId int64 `json:"driverId"` //司机ID

	Ctime int64 `json:"ctime"`
}

type LineTime int64

func (T LineTime) String() string {
	return fmt.Sprintf("%02d:%02d", T/3600, (T/60)%60)
}

func LineTimeFromString(value string) LineTime {
	vs := strings.Split(value, ":")
	if len(vs) > 1 {
		h, _ := strconv.Atoi(vs[0])
		m, _ := strconv.Atoi(vs[1])
		return LineTime(h*3600 + m*60)
	}
	return 0
}

type LineTimeSlice []LineTime

func LineTimeSliceFromString(value string) LineTimeSlice {

	v := LineTimeSlice{}

	for _, vv := range strings.Split(value, " ") {
		v = append(v, LineTimeFromString(vv))
	}

	return v
}

type ILineApp interface {
	route.IRouteApp
	GetLineTable() *kk.DBTable
}
