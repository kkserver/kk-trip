package stop

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lbs/lbs"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"time"
)

type StopService struct {
	app.Service

	Create *StopCreateTask
	Get    *StopTask
	Set    *StopSetTask
	Remove *StopRemoveTask
	Query  *StopQueryTask
	Nearby *StopNearbyTask
}

func (S *StopService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *StopService) HandleStopCreateTask(a IStopApp, task *StopCreateTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_STOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Stop{}

	v.Title = task.Title
	v.Longitude = task.Longitude
	v.Latitude = task.Latitude
	v.Ctime = time.Now().Unix()

	_, err = kk.DBInsert(db, a.GetStopTable(), a.GetPrefix(), &v)

	if err != nil {
		task.Result.Errno = ERROR_STOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.Stop = &v

	return nil
}

func (S *StopService) HandleStopSetTask(a IStopApp, task *StopSetTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_STOP_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_STOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Stop{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetStopTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_STOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_STOP
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

		_, err = kk.DBUpdate(db, a.GetStopTable(), a.GetPrefix(), &v)

		if err != nil {
			task.Result.Errno = ERROR_STOP
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Stop = &v

	} else {
		task.Result.Errno = ERROR_STOP_NOT_FOUND_STOP
		task.Result.Errmsg = "Not found stop"
	}

	return nil
}

func (S *StopService) HandleStopTask(a IStopApp, task *StopTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_STOP_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_STOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Stop{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetStopTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_STOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_STOP
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Stop = &v

	} else {

		task.Result.Errno = ERROR_STOP_NOT_FOUND_STOP
		task.Result.Errmsg = "Not found stop"
	}

	return nil
}

func (S *StopService) HandleStopRemoveTask(a IStopApp, task *StopRemoveTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_STOP_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_STOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Stop{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetStopTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_STOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_STOP
			task.Result.Errmsg = err.Error()
			return nil
		}

		{

			trigger := TriggerStopRemovingTask{}
			trigger.Stop = &v

			err = app.Handle(a, &trigger)

			if err != nil {
				task.Result.Errno = ERROR_STOP
				task.Result.Errmsg = err.Error()
				return nil
			}
		}

		_, err = kk.DBDelete(db, a.GetStopTable(), a.GetPrefix(), " WHERE id=?", task.Id)

		if err != nil {
			task.Result.Errno = ERROR_STOP
			task.Result.Errmsg = err.Error()
			return nil
		}

		{
			trigger := TriggerStopRemovedTask{}
			trigger.Stop = &v
			app.Handle(a, &trigger)
		}

		task.Result.Stop = &v

	} else {
		task.Result.Errno = ERROR_STOP_NOT_FOUND_STOP
		task.Result.Errmsg = "Not found stop"
	}

	return nil
}

func (S *StopService) HandleStopQueryTask(a IStopApp, task *StopQueryTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_STOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	var stops = []Stop{}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(" WHERE 1")

	if task.Id != 0 {
		sql.WriteString(" AND id=?")
		args = append(args, task.Id)
	} else {

		if task.Keyword != "" {
			q := "%" + task.Keyword + "%"
			sql.WriteString(" AND ( title LIKE ? OR id=? )")
			args = append(args, q, task.Keyword)
		}

		var pageIndex = task.PageIndex
		var pageSize = task.PageSize

		if pageIndex < 1 {
			pageIndex = 1
		}

		if pageSize < 1 {
			pageSize = 10
		}

		if task.Counter {
			var counter = StopQueryCounter{}
			counter.PageIndex = pageIndex
			counter.PageSize = pageSize
			counter.PageSize, err = kk.DBQueryCount(db, a.GetStopTable(), a.GetPrefix(), sql.String(), args...)
			if err != nil {
				task.Result.Errno = ERROR_STOP
				task.Result.Errmsg = err.Error()
				return nil
			}
			task.Result.Counter = &counter
		}

		if task.OrderBy == "asc" {
			sql.WriteString(" ORDER BY id ASC")
		} else {
			sql.WriteString(" ORDER BY id DESC")
		}

		sql.WriteString(fmt.Sprintf(" LIMIT %d,%d", (pageIndex-1)*pageSize, pageSize))

		var v = Stop{}
		var scanner = kk.NewDBScaner(&v)

		rows, err := kk.DBQuery(db, a.GetStopTable(), a.GetPrefix(), sql.String(), args...)

		if err != nil {
			task.Result.Errno = ERROR_STOP
			task.Result.Errmsg = err.Error()
			return nil
		}

		defer rows.Close()

		for rows.Next() {

			err = scanner.Scan(rows)

			if err != nil {
				task.Result.Errno = ERROR_STOP
				task.Result.Errmsg = err.Error()
				return nil
			}

			stops = append(stops, v)
		}
	}

	task.Result.Stops = stops

	return nil
}

func (S *StopService) HandleStopNearbyTask(a IStopApp, task *StopNearbyTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_STOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	var stops = NearbyStopSlice{}

	loc := lbs.LngLat{task.Longitude, task.Latitude}

	box := lbs.BoxFromCenter(loc, task.Distance)

	sql := bytes.NewBuffer(nil)

	args := []interface{}{}

	sql.WriteString("FROM (SELECT *,( 3956*2*ASIN(SQRT(POWER(SIN((?-latitude) * PI()/180/2),2)+COS(?*pi()/180)*COS(latitude*pi()/180)*POWER(SIN((?-longitude)*pi()/180/2),2)))) as distance")

	sql.WriteString(fmt.Sprintf(" FROM %s%s", a.GetPrefix(), a.GetStopTable().Name))

	args = append(args, loc.Latitude, loc.Latitude, loc.Longitude)

	sql.WriteString(" WHERE longitude >= ? AND longitude <= ? AND latitude >= ? AND latitude <= ?) as v WHERE v.distance<=?")

	args = append(args, box.Min.Longitude, box.Max.Longitude, box.Min.Latitude, box.Max.Latitude, task.Distance)

	var pageIndex = task.PageIndex
	var pageSize = task.PageSize

	if pageIndex < 1 {
		pageIndex = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	if task.Counter {
		var counter = StopNearbyCounter{}
		counter.PageIndex = pageIndex
		counter.PageSize = pageSize

		rs, err := db.Query("SELECT COUNT(*) "+sql.String(), args...)

		if err != nil {
			task.Result.Errno = ERROR_STOP
			task.Result.Errmsg = err.Error()
			return nil
		}

		total := 0

		if rs.Next() {
			rs.Scan(&total)
		}

		rs.Close()

		if total%pageSize == 0 {
			counter.PageCount = total / pageSize
		} else {
			counter.PageCount = total/pageSize + 1
		}

		task.Result.Counter = &counter
	}

	sql.WriteString(" ORDER BY distance ASC, id DESC")

	sql.WriteString(fmt.Sprintf(" LIMIT %d,%d", (pageIndex-1)*pageSize, pageSize))

	rows, err := db.Query("SELECT * "+sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_STOP
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	var v = NearbyStop{}
	var scanner = kk.NewDBScaner(&v)

	for rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_STOP
			task.Result.Errmsg = err.Error()
			return nil
		}

		stops = append(stops, v)
	}

	task.Result.Stops = stops

	return nil
}
