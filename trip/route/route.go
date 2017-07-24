package route

import (
	"database/sql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
)

const RouteStatusNone = 0    //可用
const RouteStatusTrash = 300 //已回收
const RouteDirectionNone = 0 //去程
const RouteDirectionBack = 1 //返程

/**
 * 路线
 */
type Route struct {
	Id            int64   `json:"id"`
	Start         string  `json:"start"`         //开始名
	End           string  `json:"end"`           //结束名
	Alias         string  `json:"alias"`         //别名
	Tags          string  `json:"tags"`          //搜索标签
	Distance      float64 `json:"distance"`      //路面距离 km
	Status        int     `json:"status"`        //状态
	StartCityId   int64   `json:"startCityId"`   //开始城市ID
	StartCityPath string  `json:"startCityPath"` //开始城市
	EndCityId     int64   `json:"endCityId"`     //结束城市ID
	EndCityPath   string  `json:"endCityPath"`   //结束城市
	WhiteList     string  `json:"whiteList"`     //白名单
	Ctime         int64   `json:"ctime"`
}

const RouteStopTypeNone = 0
const RouteStopTypeOn = 1    //上车站
const RouteStopTypeUnder = 2 //下车站

/**
 * 路线站点
 */
type RouteStop struct {
	Id        int64   `json:"id"`
	Title     string  `json:"title"`     //站点备注名
	StopId    int64   `json:"stopId"`    //站点ID
	RouteId   int64   `json:"routeId"`   //路线ID
	Longitude float64 `json:"longitude"` //经度
	Latitude  float64 `json:"latitude"`  //纬度
	Type      int     `json:"type"`      //站点类型
	Direction int     `json:"direction"` //方向
	Body      string  `json:"body"`      //站点描述
	Images    string  `json:"images"`    //站点图片
	Ctime     int64   `json:"ctime"`     //创建时间
}

type IRouteApp interface {
	app.IApp
	GetDB() (*sql.DB, error)
	GetPrefix() string
	GetRouteTable() *kk.DBTable
	GetRouteStopTable() *kk.DBTable
}
