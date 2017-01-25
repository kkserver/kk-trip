package suggest

import (
	"database/sql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
)

type ISuggestApp interface {
	app.IApp
	GetDB() (*sql.DB, error)
	GetPrefix() string
	GetTicketTable() *kk.DBTable
	GetLineTable() *kk.DBTable
}
