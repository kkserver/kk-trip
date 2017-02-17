package schedule

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"log"
	"strconv"
	"strings"
	"time"
)

type ScheduleService struct {
	app.Service

	Create   *ScheduleCreateTask
	Get      *ScheduleTask
	Set      *ScheduleSetTask
	Remove   *ScheduleRemoveTask
	Query    *ScheduleQueryTask
	In       *ScheduleInTask
	Off      *ScheduleOffTask
	Start    *ScheduleStartTask
	End      *ScheduleEndTask
	Fail     *ScheduleFailTask
	BatchSet *ScheduleBatchSetTask
	Count    *ScheduleCountTask
}

func (S *ScheduleService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *ScheduleService) HandleRunloopTask(a IScheduleApp, task *app.RunloopTask) error {

	var db, err = a.GetDB()

	if err != nil {
		return err
	}

	var fn func() = nil

	fn = func() {

		now := time.Now()
		now = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

		log.Println("ScheduleService", "Runloop", "SQL", fmt.Sprintf("UPDATE %s%s SET status=? WHERE status=? AND intime !=0 AND intime<=?", a.GetPrefix(), a.GetScheduleTable().Name))

		_, err := db.Exec(fmt.Sprintf("UPDATE %s%s SET status=? WHERE status=? AND intime !=0 AND intime<=?", a.GetPrefix(), a.GetScheduleTable().Name), ScheduleStatusIn, ScheduleStatusNone, now.Unix())

		if err != nil {
			log.Println("ScheduleService", "Runloop", "Fail", err.Error())
		}

		log.Println("ScheduleService", "Runloop", "SQL", fmt.Sprintf("UPDATE %s%s as s INNER JOIN %s%s as l ON s.lineid=l.id SET s.status=? WHERE ? > s.date + l.time AND s.status=?", a.GetPrefix(), a.GetScheduleTable().Name, a.GetPrefix(), a.GetLineTable().Name))

		_, err = db.Exec(fmt.Sprintf("UPDATE %s%s as s INNER JOIN %s%s as l ON s.lineid=l.id SET s.status=? WHERE ? > s.date + l.time AND s.status=?", a.GetPrefix(), a.GetScheduleTable().Name, a.GetPrefix(), a.GetLineTable().Name), ScheduleStatusStart, time.Now().Unix(), ScheduleStatusIn)

		if err != nil {
			log.Println("ScheduleService", "Runloop", "Fail", err.Error())
		}

		log.Println("ScheduleService", "Runloop", "OK")

		a.GetRunloop().AsyncDelay(fn, 10*time.Second)

	}

	fn()

	return nil
}

func (S *ScheduleService) HandleScheduleCreateTask(a IScheduleApp, task *ScheduleCreateTask) error {

	if task.LineId == 0 {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND_LINE_ID
		task.Result.Errmsg = "Not found lineId"
		return nil
	}

	if task.Date == 0 {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND_DATE
		task.Result.Errmsg = "Not found date"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	rows, err := kk.DBQuery(db, a.GetScheduleTable(), a.GetPrefix(), " WHERE lineid=? AND date=?", task.LineId, task.Date)

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {
		task.Result.Errno = ERROR_SCHEDULE_EXISTS
		task.Result.Errmsg = "The date already exists"
		return nil
	}

	v := Schedule{}

	v.LineId = task.LineId
	v.Date = task.Date
	v.MaxCount = task.MaxCount
	v.UMaxCount = task.UMaxCount
	v.CarId = task.CarId
	v.DriverId = task.DriverId
	v.InTime = task.InTime
	v.Ctime = time.Now().Unix()

	_, err = kk.DBInsert(db, a.GetScheduleTable(), a.GetPrefix(), &v)

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.Schedule = &v

	if v.InTime != 0 && v.Ctime >= v.InTime {
		t := ScheduleInTask{}
		t.Id = v.Id
		app.Handle(a, &t)
	}

	return nil
}

func (S *ScheduleService) HandleScheduleBatchSetTask(a IScheduleApp, task *ScheduleBatchSetTask) error {

	if task.LineIds == "" {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND_LINE_ID
		task.Result.Errmsg = "Not found lineIds"
		return nil
	}

	if task.Dates == "" {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND_DATE
		task.Result.Errmsg = "Not found dates"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	tx, err := db.Begin()

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	now := time.Now()

	err = func() error {

		for _, slineId := range strings.Split(task.LineIds, ",") {

			for _, sdate := range strings.Split(task.Dates, ",") {

				lineId, err := strconv.ParseInt(slineId, 10, 64)

				if err != nil {
					return err
				}

				date, err := time.Parse("2006-01-02", sdate)

				if err != nil {
					return err
				}

				v := Schedule{}

				rows, err := kk.DBQuery(db, a.GetScheduleTable(), a.GetPrefix(), " WHERE lineid=? AND date=?", lineId, date.Unix())

				if err != nil {
					return err
				}

				if rows.Next() {

					scanner := kk.NewDBScaner(&v)
					err = scanner.Scan(rows)
					if err != nil {
						return err
					}

					v.InTime = date.AddDate(0, 0, -task.AdvanceDays).Unix()

					if v.Status == ScheduleStatusNone {
						if now.Unix() >= v.InTime {
							v.Status = ScheduleStatusIn
						}
					}

					v.MaxCount = task.MaxCount
					v.UMaxCount = task.UMaxCount

					_, err = kk.DBUpdate(tx, a.GetScheduleTable(), a.GetPrefix(), &v)

					if err != nil {
						return err
					}

				} else {

					v.LineId = lineId
					v.Date = date.Unix()
					v.MaxCount = task.MaxCount
					v.UMaxCount = task.UMaxCount
					v.InTime = date.AddDate(0, 0, -task.AdvanceDays).Unix()
					v.Ctime = time.Now().Unix()

					if now.Unix() >= v.InTime {
						v.Status = ScheduleStatusIn
					}

					_, err = kk.DBInsert(tx, a.GetScheduleTable(), a.GetPrefix(), &v)

					if err != nil {
						return err
					}
				}

				rows.Close()
			}

		}

		return nil
	}()

	if err != nil {

		e, ok := err.(*app.Error)

		if ok {
			task.Result.Errno = e.Errno
			task.Result.Errmsg = e.Errmsg
		} else {
			task.Result.Errno = ERROR_SCHEDULE
			task.Result.Errmsg = err.Error()
		}

		tx.Rollback()

		return nil
	} else {

		err = tx.Commit()

		if err != nil {
			tx.Rollback()
			task.Result.Errno = ERROR_SCHEDULE
			task.Result.Errmsg = err.Error()
			return nil
		}
	}

	return nil
}

func (S *ScheduleService) HandleScheduleSetTask(a IScheduleApp, task *ScheduleSetTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Schedule{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetScheduleTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_SCHEDULE
			task.Result.Errmsg = err.Error()
			return nil
		}

		if task.MaxCount != nil {
			v.MaxCount = int(dynamic.IntValue(task.MaxCount, int64(v.MaxCount)))
		}

		if task.UMaxCount != nil {
			v.UMaxCount = int(dynamic.IntValue(task.UMaxCount, int64(v.UMaxCount)))
		}

		if task.CarId != nil {
			v.CarId = dynamic.IntValue(task.CarId, v.CarId)
		}

		if task.DriverId != nil {
			v.DriverId = dynamic.IntValue(task.DriverId, v.DriverId)
		}

		if task.InTime != nil {
			v.InTime = dynamic.IntValue(task.InTime, v.InTime)
		}

		_, err = kk.DBUpdate(db, a.GetScheduleTable(), a.GetPrefix(), &v)

		if err != nil {
			task.Result.Errno = ERROR_SCHEDULE
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Schedule = &v

	} else {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND
		task.Result.Errmsg = "Not found schedule"
	}

	return nil
}

func (S *ScheduleService) HandleScheduleTask(a IScheduleApp, task *ScheduleTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Schedule{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetScheduleTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_SCHEDULE
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Schedule = &v

	} else {

		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND
		task.Result.Errmsg = "Not found schedule"
	}

	return nil
}

func (S *ScheduleService) HandleScheduleRemoveTask(a IScheduleApp, task *ScheduleRemoveTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Schedule{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetScheduleTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_SCHEDULE
			task.Result.Errmsg = err.Error()
			return nil
		}

		if v.Count > 0 {
			task.Result.Errno = ERROR_SCHEDULE_HAS_PAY
			task.Result.Errmsg = "The schedule has pay"
			return nil
		}

		_, err = kk.DBDelete(db, a.GetScheduleTable(), a.GetPrefix(), " WHERE id=?", task.Id)

		if err != nil {
			task.Result.Errno = ERROR_SCHEDULE
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Schedule = &v

	} else {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND
		task.Result.Errmsg = "Not found route"
	}

	return nil
}

func (S *ScheduleService) HandleScheduleQueryTask(a IScheduleApp, task *ScheduleQueryTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var schedules = []Schedule{}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(" WHERE 1")

	if task.Id != 0 {
		sql.WriteString(" AND id=?")
		args = append(args, task.Id)
	} else {

		if task.LineId != 0 {
			sql.WriteString(" AND lineid=?")
			args = append(args, task.LineId)
		}

		if task.StartDate != nil {
			sql.WriteString(" AND date>=?")
			args = append(args, task.StartDate)
		}

		if task.EndDate != nil {
			sql.WriteString(" AND date<?")
			args = append(args, task.EndDate)
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

		if task.CarId != nil {
			sql.WriteString(" AND carid=?")
			args = append(args, task.CarId)
		}

		if task.DriverId != nil {
			sql.WriteString(" AND dirverid=?")
			args = append(args, task.DriverId)
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
			var counter = ScheduleQueryCounter{}
			counter.PageIndex = pageIndex
			counter.PageSize = pageSize
			total, err := kk.DBQueryCount(db, a.GetScheduleTable(), a.GetPrefix(), sql.String(), args...)
			if err != nil {
				task.Result.Errno = ERROR_SCHEDULE
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
		} else if task.OrderBy == "count" {
			sql.WriteString(" ORDER BY count DESC,id DESC")
		} else {
			sql.WriteString(" ORDER BY id DESC")
		}

		sql.WriteString(fmt.Sprintf(" LIMIT %d,%d", (pageIndex-1)*pageSize, pageSize))

	}

	fmt.Println("SQL", sql.String())

	var v = Schedule{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetScheduleTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	for rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_SCHEDULE
			task.Result.Errmsg = err.Error()
			return nil
		}

		schedules = append(schedules, v)
	}

	task.Result.Schedules = schedules

	return nil
}

func (S *ScheduleService) HandleScheduleCountTask(a IScheduleApp, task *ScheduleCountTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var dates = []ScheduleCountDate{}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(fmt.Sprintf("SELECT v.date,COUNT(v.lineId) as lineCount FROM %s%s as v WHERE 1", a.GetPrefix(), a.GetScheduleTable().Name))

	if task.StartDate != nil {
		sql.WriteString(" AND v.date>=?")
		args = append(args, task.StartDate)
	}

	if task.EndDate != nil {
		sql.WriteString(" AND v.date<?")
		args = append(args, task.EndDate)
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

	if task.CarId != nil {
		sql.WriteString(" AND v.carid=?")
		args = append(args, task.CarId)
	}

	if task.DriverId != nil {
		sql.WriteString(" AND v.dirverid=?")
		args = append(args, task.DriverId)
	}

	sql.WriteString(" GROUP BY v.date ORDER BY v.date ASC")

	rows, err := db.Query(sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	v := ScheduleCountDate{}

	for rows.Next() {

		err = rows.Scan(&v.Date, &v.LineCount)

		if err != nil {
			task.Result.Errno = ERROR_SCHEDULE
			task.Result.Errmsg = err.Error()
			return nil
		}

		dates = append(dates, v)
	}

	task.Result.Dates = dates

	return nil
}

func (S *ScheduleService) HandleScheduleInTask(a IScheduleApp, task *ScheduleInTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Schedule{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetScheduleTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_SCHEDULE
			task.Result.Errmsg = err.Error()
			return nil
		}

		if v.Status == ScheduleStatusNone {

			v.Status = ScheduleStatusIn
			v.InTime = time.Now().Unix()

			_, err = kk.DBUpdateWithKeys(db, a.GetScheduleTable(), a.GetPrefix(), &v, map[string]bool{"status": true, "intime": true})

			if err != nil {
				task.Result.Errno = ERROR_SCHEDULE
				task.Result.Errmsg = err.Error()
				return nil
			}

		} else if v.Status != ScheduleStatusIn {
			task.Result.Errno = ERROR_SCHEDULE_STATUS
			task.Result.Errmsg = "The status can not be modified"
			return nil
		}

		task.Result.Schedule = &v

	} else {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND
		task.Result.Errmsg = "Not found schedule"
	}

	return nil
}

func (S *ScheduleService) HandleScheduleOffTask(a IScheduleApp, task *ScheduleOffTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Schedule{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetScheduleTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_SCHEDULE
			task.Result.Errmsg = err.Error()
			return nil
		}

		if v.Status == ScheduleStatusIn {

			v.Status = ScheduleStatusNone
			v.InTime = 0

			_, err = kk.DBUpdateWithKeys(db, a.GetScheduleTable(), a.GetPrefix(), &v, map[string]bool{"status": true, "intime": true})

			if err != nil {
				task.Result.Errno = ERROR_SCHEDULE
				task.Result.Errmsg = err.Error()
				return nil
			}

		} else if v.Status != ScheduleStatusNone {
			task.Result.Errno = ERROR_SCHEDULE_STATUS
			task.Result.Errmsg = "The status can not be modified"
			return nil
		}

		task.Result.Schedule = &v

	} else {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND
		task.Result.Errmsg = "Not found schedule"
	}

	return nil
}

func (S *ScheduleService) HandleScheduleStartTask(a IScheduleApp, task *ScheduleStartTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Schedule{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetScheduleTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_SCHEDULE
			task.Result.Errmsg = err.Error()
			return nil
		}

		if v.Status == ScheduleStatusIn || v.Status == ScheduleStatusFail || v.Status == ScheduleStatusEnd {

			v.Status = ScheduleStatusStart

			_, err = kk.DBUpdateWithKeys(db, a.GetScheduleTable(), a.GetPrefix(), &v, map[string]bool{"status": true})

			if err != nil {
				task.Result.Errno = ERROR_SCHEDULE
				task.Result.Errmsg = err.Error()
				return nil
			}

		} else if v.Status != ScheduleStatusStart {
			task.Result.Errno = ERROR_SCHEDULE_STATUS
			task.Result.Errmsg = "The status can not be modified"
			return nil
		}

		task.Result.Schedule = &v

	} else {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND
		task.Result.Errmsg = "Not found schedule"
	}

	return nil
}

func (S *ScheduleService) HandleScheduleEndTask(a IScheduleApp, task *ScheduleEndTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Schedule{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetScheduleTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_SCHEDULE
			task.Result.Errmsg = err.Error()
			return nil
		}

		if v.Status == ScheduleStatusIn || v.Status == ScheduleStatusFail || v.Status == ScheduleStatusStart {

			v.Status = ScheduleStatusEnd

			_, err = kk.DBUpdateWithKeys(db, a.GetScheduleTable(), a.GetPrefix(), &v, map[string]bool{"status": true})

			if err != nil {
				task.Result.Errno = ERROR_SCHEDULE
				task.Result.Errmsg = err.Error()
				return nil
			}

		} else if v.Status != ScheduleStatusEnd {
			task.Result.Errno = ERROR_SCHEDULE_STATUS
			task.Result.Errmsg = "The status can not be modified"
			return nil
		}

		task.Result.Schedule = &v

	} else {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND
		task.Result.Errmsg = "Not found schedule"
	}

	return nil
}

func (S *ScheduleService) HandleScheduleFailTask(a IScheduleApp, task *ScheduleFailTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Schedule{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetScheduleTable(), a.GetPrefix(), " WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_SCHEDULE
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_SCHEDULE
			task.Result.Errmsg = err.Error()
			return nil
		}

		if v.Status == ScheduleStatusIn || v.Status == ScheduleStatusStart {

			v.Status = ScheduleStatusFail

			_, err = kk.DBUpdateWithKeys(db, a.GetScheduleTable(), a.GetPrefix(), &v, map[string]bool{"status": true})

			if err != nil {
				task.Result.Errno = ERROR_SCHEDULE
				task.Result.Errmsg = err.Error()
				return nil
			}

		} else if v.Status != ScheduleStatusFail {
			task.Result.Errno = ERROR_SCHEDULE_STATUS
			task.Result.Errmsg = "The status can not be modified"
			return nil
		}

		task.Result.Schedule = &v

	} else {
		task.Result.Errno = ERROR_SCHEDULE_NOT_FOUND
		task.Result.Errmsg = "Not found schedule"
	}

	return nil
}
