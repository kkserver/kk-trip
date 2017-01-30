package route

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lbs/lbs"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"strings"
	"time"
)

type RouteService struct {
	app.Service

	Create *RouteCreateTask
	Get    *RouteTask
	Set    *RouteSetTask
	Remove *RouteRemoveTask
	Query  *RouteQueryTask
	Count  *RouteCountTask
	Nearby *RouteNearbyTask
}

func (S *RouteService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *RouteService) HandleRouteCreateTask(a IRouteApp, task *RouteCreateTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Route{}

	v.Alias = task.Alias
	v.Start = task.Start
	v.End = task.End
	v.Distance = task.Distance
	v.Tags = task.Tags
	v.StartCityId = task.StartCityId
	v.EndCityId = task.EndCityId
	v.Ctime = time.Now().Unix()

	_, err = kk.DBInsert(db, a.GetRouteTable(), a.GetPrefix(), &v)

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.Route = &v

	return nil
}

func (S *RouteService) HandleRouteSetTask(a IRouteApp, task *RouteSetTask) error {

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

	var v = Route{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetRouteTable(), a.GetPrefix(), " WHERE id=?", task.Id)

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

		keys := map[string]bool{}

		if task.Start != nil {
			v.Start = dynamic.StringValue(task.Start, v.Start)
			keys["start"] = true
		}

		if task.End != nil {
			v.End = dynamic.StringValue(task.End, v.End)
			keys["end"] = true
		}

		if task.Alias != nil {
			v.Alias = dynamic.StringValue(task.Alias, v.Alias)
			keys["alias"] = true
		}

		if task.Tags != nil {
			v.Tags = dynamic.StringValue(task.Tags, v.Tags)
			keys["tags"] = true
		}

		if task.Distance != nil {
			v.Distance = dynamic.FloatValue(task.Distance, v.Distance)
			keys["distance"] = true
		}

		if task.Status != nil {
			v.Status = int(dynamic.IntValue(task.Status, int64(v.Status)))
			keys["status"] = true
		}

		if task.StartCityId != nil {
			v.StartCityId = dynamic.IntValue(task.StartCityId, int64(v.StartCityId))
			keys["startcityid"] = true
		}

		if task.EndCityId != nil {
			v.EndCityId = dynamic.IntValue(task.EndCityId, int64(v.EndCityId))
			keys["endcityid"] = true
		}

		_, err = kk.DBUpdateWithKeys(db, a.GetRouteTable(), a.GetPrefix(), &v, keys)

		if err != nil {
			task.Result.Errno = ERROR_ROUTE
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Route = &v

	} else {
		task.Result.Errno = ERROR_ROUTE_NOT_FOUND_ROUTE
		task.Result.Errmsg = "Not found route"
	}

	return nil
}

func (S *RouteService) HandleRouteTask(a IRouteApp, task *RouteTask) error {

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

	var v = Route{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetRouteTable(), a.GetPrefix(), " WHERE id=?", task.Id)

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

		task.Result.Route = &v

	} else {

		task.Result.Errno = ERROR_ROUTE_NOT_FOUND_ROUTE
		task.Result.Errmsg = "Not found route"
	}

	return nil
}

func (S *RouteService) HandleRouteRemoveTask(a IRouteApp, task *RouteRemoveTask) error {

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

	var v = Route{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetRouteTable(), a.GetPrefix(), " WHERE id=?", task.Id)

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

		{

			trigger := TriggerRouteRemovingTask{}
			trigger.Route = &v

			err = app.Handle(a, &trigger)

			if err != nil {
				task.Result.Errno = ERROR_ROUTE
				task.Result.Errmsg = err.Error()
				return nil
			}
		}

		_, err = kk.DBDelete(db, a.GetRouteTable(), a.GetPrefix(), " WHERE id=?", task.Id)

		if err != nil {
			task.Result.Errno = ERROR_ROUTE
			task.Result.Errmsg = err.Error()
			return nil
		}

		_, err = kk.DBDelete(db, a.GetRouteStopTable(), a.GetPrefix(), " WHERE routeid=?", task.Id)

		if err != nil {
			task.Result.Errno = ERROR_ROUTE
			task.Result.Errmsg = err.Error()
			return nil
		}

		{
			trigger := TriggerRouteRemovedTask{}
			trigger.Route = &v
			app.Handle(a, &trigger)
		}

		task.Result.Route = &v

	} else {
		task.Result.Errno = ERROR_ROUTE_NOT_FOUND_ROUTE
		task.Result.Errmsg = "Not found route"
	}

	return nil
}

func (S *RouteService) HandleRouteQueryTask(a IRouteApp, task *RouteQueryTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var routes = []Route{}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(" WHERE 1")

	if task.Id != 0 {
		sql.WriteString(" AND id=?")
		args = append(args, task.Id)
	} else {

		if task.Keyword != "" {
			q := "%" + task.Keyword + "%"
			sql.WriteString(" AND ( tags LIKE ? OR id=? OR start LIKE ? OR end LIKE ? OR alias LIKE ? OR CONCAT(start,'-',end,alias) LIKE ?)")
			args = append(args, q, task.Keyword, q, q, q, q)
		}

		if task.Status != "" {
			vs := strings.Split(task.Status, ",")
			sql.WriteString(" AND status IN (")
			for i, v := range vs {
				if i != 0 {
					sql.WriteString(",")
				}
				sql.WriteString("?")
				args = append(args, v)
			}
			sql.WriteString(")")
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
			var counter = RouteQueryCounter{}
			counter.PageIndex = pageIndex
			counter.PageSize = pageSize
			total, err := kk.DBQueryCount(db, a.GetRouteTable(), a.GetPrefix(), sql.String(), args...)
			if err != nil {
				task.Result.Errno = ERROR_ROUTE
				task.Result.Errmsg = err.Error()
				return nil
			}
			if total%pageSize == 0 {
				counter.PageCount = total / pageSize
			} else {
				counter.PageCount = total/pageSize + 1
			}
			task.Result.Counter = &counter
		}

		if task.OrderBy == "asc" {
			sql.WriteString(" ORDER BY id ASC")
		} else {
			sql.WriteString(" ORDER BY id DESC")
		}

		sql.WriteString(fmt.Sprintf(" LIMIT %d,%d", (pageIndex-1)*pageSize, pageSize))

	}

	fmt.Println("SQL", sql.String(), task)

	var v = Route{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetRouteTable(), a.GetPrefix(), sql.String(), args...)

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

		routes = append(routes, v)
	}

	task.Result.Routes = routes

	return nil
}

func (S *RouteService) HandleRouteCountTask(a IRouteApp, task *RouteCountTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(" WHERE 1")

	if task.Id != 0 {
		sql.WriteString(" AND id=?")
		args = append(args, task.Id)
	} else {

		if task.Keyword != "" {
			q := "%" + task.Keyword + "%"
			sql.WriteString(" AND ( tags LIKE ? OR id=? OR start LIKE ? OR end LIKE ? OR alias LIKE ? OR CONCAT(start,'-',end,alias) LIKE ?)")
			args = append(args, q, task.Keyword, q, q, q, q)
		}

		if task.Status != "" {
			vs := strings.Split(task.Status, ",")
			sql.WriteString(" AND status IN (")
			for i, v := range vs {
				if i != 0 {
					sql.WriteString(",")
				}
				sql.WriteString("?")
				args = append(args, v)
			}
			sql.WriteString(")")
		}

	}

	v, err := kk.DBQueryCount(db, a.GetRouteTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.Count = v

	return nil
}

func (S *RouteService) HandleRouteNearbyTask(a IRouteApp, task *RouteNearbyTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var routes = []Route{}

	loc := lbs.LngLat{task.Longitude, task.Latitude}

	box := lbs.BoxFromCenter(loc, task.Distance)

	sql := bytes.NewBuffer(nil)

	args := []interface{}{}

	sql.WriteString("FROM (SELECT *,( 3956*2*ASIN(SQRT(POWER(SIN((?-latitude) * PI()/180/2),2)+COS(?*pi()/180)*COS(latitude*pi()/180)*POWER(SIN((?-longitude)*pi()/180/2),2)))) as distance")

	sql.WriteString(fmt.Sprintf(" FROM %s%s", a.GetPrefix(), a.GetRouteStopTable().Name))

	args = append(args, loc.Latitude, loc.Latitude, loc.Longitude)

	sql.WriteString(fmt.Sprintf(" WHERE longitude >= ? AND longitude <= ? AND latitude >= ? AND latitude <= ?) as v LEFT JOIN %s%s as r ON v.routeid=r.id WHERE v.distance<=? GROUP BY r.id HAVING 1", a.GetPrefix(), a.GetRouteTable().Name))

	args = append(args, box.Min.Longitude, box.Max.Longitude, box.Min.Latitude, box.Max.Latitude, task.Distance)

	if task.Keyword != "" {
		q := "%" + task.Keyword + "%"
		sql.WriteString(" AND ( r.tags LIKE ? OR r.id=? OR r.start LIKE ? OR r.end LIKE ? OR r.alias LIKE ? OR CONCAT(r.start,'-',r.end,r.alias) LIKE ?)")
		args = append(args, q, task.Keyword, q, q, q, q)
	}

	if task.Status != "" {
		vs := strings.Split(task.Status, ",")
		sql.WriteString(" AND r.status IN (")
		for i, v := range vs {
			if i != 0 {
				sql.WriteString(",")
			}
			sql.WriteString("?")
			args = append(args, v)
		}
		sql.WriteString(")")
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
		var counter = RouteNearbyCounter{}
		counter.PageIndex = pageIndex
		counter.PageSize = pageSize

		rs, err := db.Query("SELECT COUNT(*) "+sql.String(), args...)

		if err != nil {
			task.Result.Errno = ERROR_ROUTE
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

	sql.WriteString(" ORDER BY distance ASC")

	sql.WriteString(fmt.Sprintf(" LIMIT %d,%d", (pageIndex-1)*pageSize, pageSize))

	rows, err := db.Query("SELECT r.id as id , r.start as start,r.end as end,r.alias as alias,r.tags as tags,r.status as status,r.ctime as ctime, MIN(v.distance) as distance "+sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_ROUTE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	var v = Route{}
	var scanner = kk.NewDBScaner(&v)

	for rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_ROUTE
			task.Result.Errmsg = err.Error()
			return nil
		}

		routes = append(routes, v)
	}

	task.Result.Routes = routes

	return nil
}
