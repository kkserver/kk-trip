package driver

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"time"
)

type DriverService struct {
	app.Service

	Create *DriverCreateTask
	Set    *DriverSetTask
	Get    *DriverTask
	Remove *DriverRemoveTask
	Query  *DriverQueryTask
}

func (S *DriverService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *DriverService) HandleDriverCreateTask(a IDriverApp, task *DriverCreateTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_DRIVER
		task.Result.Errmsg = err.Error()
		return nil
	}

	if task.Phone != "" {
		count, err := kk.DBQueryCount(db, a.GetDriverTable(), a.GetPrefix(), " WHERE phone=?", task.Phone)
		if err != nil {
			task.Result.Errno = ERROR_DRIVER
			task.Result.Errmsg = err.Error()
			return nil
		}
		if count > 0 {
			task.Result.Errno = ERROR_DRIVER_PHONE
			task.Result.Errmsg = "Phone number already exists"
			return nil
		}
	}

	if task.Code != "" {
		count, err := kk.DBQueryCount(db, a.GetDriverTable(), a.GetPrefix(), " WHERE code=?", task.Code)
		if err != nil {
			task.Result.Errno = ERROR_DRIVER
			task.Result.Errmsg = err.Error()
			return nil
		}
		if count > 0 {
			task.Result.Errno = ERROR_DRIVER_CODE
			task.Result.Errmsg = "ID code already exists"
			return nil
		}
	}

	if task.LicenceCode != "" {
		count, err := kk.DBQueryCount(db, a.GetDriverTable(), a.GetPrefix(), " WHERE licencecode=?", task.LicenceCode)
		if err != nil {
			task.Result.Errno = ERROR_DRIVER
			task.Result.Errmsg = err.Error()
			return nil
		}
		if count > 0 {
			task.Result.Errno = ERROR_DRIVER_LICENCE_CODE
			task.Result.Errmsg = "Licence code already exists"
			return nil
		}
	}

	v := Driver{}

	v.Name = task.Name
	v.Phone = task.Phone
	v.Code = task.Code
	v.LicenceCode = task.LicenceCode
	v.Ctime = time.Now().Unix()

	_, err = kk.DBInsert(db, a.GetDriverTable(), a.GetPrefix(), &v)

	if err != nil {
		task.Result.Errno = ERROR_DRIVER
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.Driver = &v

	return nil
}

func (S *DriverService) HandleDriverTask(a IDriverApp, task *DriverTask) error {

	if task.Id == 0 && task.Phone == "" && task.Code == "" && task.LicenceCode == "" {
		task.Result.Errno = ERROR_DRIVER_NOT_FOUND_ID
		task.Result.Errmsg = "Not Found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_DRIVER
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Driver{}

	sql := bytes.NewBuffer(nil)

	args := []interface{}{}

	sql.WriteString(" WHERE 1")

	if task.Id != 0 {
		sql.WriteString(" AND id=?")
		args = append(args, task.Id)
	}

	if task.Phone != "" {
		sql.WriteString(" AND phone=?")
		args = append(args, task.Phone)
	}

	if task.Code != "" {
		sql.WriteString(" AND code=?")
		args = append(args, task.Code)
	}

	if task.LicenceCode != "" {
		sql.WriteString(" AND licencecode=?")
		args = append(args, task.LicenceCode)
	}

	rows, err := kk.DBQuery(db, a.GetDriverTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_DRIVER
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_DRIVER
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Driver = &v

	} else {
		task.Result.Errno = ERROR_DRIVER_NOT_FOUND
		task.Result.Errmsg = "Not Found drirver"
		return nil
	}

	return nil
}

func (S *DriverService) HandleDriverSetTask(a IDriverApp, task *DriverSetTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_DRIVER_NOT_FOUND_ID
		task.Result.Errmsg = "Not Found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_DRIVER
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Driver{}

	rows, err := kk.DBQuery(db, a.GetDriverTable(), a.GetPrefix(), "WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_DRIVER
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_DRIVER
			task.Result.Errmsg = err.Error()
			return nil
		}

		keys := map[string]bool{}

		if task.Name != nil {
			v.Name = dynamic.StringValue(task.Name, v.Name)
			keys["name"] = true
		}

		if task.Code != nil {

			count, err := kk.DBQueryCount(db, a.GetDriverTable(), a.GetPrefix(), " WHERE code=? AND id <> ?", task.Code, v.Id)

			if err != nil {
				task.Result.Errno = ERROR_DRIVER
				task.Result.Errmsg = err.Error()
				return nil
			}
			if count > 0 {
				task.Result.Errno = ERROR_DRIVER_CODE
				task.Result.Errmsg = "ID code already exists"
				return nil
			}

			v.Code = dynamic.StringValue(task.Code, v.Code)
			keys["code"] = true
		}

		if task.LicenceCode != nil {

			count, err := kk.DBQueryCount(db, a.GetDriverTable(), a.GetPrefix(), " WHERE licencecode=? AND id <> ?", task.LicenceCode, v.Id)
			if err != nil {
				task.Result.Errno = ERROR_DRIVER
				task.Result.Errmsg = err.Error()
				return nil
			}
			if count > 0 {
				task.Result.Errno = ERROR_DRIVER_LICENCE_CODE
				task.Result.Errmsg = "Licence code already exists"
				return nil
			}

			v.LicenceCode = dynamic.StringValue(task.LicenceCode, v.LicenceCode)
			keys["licencecode"] = true
		}

		if task.Phone != nil {

			count, err := kk.DBQueryCount(db, a.GetDriverTable(), a.GetPrefix(), " WHERE phone=? AND id <> ?", task.Phone, v.Id)

			if err != nil {
				task.Result.Errno = ERROR_DRIVER
				task.Result.Errmsg = err.Error()
				return nil
			}

			if count > 0 {
				task.Result.Errno = ERROR_DRIVER_PHONE
				task.Result.Errmsg = "Phone number already exists"
				return nil
			}

			v.Phone = dynamic.StringValue(task.Phone, v.Phone)
			keys["phone"] = true
		}

		_, err = kk.DBUpdateWithKeys(db, a.GetDriverTable(), a.GetPrefix(), &v, keys)

		if err != nil {
			task.Result.Errno = ERROR_DRIVER
			task.Result.Errmsg = err.Error()
			return nil
		}

		task.Result.Driver = &v

	} else {
		task.Result.Errno = ERROR_DRIVER_NOT_FOUND
		task.Result.Errmsg = "Not Found drirver"
		return nil
	}

	return nil
}

func (S *DriverService) HandleDriverRemoveTask(a IDriverApp, task *DriverRemoveTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_DRIVER_NOT_FOUND_ID
		task.Result.Errmsg = "Not Found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_DRIVER
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Driver{}

	rows, err := kk.DBQuery(db, a.GetDriverTable(), a.GetPrefix(), "WHERE id=?", task.Id)

	if err != nil {
		task.Result.Errno = ERROR_DRIVER
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_DRIVER
			task.Result.Errmsg = err.Error()
			return nil
		}

		{
			t := TriggerDriverRemovingTask{}
			t.Driver = &v

			err := app.Handle(a, &t)

			if err != nil {
				e, ok := err.(*app.Error)
				if ok {
					task.Result.Errno = e.Errno
					task.Result.Errmsg = e.Errmsg
					return nil
				}
				task.Result.Errno = ERROR_DRIVER
				task.Result.Errmsg = err.Error()
				return nil
			}

		}

		_, err = kk.DBDelete(db, a.GetDriverTable(), a.GetPrefix(), " WHERE id=?", task.Id)

		if err != nil {
			task.Result.Errno = ERROR_DRIVER
			task.Result.Errmsg = err.Error()
			return nil
		}

		{
			t := TriggerDriverRemovedTask{}
			t.Driver = &v

			err := app.Handle(a, &t)

			if err != nil {
				e, ok := err.(*app.Error)
				if ok {
					task.Result.Errno = e.Errno
					task.Result.Errmsg = e.Errmsg
					return nil
				}
				task.Result.Errno = ERROR_DRIVER
				task.Result.Errmsg = err.Error()
				return nil
			}

		}

		task.Result.Driver = &v

	} else {
		task.Result.Errno = ERROR_DRIVER_NOT_FOUND
		task.Result.Errmsg = "Not Found drirver"
		return nil
	}

	return nil
}

func (S *DriverService) HandleDriverQueryTask(a IDriverApp, task *DriverQueryTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_DRIVER
		task.Result.Errmsg = err.Error()
		return nil
	}

	var drivers = []Driver{}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(" WHERE 1")

	if task.Id != 0 {
		sql.WriteString(" AND id=?")
		args = append(args, task.Id)
	}

	if task.Name != "" {
		sql.WriteString(" AND name=?")
		args = append(args, task.Name)
	}

	if task.Phone != "" {
		sql.WriteString(" AND phone=?")
		args = append(args, task.Phone)
	}

	if task.LicenceCode != "" {
		sql.WriteString(" AND licencecode=?")
		args = append(args, task.LicenceCode)
	}

	if task.Code != "" {
		sql.WriteString(" AND code=?")
		args = append(args, task.Code)
	}

	if task.Keyword != "" {
		q := "%" + task.Keyword + "%"
		sql.WriteString(" AND (id=? OR code LIKE ? OR name LIKE ? OR phone LIKE ? OR licencecode LIKE ?)")
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
		var counter = DriverQueryCounter{}
		counter.PageIndex = pageIndex
		counter.PageSize = pageSize
		total, err := kk.DBQueryCount(db, a.GetDriverTable(), a.GetPrefix(), sql.String(), args...)
		if err != nil {
			task.Result.Errno = ERROR_DRIVER
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

	var v = Driver{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetDriverTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_DRIVER
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	for rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_DRIVER
			task.Result.Errmsg = err.Error()
			return nil
		}

		drivers = append(drivers, v)
	}

	task.Result.Drivers = drivers

	return nil
}
