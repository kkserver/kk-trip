package trip

import (
	"database/sql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/remote"
	"github.com/kkserver/kk-order/order"
	"github.com/kkserver/kk-trip/trip/car"
	"github.com/kkserver/kk-trip/trip/driver"
	"github.com/kkserver/kk-trip/trip/line"
	"github.com/kkserver/kk-trip/trip/route"
	"github.com/kkserver/kk-trip/trip/schedule"
	"github.com/kkserver/kk-trip/trip/stop"
	"github.com/kkserver/kk-trip/trip/suggest"
	"github.com/kkserver/kk-trip/trip/ticket"
)

type TripApp struct {
	app.App
	DB *app.DBConfig

	Runloop bool

	Route          *route.RouteService
	RouteStop      *route.RouteStopService
	RouteTable     kk.DBTable
	RouteStopTable kk.DBTable

	Stop      *stop.StopService
	StopTable kk.DBTable

	Line      *line.LineService
	LineTable kk.DBTable

	Schedule      *schedule.ScheduleService
	ScheduleTable kk.DBTable

	Order      *order.OrderService
	OrderTable kk.DBTable

	Ticket      *ticket.TicketService
	TicketTable kk.DBTable

	Driver      *driver.DriverService
	DriverTable kk.DBTable

	Car      *car.CarService
	CarTable kk.DBTable

	Suggest *suggest.SuggestService

	Remote *remote.Service

	runloop *kk.Dispatch
}

func (A *TripApp) GetDB() (*sql.DB, error) {
	return A.DB.Get(A)
}

func (A *TripApp) GetPrefix() string {
	return A.DB.Prefix
}

func (A *TripApp) GetRouteTable() *kk.DBTable {
	return &A.RouteTable
}

func (A *TripApp) GetRouteStopTable() *kk.DBTable {
	return &A.RouteStopTable
}

func (A *TripApp) GetStopTable() *kk.DBTable {
	return &A.StopTable
}

func (A *TripApp) GetScheduleTable() *kk.DBTable {
	return &A.ScheduleTable
}

func (A *TripApp) GetLineTable() *kk.DBTable {
	return &A.LineTable
}

func (A *TripApp) GetOrderTable() *kk.DBTable {
	return &A.OrderTable
}

func (A *TripApp) GetTicketTable() *kk.DBTable {
	return &A.TicketTable
}

func (A *TripApp) GetDriverTable() *kk.DBTable {
	return &A.DriverTable
}

func (A *TripApp) GetCarTable() *kk.DBTable {
	return &A.CarTable
}

func (A *TripApp) GetRunloop() *kk.Dispatch {
	if A.runloop == nil {
		A.runloop = kk.NewDispatch()
	}
	return A.runloop
}
