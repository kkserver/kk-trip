package line

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"github.com/kkserver/kk-trip/trip/route"
	"strings"
	"time"
)

type LineService struct {
	app.Service

	Create *LineCreateTask
	Get    *LineTask
	Set    *LineSetTask
	Remove *LineRemoveTask
	Query  *LineQueryTask
	Count  *LineCountTask
}

func (S *LineService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *LineService) HandleLineCreateTask(a ILineApp, task *LineCreateTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_LINE
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Line{}

	v.RouteId = task.RouteId
	v.Alias = task.Alias
	v.Direction = task.Direction
	v.Price = task.Price
	v.Times = task.Times

	ts := LineTimeSliceFromString(task.Times)

	if len(ts) > 0 {
		v.Time = int64(ts[0])
	}

	v.Ctime = time.Now().Unix()

	_, err = kk.DBInsert(db, a.GetLineTable(), a.GetPrefix(), &v)

	if err != nil {
		task.Result.Errno = ERROR_LINE
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.Line = &v

	return nil
}

func (S *LineService) HandleLineSetTask(a ILineApp, task *LineSetTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_LINE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_LINE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Line{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetLineTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_LINE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_LINE
			task.Result.Errmsg = err.Error()
			return nil
		}

		if task.Price != nil {
			v.Price = dynamic.IntValue(task.Price, v.Price)
		}

		if task.Alias != nil {
			v.Alias = dynamic.StringValue(task.Alias, v.Alias)
		}

		if task.Direction != nil {
			v.Direction = int(dynamic.IntValue(task.Direction, int64(v.Direction)))
		}

		if task.Times != nil {
			v.Times = dynamic.StringValue(task.Times, v.Times)
			ts := LineTimeSliceFromString(v.Times)
			if len(ts) > 0 {
				v.Time = int64(ts[0])
			}
		}

		if task.Status != nil {
			v.Status = int(dynamic.IntValue(task.Status, int64(v.Status)))
		}

		_, err = kk.DBUpdate(db, a.GetLineTable(), a.GetPrefix(), &v)

		if err != nil {
			task.Result.Errno = ERROR_LINE
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Line = &v

	} else {
		task.Result.Errno = ERROR_LINE_NOT_FOUND_LINE
		task.Result.Errmsg = "Not found line"
	}

	return nil
}

func (S *LineService) HandleLineTask(a ILineApp, task *LineTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_LINE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_LINE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Line{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetLineTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_LINE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_LINE
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Line = &v

	} else {

		task.Result.Errno = ERROR_LINE_NOT_FOUND_LINE
		task.Result.Errmsg = "Not found line"
	}

	return nil
}

func (S *LineService) HandleLineRemoveTask(a ILineApp, task *LineRemoveTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_LINE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_LINE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Line{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetLineTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_LINE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_LINE
			task.Result.Errmsg = err.Error()
			return nil
		}

		{

			trigger := TriggerLineRemovingTask{}
			trigger.Line = &v

			err = app.Handle(a, &trigger)

			if err != nil {
				task.Result.Errno = ERROR_LINE
				task.Result.Errmsg = err.Error()
				return nil
			}
		}

		_, err = kk.DBDelete(db, a.GetLineTable(), a.GetPrefix(), " WHERE id=?", task.Id)

		if err != nil {
			task.Result.Errno = ERROR_LINE
			task.Result.Errmsg = err.Error()
			return nil
		}

		{
			trigger := TriggerLineRemovedTask{}
			trigger.Line = &v
			app.Handle(a, &trigger)
		}

		task.Result.Line = &v

	} else {
		task.Result.Errno = ERROR_LINE_NOT_FOUND_LINE
		task.Result.Errmsg = "Not found line"
	}

	return nil
}

func (S *LineService) HandleTriggerRouteRemovingTask(a ILineApp, task *route.TriggerRouteRemovingTask) error {

	var db, err = a.GetDB()

	if err != nil {
		return err
	}

	rows, err := kk.DBQuery(db, a.GetLineTable(), a.GetPrefix(), " WHERE routeid=?", task.Route.Id)

	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		return app.NewError(ERROR_LINE, "The route in use")
	}

	return nil
}

func (S *LineService) HandleLineQueryTask(a ILineApp, task *LineQueryTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_LINE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var lines = []Line{}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(fmt.Sprintf(" FROM %s%s as v LEFT JOIN %s%s as r ON v.routeid=r.id WHERE 1",
		a.GetPrefix(), a.GetLineTable().Name, a.GetPrefix(), a.GetRouteTable().Name))

	if task.Id != 0 {
		sql.WriteString(" AND v.id=?")
		args = append(args, task.Id)
	} else {

		if task.Keyword != "" {
			q := "%" + task.Keyword + "%"
			sql.WriteString(" AND ( v.alias LIKE ? OR v.id=? OR r.start LIKE ? OR r.end LIKE ? OR r.alias LIKE ? OR CONCAT(r.start,'-',r.end,r.alias) LIKE ? OR CONCAT(r.start,'-',r.end,v.alias) LIKE ?)")
			args = append(args, q, task.Keyword, q, q, q, q, q)
		}

		if task.Status != "" {
			vs := strings.Split(task.Status, ",")
			sql.WriteString(" AND v.status IN (")
			for i, v := range vs {
				if i != 0 {
					sql.WriteString(",")
				}
				sql.WriteString("?")
				args = append(args, v)
			}
			sql.WriteString(")")
		}

		if task.Direction != "" {
			vs := strings.Split(task.Direction, ",")
			sql.WriteString(" AND v.direction IN (")
			for i, v := range vs {
				if i != 0 {
					sql.WriteString(",")
				}
				sql.WriteString("?")
				args = append(args, v)
			}
			sql.WriteString(")")
		}

		if task.RouteId != 0 {
			sql.WriteString(" AND v.routeid=?")
			args = append(args, task.RouteId)
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

			var counter = LineQueryCounter{}

			counter.PageIndex = pageIndex
			counter.PageSize = pageSize

			rows, err := db.Query("SELECT COUNT(*)"+sql.String(), args...)

			if err != nil {
				task.Result.Errno = ERROR_LINE
				task.Result.Errmsg = err.Error()
				return nil
			}

			rowCount := 0

			if rows.Next() {
				err = rows.Scan(&rowCount)
				if err != nil {
					task.Result.Errno = ERROR_LINE
					task.Result.Errmsg = err.Error()
					return nil
				}
			}

			rows.Close()

			if rowCount%pageSize == 0 {
				counter.PageCount = rowCount / pageSize
			} else {
				counter.PageCount = rowCount/pageSize + 1
			}

			counter.RowCount = rowCount

			task.Result.Counter = &counter

		}

		if task.OrderBy == "asc" {
			sql.WriteString(" ORDER BY v.id ASC")
		} else if task.OrderBy == "time" {
			sql.WriteString(" ORDER BY v.time ASC")
		} else {
			sql.WriteString(" ORDER BY v.id DESC")
		}

		sql.WriteString(fmt.Sprintf(" LIMIT %d,%d", (pageIndex-1)*pageSize, pageSize))

	}

	var v = Line{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := db.Query("SELECT v.* "+sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_LINE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	for rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_LINE
			task.Result.Errmsg = err.Error()
			return nil
		}

		lines = append(lines, v)
	}

	task.Result.Lines = lines

	return nil
}

func (S *LineService) HandleLineCountTask(a ILineApp, task *LineCountTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_LINE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(fmt.Sprintf(" FROM %s%s as v LEFT JOIN %s%s as r ON v.routeid=r.id WHERE 1",
		a.GetPrefix(), a.GetLineTable().Name, a.GetPrefix(), a.GetRouteTable().Name))

	if task.Id != 0 {
		sql.WriteString(" AND v.id=?")
		args = append(args, task.Id)
	} else {

		if task.Keyword != "" {
			q := "%" + task.Keyword + "%"
			sql.WriteString(" AND ( v.alias LIKE ? OR v.id=? OR r.start LIKE ? OR r.end LIKE ? OR r.alias LIKE ? OR CONCAT(r.start,'-',r.end,r.alias) LIKE ? OR CONCAT(r.start,'-',r.end,v.alias) LIKE ?)")
			args = append(args, q, task.Keyword, q, q, q, q, q)
		}

		if task.Status != "" {
			vs := strings.Split(task.Status, ",")
			sql.WriteString(" AND v.status IN (")
			for i, v := range vs {
				if i != 0 {
					sql.WriteString(",")
				}
				sql.WriteString("?")
				args = append(args, v)
			}
			sql.WriteString(")")
		}

		if task.Direction != "" {
			vs := strings.Split(task.Direction, ",")
			sql.WriteString(" AND v.direction IN (")
			for i, v := range vs {
				if i != 0 {
					sql.WriteString(",")
				}
				sql.WriteString("?")
				args = append(args, v)
			}
			sql.WriteString(")")
		}

		if task.RouteId != 0 {
			sql.WriteString(" AND v.routeid=?")
			args = append(args, task.RouteId)
		}

	}

	rows, err := db.Query("SELECT COUNT(*)"+sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_LINE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	v := 0

	if rows.Next() {
		err = rows.Scan(&v)
		if err != nil {
			task.Result.Errno = ERROR_LINE
			task.Result.Errmsg = err.Error()
			return nil
		}
	}

	task.Result.Count = v

	return nil
}
