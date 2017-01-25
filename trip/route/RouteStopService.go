package route

import (
	"bytes"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"github.com/kkserver/kk-trip/trip/stop"
	"time"
)

const StopNearbyDistance = float64(0.01)

type RouteStopService struct {
	app.Service
	Create   *RouteStopCreateTask
	Set      *RouteStopSetTask
	Remove   *RouteStopRemoveTask
	Query    *RouteStopQueryTask
	Exchange *RouteStopExchangeTask
}

func (S *RouteStopService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *RouteStopService) HandleRouteStopCreateTask(a IRouteApp, task *RouteStopCreateTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := RouteStop{}

	v.Title = task.Title
	v.Direction = task.Direction
	v.Type = task.Type
	v.Latitude = task.Latitude
	v.Longitude = task.Longitude
	v.RouteId = task.RouteId
	v.Ctime = time.Now().Unix()

	{
		t := stop.StopNearbyTask{}
		t.Longitude = task.Longitude
		t.Latitude = task.Latitude
		t.Distance = StopNearbyDistance
		t.PageIndex = 1
		t.PageSize = 1
		app.Handle(a, &t)
		if t.Result.Stops != nil && len(t.Result.Stops) > 0 {
			v.StopId = t.Result.Stops[0].Id
		} else {
			tt := stop.StopCreateTask{}
			tt.Longitude = task.Longitude
			tt.Latitude = task.Latitude
			tt.Title = task.Title
			app.Handle(a, &tt)
			if tt.Result.Stop != nil {
				v.StopId = tt.Result.Stop.Id
			}
		}
	}

	_, err = kk.DBInsert(db, a.GetRouteStopTable(), a.GetPrefix(), &v)

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.Stop = &v

	return nil
}

func (S *RouteStopService) HandleRouteStopSetTask(a IRouteApp, task *RouteStopSetTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_ROUTE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = RouteStop{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetRouteStopTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_ROUTE
			task.Result.Errmsg = err.Error()
			return nil
		}

		if task.Title != nil {
			v.Title = dynamic.StringValue(task.Title, v.Title)
		}

		if task.Longitude != nil {
			v.Longitude = dynamic.FloatValue(task.Longitude, v.Longitude)
		}

		if task.Latitude != nil {
			v.Latitude = dynamic.FloatValue(task.Latitude, v.Latitude)
		}

		if task.Latitude != nil || task.Longitude != nil {

			{
				t := stop.StopNearbyTask{}
				t.Longitude = v.Longitude
				t.Latitude = v.Latitude
				t.Distance = StopNearbyDistance
				t.PageIndex = 1
				t.PageSize = 1
				app.Handle(a, &t)
				if t.Result.Stops != nil && len(t.Result.Stops) > 0 {
					v.StopId = t.Result.Stops[0].Id
				} else {
					tt := stop.StopCreateTask{}
					tt.Longitude = v.Longitude
					tt.Latitude = v.Latitude
					tt.Title = v.Title
					app.Handle(a, &tt)
					if tt.Result.Stop != nil {
						v.StopId = tt.Result.Stop.Id
					}
				}
			}

		}

		_, err = kk.DBUpdate(db, a.GetRouteStopTable(), a.GetPrefix(), &v)

		if err != nil {
			task.Result.Errno = ERROR_ROUTE
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Stop = &v

	} else {
		task.Result.Errno = ERROR_ROUTE_NOT_FOUND_ROUTE
		task.Result.Errmsg = "Not found route"
	}

	return nil
}

func (S *RouteStopService) HandleRouteStopRemoveTask(a IRouteApp, task *RouteStopRemoveTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_ROUTE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = RouteStop{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetRouteStopTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_ROUTE
			task.Result.Errmsg = err.Error()
			return nil
		}

		_, err = kk.DBDelete(db, a.GetRouteStopTable(), a.GetPrefix(), " WHERE id=?", v.Id)

		if err != nil {
			task.Result.Errno = ERROR_ROUTE
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Stop = &v

	} else {

		task.Result.Errno = ERROR_ROUTE_NOT_FOUND_STOP
		task.Result.Errmsg = "Not found route stop"
	}

	return nil
}

func (S *RouteStopService) HandleTriggerStopRemovingTask(a IRouteApp, task *stop.TriggerStopRemovingTask) error {

	var db, err = a.GetDB()

	if err != nil {
		return err
	}

	rows, err := kk.DBQuery(db, a.GetRouteStopTable(), a.GetPrefix(), " WHERE stopid=?", task.Stop.Id)

	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		return app.NewError(ERROR_ROUTE, "The stop in use")
	}

	return nil
}

func (S *RouteStopService) HandleRouteStopQueryTask(a IRouteApp, task *RouteStopQueryTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var stops = []RouteStop{}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(" WHERE routeid=?")

	args = append(args, task.RouteId)

	if task.Type != nil {
		sql.WriteString(" AND type=?")
		args = append(args, task.Type)
	}

	if task.Direction != nil {
		sql.WriteString(" AND direction=?")
		args = append(args, task.Direction)
	}

	sql.WriteString(" ORDER BY direction ASC,type ASC,id ASC")

	var v = RouteStop{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetRouteStopTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	for rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_ROUTE
			task.Result.Errmsg = err.Error()
			return nil
		}

		stops = append(stops, v)
	}

	task.Result.Stops = stops

	return nil
}

func (S *RouteStopService) HandleRouteStopExchangeTask(a IRouteApp, task *RouteStopExchangeTask) error {

	if task.FromId == 0 {
		task.Result.Errno = ERROR_ROUTE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found fromId"
		return nil
	}

	if task.ToId == 0 {
		task.Result.Errno = ERROR_ROUTE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found toId"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var fstop = RouteStop{}
	var tstop = RouteStop{}

	rows, err := kk.DBQuery(db, a.GetRouteStopTable(), a.GetPrefix(), " WHERE id IN (?,?)", task.FromId, task.ToId)

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		var scanner = kk.NewDBScaner(&fstop)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_ROUTE
			task.Result.Errmsg = err.Error()
			return nil
		}

	} else {

		task.Result.Errno = ERROR_ROUTE_NOT_FOUND_STOP
		task.Result.Errmsg = "Not found route from stop"
		return nil
	}

	if rows.Next() {

		var scanner = kk.NewDBScaner(&tstop)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_ROUTE
			task.Result.Errmsg = err.Error()
			return nil
		}

	} else {
		task.Result.Errno = ERROR_ROUTE_NOT_FOUND_STOP
		task.Result.Errmsg = "Not found route to stop"
		return nil
	}

	id := fstop.Id
	fstop.Id = tstop.Id
	tstop.Id = id

	_, err = kk.DBUpdate(db, a.GetRouteStopTable(), a.GetPrefix(), &fstop)

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	_, err = kk.DBUpdate(db, a.GetRouteStopTable(), a.GetPrefix(), &tstop)

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	return nil
}
