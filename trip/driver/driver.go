package driver

import (
	"database/sql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
)

/**
 * 司机
 */
type Driver struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`        //姓名
	Code        string `json:"code"`        //身份证
	LicenceCode string `json:"licenceCode"` //驾驶证
	Phone       string `json:"phone"`       //手机号
	Ctime       int64  `json:"ctime"`
}

type IDriverApp interface {
	app.IApp
	GetDB() (*sql.DB, error)
	GetPrefix() string
	GetDriverTable() *kk.DBTable
}
