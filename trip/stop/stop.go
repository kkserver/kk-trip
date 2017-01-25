package stop

import (
	"database/sql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
)

/**
 * 站点
 */
type Stop struct {
	Id        int64   `json:"id"`
	Title     string  `json:"title"`     //站点名
	Longitude float64 `json:"longitude"` //经度
	Latitude  float64 `json:"latitude"`  //纬度
	Ctime     int64   `json:"ctime"`     //创建时间
}

type IStopApp interface {
	app.IApp
	GetDB() (*sql.DB, error)
	GetPrefix() string
	GetStopTable() *kk.DBTable
}
