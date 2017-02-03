package car

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"time"
)

type CarService struct {
	app.Service

	Create      *CarCreateTask
	Set         *CarSetTask
	Get         *CarTask
	Remove      *CarRemoveTask
	Query       *CarQueryTask
	LocationSet *CarLocationSetTask
}

func (S *CarService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *CarService) HandleCarCreateTask(a ICarApp, task *CarCreateTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CAR
		task.Result.Errmsg = err.Error()
		return nil
	}

	if task.PlateNo != "" {
		count, err := kk.DBQueryCount(db, a.GetCarTable(), a.GetPrefix(), " WHERE plateno=?", task.PlateNo)
		if err != nil {
			task.Result.Errno = ERROR_CAR
			task.Result.Errmsg = err.Error()
			return nil
		}
		if count > 0 {
			task.Result.Errno = ERROR_CAR_PLATENO
			task.Result.Errmsg = "Plate no already exists"
			return nil
		}
	}

	if task.LicenceCode != "" {
		count, err := kk.DBQueryCount(db, a.GetCarTable(), a.GetPrefix(), " WHERE licencecode=?", task.LicenceCode)
		if err != nil {
			task.Result.Errno = ERROR_CAR
			task.Result.Errmsg = err.Error()
			return nil
		}
		if count > 0 {
			task.Result.Errno = ERROR_CAR_LICENCE_CODE
			task.Result.Errmsg = "Licence code already exists"
			return nil
		}
	}

	v := Car{}

	v.LicenceCode = task.LicenceCode
	v.Brand = task.Brand
	v.Name = task.Name
	v.PlateNo = task.PlateNo
	v.Ctime = time.Now().Unix()
	v.Atime = v.Ctime

	_, err = kk.DBInsert(db, a.GetCarTable(), a.GetPrefix(), &v)

	if err != nil {
		task.Result.Errno = ERROR_CAR
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.Car = &v

	return nil
}

func (S *CarService) HandleCarTask(a ICarApp, task *CarTask) error {

	if task.Id == 0 && task.PlateNo == "" && task.LicenceCode == "" {
		task.Result.Errno = ERROR_CAR_NOT_FOUND_ID
		task.Result.Errmsg = "Not Found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CAR
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Car{}

	sql := bytes.NewBuffer(nil)

	args := []interface{}{}

	sql.WriteString(" WHERE 1")

	if task.Id != 0 {
		sql.WriteString(" AND id=?")
		args = append(args, task.Id)
	}

	if task.PlateNo != "" {
		sql.WriteString(" AND plateno=?")
		args = append(args, task.PlateNo)
	}

	if task.LicenceCode != "" {
		sql.WriteString(" AND licencecode=?")
		args = append(args, task.LicenceCode)
	}

	rows, err := kk.DBQuery(db, a.GetCarTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_CAR
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_CAR
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Car = &v

	} else {
		task.Result.Errno = ERROR_CAR_NOT_FOUND
		task.Result.Errmsg = "Not Found car"
		return nil
	}

	return nil
}

func (S *CarService) HandleCarSetTask(a ICarApp, task *CarSetTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_CAR_NOT_FOUND_ID
		task.Result.Errmsg = "Not Found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CAR
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Car{}

	rows, err := kk.DBQuery(db, a.GetCarTable(), a.GetPrefix(), "WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_CAR
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_CAR
			task.Result.Errmsg = err.Error()
			return nil
		}

		keys := map[string]bool{}

		if task.Name != nil {
			v.Name = dynamic.StringValue(task.Name, v.Name)
			keys["name"] = true
		}

		if task.PlateNo != nil {

			v.PlateNo = dynamic.StringValue(task.PlateNo, v.PlateNo)

			if v.PlateNo != "" {

				count, err := kk.DBQueryCount(db, a.GetCarTable(), a.GetPrefix(), " WHERE plateno=? AND id <> ?", v.PlateNo, v.Id)

				if err != nil {
					task.Result.Errno = ERROR_CAR
					task.Result.Errmsg = err.Error()
					return nil
				}
				if count > 0 {
					task.Result.Errno = ERROR_CAR_PLATENO
					task.Result.Errmsg = "Plate noalready exists"
					return nil
				}

			}

			keys["plateno"] = true
		}

		if task.LicenceCode != nil {

			v.LicenceCode = dynamic.StringValue(task.LicenceCode, v.LicenceCode)

			if v.LicenceCode != "" {
				count, err := kk.DBQueryCount(db, a.GetCarTable(), a.GetPrefix(), " WHERE licencecode=? AND id <> ?", v.LicenceCode, v.Id)
				if err != nil {
					task.Result.Errno = ERROR_CAR
					task.Result.Errmsg = err.Error()
					return nil
				}
				if count > 0 {
					task.Result.Errno = ERROR_CAR_LICENCE_CODE
					task.Result.Errmsg = "Licence code already exists"
					return nil
				}

			}

			keys["licencecode"] = true
		}

		if task.Brand != nil {
			v.Brand = dynamic.StringValue(task.Brand, v.Brand)
			keys["brand"] = true
		}

		if task.Capacity != nil {
			v.Capacity = dynamic.StringValue(task.Capacity, v.Capacity)
			keys["capacity"] = true
		}

		_, err = kk.DBUpdateWithKeys(db, a.GetCarTable(), a.GetPrefix(), &v, keys)

		if err != nil {
			task.Result.Errno = ERROR_CAR
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Car = &v

	} else {
		task.Result.Errno = ERROR_CAR_NOT_FOUND
		task.Result.Errmsg = "Not Found car"
		return nil
	}

	return nil
}

func (S *CarService) HandleCarRemoveTask(a ICarApp, task *CarRemoveTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_CAR_NOT_FOUND_ID
		task.Result.Errmsg = "Not Found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CAR
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Car{}

	rows, err := kk.DBQuery(db, a.GetCarTable(), a.GetPrefix(), "WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_CAR
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_CAR
			task.Result.Errmsg = err.Error()
			return nil
		}

		{
			t := TriggerCarRemovingTask{}
			t.Car = &v

			err := app.Handle(a, &t)

			if err != nil {
				e, ok := err.(*app.Error)
				if ok {
					task.Result.Errno = e.Errno
					task.Result.Errmsg = e.Errmsg
					return nil
				}
				task.Result.Errno = ERROR_CAR
				task.Result.Errmsg = err.Error()
				return nil
			}

		}

		_, err = kk.DBDelete(db, a.GetCarTable(), a.GetPrefix(), " WHERE id=?", task.Id)

		if err != nil {
			task.Result.Errno = ERROR_CAR
			task.Result.Errmsg = err.Error()
			return nil
		}

		{
			t := TriggerCarRemovedTask{}
			t.Car = &v

			err := app.Handle(a, &t)

			if err != nil {
				e, ok := err.(*app.Error)
				if ok {
					task.Result.Errno = e.Errno
					task.Result.Errmsg = e.Errmsg
					return nil
				}
				task.Result.Errno = ERROR_CAR
				task.Result.Errmsg = err.Error()
				return nil
			}

		}

		task.Result.Car = &v

	} else {
		task.Result.Errno = ERROR_CAR_NOT_FOUND
		task.Result.Errmsg = "Not Found car"
		return nil
	}

	return nil
}

func (S *CarService) HandleCarQueryTask(a ICarApp, task *CarQueryTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CAR
		task.Result.Errmsg = err.Error()
		return nil
	}

	var cars = []Car{}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(" WHERE 1")

	if task.Id != 0 {
		sql.WriteString(" AND id=?")
		args = append(args, task.Id)
	}

	if task.PlateNo != "" {
		sql.WriteString(" AND plateno=?")
		args = append(args, task.PlateNo)
	}

	if task.LicenceCode != "" {
		sql.WriteString(" AND licencecode=?")
		args = append(args, task.LicenceCode)
	}

	if task.Keyword != "" {
		q := "%" + task.Keyword + "%"
		sql.WriteString(" AND (id=? OR plateno LIKE ? OR licencecode LIKE ? OR name LIKE ? OR brand LIKE ?)")
		args = append(args, task.Keyword, q, q, q, q)
	}

	if task.OrderBy == "asc" {
		sql.WriteString(" ORDER BY id ASC")
	} else {
		sql.WriteString(" ORDER BY id DESC")
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
		var counter = CarQueryCounter{}
		counter.PageIndex = pageIndex
		counter.PageSize = pageSize
		total, err := kk.DBQueryCount(db, a.GetCarTable(), a.GetPrefix(), sql.String(), args...)
		if err != nil {
			task.Result.Errno = ERROR_CAR
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

	sql.WriteString(fmt.Sprintf(" LIMIT %d,%d", (pageIndex-1)*pageSize, pageSize))

	var v = Car{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetCarTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_CAR
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	for rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_CAR
			task.Result.Errmsg = err.Error()
			return nil
		}

		cars = append(cars, v)
	}

	task.Result.Cars = cars

	return nil
}

func (S *CarService) HandleCarLocationSetTask(a ICarApp, task *CarLocationSetTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_CAR_NOT_FOUND_ID
		task.Result.Errmsg = "Not Found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_CAR
		task.Result.Errmsg = err.Error()
		return nil
	}

	_, err = db.Exec(fmt.Sprintf("UPDATE %s%s SET longitude=?,latitude=?,ip=? WHERE id=?", a.GetPrefix(), a.GetCarTable().Name), task.Longitude, task.Latitude, task.Ip, task.Id)

	if err != nil {
		task.Result.Errno = ERROR_CAR
		task.Result.Errmsg = err.Error()
		return nil
	}

	return nil
}
