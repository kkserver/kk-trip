package car

import (
	"database/sql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
)

type Car struct {
	Id int64 `json:"id"`

	LicenceCode string `json:"licenceCode"` //行驶证 唯一

	Brand    string `json:"brand"`    //品牌名
	Name     string `json:"name"`     //车辆明
	PlateNo  string `json:"plateNo"`  //车牌号 唯一
	Capacity string `json:"capacity"` //排量

	Latitude  float64 `json:"latitude"`  //经度
	Longitude float64 `json:"longitude"` //纬度
	Ip        string  `json:"ip"`        //IP地址

	Atime int64 `json:"atime"`
	Ctime int64 `json:"ctime"`
}

type ICarApp interface {
	app.IApp
	GetDB() (*sql.DB, error)
	GetPrefix() string
	GetCarTable() *kk.DBTable
}
