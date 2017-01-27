package ticket

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"github.com/kkserver/kk-order/order"
	"github.com/kkserver/kk-trip/trip/line"
	"github.com/kkserver/kk-trip/trip/schedule"
	"log"
	"strings"
	"time"
)

type TicketService struct {
	app.Service

	Create    *TicketCreateTask
	PreCreate *TicketPreCreateTask
	Get       *TicketTask
	Refund    *TicketRefundTask
	Query     *TicketQueryTask
	Calendar  *TicketCalendarTask
}

func (S *TicketService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *TicketService) HandleTicketPreCreateTask(a ITicketApp, task *TicketPreCreateTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TICKET
		task.Result.Errmsg = err.Error()
		return nil
	}

	items := task.Items

	if items == nil {
		items = []TicketItem{}
	}

	if task.Text != "" {
		vs := strings.Split(task.Text, ";")

		for _, v := range vs {
			is := strings.Split(v, ":")
			i := TicketItem{}
			if len(is) > 0 {
				i.LineId = dynamic.IntValue(is[0], 0)
			}
			if len(is) > 1 {
				i.Date = dynamic.IntValue(is[1], 0)
			}
			if len(is) > 2 {
				i.SeatNo = dynamic.StringValue(is[2], "")
			}
			items = append(items, i)
		}
	}

	if len(items) == 0 {
		task.Result.Errno = ERROR_TICKET_NOT_FOUND
		task.Result.Errmsg = "Not Found Ticket"
		return nil
	}

	values := []TicketValue{}
	value := int64(0)
	payValue := int64(0)

	err = func() error {

		schedules := []*schedule.Schedule{}
		lines := map[int64]*line.Line{}

		getSchedule := func(lineId int64, date int64) (*schedule.Schedule, error) {

			for _, v := range schedules {
				if v.LineId == lineId && v.Date == date {
					return v, nil
				}
			}

			rows, err := kk.DBQuery(db, a.GetScheduleTable(), a.GetPrefix(), " WHERE lineid=? AND date=?", lineId, date)

			if err != nil {
				return nil, err
			}

			defer rows.Close()

			if rows.Next() {

				v := schedule.Schedule{}
				scanner := kk.NewDBScaner(&v)
				err = scanner.Scan(rows)

				if err != nil {
					return nil, err
				}

				schedules = append(schedules, &v)

				return &v, nil

			} else {
				return nil, app.NewError(ERROR_TICKET_NOT_FOUND_SCHEDULE, "Not Found Schedule")
			}
		}

		getLine := func(lineId int64) (*line.Line, error) {

			v, ok := lines[lineId]

			if ok {
				return v, nil
			}

			rows, err := kk.DBQuery(db, a.GetLineTable(), a.GetPrefix(), " WHERE id=?", lineId)

			if err != nil {
				return nil, err
			}

			defer rows.Close()

			if rows.Next() {

				v := line.Line{}
				scanner := kk.NewDBScaner(&v)
				err = scanner.Scan(rows)

				if err != nil {
					return nil, err
				}

				lines[lineId] = &v

				return &v, nil

			} else {
				return nil, app.NewError(ERROR_TICKET_NOT_FOUND_LINE, "Not Found Line")
			}
		}

		for _, i := range items {

			line, err := getLine(i.LineId)

			if err != nil {
				return err
			}

			v, err := getSchedule(i.LineId, i.Date)

			if err != nil {
				return err
			}

			if v.Count >= v.MaxCount {
				return app.NewError(ERROR_TICKET_SCHEDULE_MAX_COUNT, "No tickets available")
			}

			count, err := kk.DBQueryCount(db, a.GetTicketTable(), a.GetPrefix(), " WHERE scheduleid=? AND uid=? AND status IN (?,?)", v.Id, task.Uid, TicketStatusNone, TicketStatusPay)

			if count >= v.UMaxCount {
				return app.NewError(ERROR_TICKET_SCHEDULE_USER_MAX_COUNT, "No tickets available")
			}

			if i.SeatNo == "" {
				i.SeatNo = fmt.Sprintf("%d", v.Count+1)
			}

			count, err = kk.DBQueryCount(db, a.GetTicketTable(), a.GetPrefix(), " WHERE scheduleid=? AND seatno=? AND status IN (?,?)", v.Id, i.SeatNo, TicketStatusNone, TicketStatusPay)

			if count > 0 {
				return app.NewError(ERROR_TICKET_SCHEDULE_SEATNO, "Seat number already exists")
			}

			vv := TicketValue{}
			vv.Value = line.Price
			vv.PayValue = line.Price
			vv.RefundValue = line.Price

			values = append(values, vv)

			value = value + vv.Value
			payValue = payValue + vv.PayValue

		}

		return nil

	}()

	if err != nil {

		e, ok := err.(*app.Error)
		if ok {
			task.Result.Errno = e.Errno
			task.Result.Errmsg = e.Errmsg
			return nil
		} else {
			task.Result.Errno = ERROR_TICKET
			task.Result.Errmsg = err.Error()
			return nil
		}
	}

	task.Result.Value = value
	task.Result.PayValue = payValue
	task.Result.Values = values

	return nil
}

func (S *TicketService) HandleTicketCreateTask(a ITicketApp, task *TicketCreateTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TICKET
		task.Result.Errmsg = err.Error()
		return nil
	}

	items := task.Items

	if items == nil {
		items = []TicketItem{}
	}

	if task.Text != "" {
		vs := strings.Split(task.Text, ";")

		for _, v := range vs {
			is := strings.Split(v, ":")
			i := TicketItem{}
			if len(is) > 0 {
				i.LineId = dynamic.IntValue(is[0], 0)
			}
			if len(is) > 1 {
				i.Date = dynamic.IntValue(is[1], 0)
			}
			if len(is) > 2 {
				i.SeatNo = dynamic.StringValue(is[2], "")
			}
			items = append(items, i)
		}
	}

	if len(items) == 0 {
		task.Result.Errno = ERROR_TICKET_NOT_FOUND
		task.Result.Errmsg = "Not Found Ticket"
		return nil
	}

	var tickets = []Ticket{}
	var value int64 = 0
	var payValue int64 = 0
	var values []TicketValue = nil

	{
		t := TicketPreCreateTask{}
		t.Items = items
		app.Handle(a, &t)
		if t.Result.Values != nil {
			values = t.Result.Values
			value = t.Result.Value
			payValue = t.Result.PayValue
		} else {
			task.Result.Errno = t.Result.Errno
			task.Result.Errmsg = t.Result.Errmsg
			return nil
		}
	}

	tx, err := db.Begin()

	if err != nil {
		task.Result.Errno = ERROR_TICKET
		task.Result.Errmsg = err.Error()
		return nil
	}

	var odr *order.Order = nil

	if task.Expires == 0 {
		task.Expires = 3600
	}

	{
		t := order.OrderCreateTask{}
		t.Uid = task.Uid
		t.Expires = task.Expires
		t.Value = value
		t.PayValue = payValue
		t.RefundValue = payValue
		t.Title = "[车票]"

		app.Handle(a, &t)

		if t.Result.Order != nil {
			odr = t.Result.Order
		} else {
			task.Result.Errno = ERROR_TICKET_ORDER_CREATE
			task.Result.Errmsg = "Can not create order"
			return nil
		}
	}

	err = func() error {

		schedules := []*schedule.Schedule{}

		getSchedule := func(lineId int64, date int64) (*schedule.Schedule, error) {

			for _, v := range schedules {
				if v.LineId == lineId && v.Date == date {
					return v, nil
				}
			}

			rows, err := kk.DBQuery(db, a.GetScheduleTable(), a.GetPrefix(), " WHERE lineid=? AND date=? FOR UPDATE", lineId, date)

			if err != nil {
				return nil, err
			}

			defer rows.Close()

			if rows.Next() {

				v := schedule.Schedule{}
				scanner := kk.NewDBScaner(&v)
				err = scanner.Scan(rows)

				if err != nil {
					return nil, err
				}

				schedules = append(schedules, &v)

				return &v, nil

			} else {
				return nil, app.NewError(ERROR_TICKET_NOT_FOUND_SCHEDULE, "Not Found Schedule")
			}
		}

		for i, item := range items {

			v, err := getSchedule(item.LineId, item.Date)

			if err != nil {
				return err
			}

			if v.Count >= v.MaxCount {
				return app.NewError(ERROR_TICKET_SCHEDULE_MAX_COUNT, "No tickets available")
			}

			count, err := kk.DBQueryCount(db, a.GetTicketTable(), a.GetPrefix(), " WHERE scheduleid=? AND uid=? AND status IN (?,?)", v.Id, task.Uid, TicketStatusNone, TicketStatusPay)

			if count >= v.UMaxCount {
				return app.NewError(ERROR_TICKET_SCHEDULE_USER_MAX_COUNT, "No tickets available")
			}

			if item.SeatNo == "" {
				item.SeatNo = fmt.Sprintf("%d", v.Count+1)
			}

			count, err = kk.DBQueryCount(db, a.GetTicketTable(), a.GetPrefix(), " WHERE scheduleid=? AND seatno=? AND status IN (?,?)", v.Id, item.SeatNo, TicketStatusNone, TicketStatusPay)

			if count > 0 {
				return app.NewError(ERROR_TICKET_SCHEDULE_SEATNO, "Seat number already exists")
			}

			vv := Ticket{}
			vv.ScheduleId = v.Id
			vv.LineId = v.LineId
			vv.Date = v.Date
			vv.Uid = task.Uid
			vv.OrderId = odr.Id
			vv.SeatNo = item.SeatNo
			vv.PayValue = values[i].PayValue
			vv.Value = values[i].Value
			vv.RefundValue = values[i].RefundValue
			vv.Ctime = time.Now().Unix()

			_, err = kk.DBInsert(db, a.GetTicketTable(), a.GetPrefix(), &vv)

			if err != nil {
				return err
			}

			v.Count = v.Count + 1

			tickets = append(tickets, vv)
		}

		keys := map[string]bool{"count": true}

		for _, schedule := range schedules {

			_, err = kk.DBUpdateWithKeys(db, a.GetScheduleTable(), a.GetPrefix(), schedule, keys)

			if err != nil {
				return err
			}
		}

		return nil
	}()

	if err == nil {
		err = tx.Commit()
	}

	if err != nil {

		tx.Rollback()

		{
			t := order.OrderCancelTask{}
			t.Id = odr.Id
			app.Handle(a, &t)
		}

		e, ok := err.(*app.Error)
		if ok {
			task.Result.Errno = e.Errno
			task.Result.Errmsg = e.Errmsg
			return nil
		} else {
			task.Result.Errno = ERROR_TICKET
			task.Result.Errmsg = err.Error()
			return nil
		}
	} else {

		t := order.OrderSetTask{}
		t.Id = odr.Id
		t.Type = "ticket"

		options := map[interface{}]interface{}{}

		items := []map[interface{}]interface{}{}

		for _, ticket := range tickets {
			items = append(items, map[interface{}]interface{}{"id": ticket.Id, "date": ticket.Date, "lineId": ticket.LineId})
		}

		options["items"] = items

		t.Options = options

		app.Handle(a, &t)

		if t.Result.Order != nil {
			odr = t.Result.Order
		}
	}

	task.Result.Order = odr
	task.Result.Tickets = tickets

	return nil
}

func (S *TicketService) HandleTicketTask(a ITicketApp, task *TicketTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_TICKET_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TICKET
		task.Result.Errmsg = err.Error()
		return nil
	}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(" WHERE id=?")

	args = append(args, task.Id)

	if task.Uid != nil {
		sql.WriteString(" WHERE uid=?")
		args = append(args, task.Uid)
	}

	rows, err := kk.DBQuery(db, a.GetTicketTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_TICKET
		task.Result.Errmsg = err.Error()
		return nil
	}

	v := Ticket{}

	defer rows.Close()

	if rows.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_TICKET
			task.Result.Errmsg = err.Error()
			return nil
		}

	} else {
		return app.NewError(ERROR_TICKET_NOT_FOUND, "Not Found Ticket")
	}

	task.Result.Ticket = &v

	return nil
}

func (S *TicketService) HandleTicketQueryTask(a ITicketApp, task *TicketQueryTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TICKET
		task.Result.Errmsg = err.Error()
		return nil
	}

	var tickets = []Ticket{}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(" WHERE 1")

	if task.Id != 0 {
		sql.WriteString(" AND id=?")
		args = append(args, task.Id)
	}

	if task.Uid != nil {
		sql.WriteString(" AND uid=?")
		args = append(args, task.Uid)
	}

	if task.OrderId != nil {
		sql.WriteString(" AND orderid=?")
		args = append(args, task.OrderId)
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
		var counter = TicketQueryCounter{}
		counter.PageIndex = pageIndex
		counter.PageSize = pageSize
		total, err := kk.DBQueryCount(db, a.GetTicketTable(), a.GetPrefix(), sql.String(), args...)
		if err != nil {
			task.Result.Errno = ERROR_TICKET
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

	var v = Ticket{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetTicketTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_TICKET
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	for rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_TICKET
			task.Result.Errmsg = err.Error()
			return nil
		}

		tickets = append(tickets, v)
	}

	task.Result.Tickets = tickets

	return nil
}

func (S *TicketService) HandleTicketRefundTask(a ITicketApp, task *TicketRefundTask) error {

	if task.Id == 0 {
		task.Result.Errno = ERROR_TICKET_NOT_FOUND_ID
		task.Result.Errmsg = "Not found id"
		return nil
	}

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TICKET
		task.Result.Errmsg = err.Error()
		return nil
	}

	tx, err := db.Begin()

	if err != nil {
		task.Result.Errno = ERROR_TICKET
		task.Result.Errmsg = err.Error()
		return nil
	}

	var v = Ticket{}

	err = func() error {

		rows, err := kk.DBQuery(db, a.GetTicketTable(), a.GetPrefix(), " WHERE id=? FOR UPDATE", task.Id)

		if err != nil {
			return err
		}

		defer rows.Close()

		if rows.Next() {

			scanner := kk.NewDBScaner(&v)

			err = scanner.Scan(rows)

			if err != nil {
				return err
			}

			if v.Status != TicketStatusPay {
				return app.NewError(ERROR_TICKET_STATUS, "The current state can not be modified")
			}

			v.Status = TicketStatusRefund

			_, err = kk.DBUpdateWithKeys(db, a.GetTicketTable(), a.GetPrefix(), &v, map[string]bool{"status": true})

			if err != nil {
				return err
			}

		} else {
			return app.NewError(ERROR_TICKET_NOT_FOUND, "Not Found Ticket")
		}

		return nil
	}()

	if err == nil {
		err = tx.Commit()
	}

	if err != nil {
		tx.Rollback()
		e, ok := err.(*app.Error)
		if ok {
			task.Result.Errno = e.Errno
			task.Result.Errmsg = e.Errmsg
			return nil
		} else {
			task.Result.Errno = ERROR_TICKET
			task.Result.Errmsg = err.Error()
			return nil
		}
	}

	task.Result.Ticket = &v

	return nil
}

func (S *TicketService) HandleTriggerOrderTimeoutDidTask(a ITicketApp, task *order.TriggerOrderTimeoutDidTask) error {

	log.Println("TriggerOrderTimeoutDidTask")

	var db, err = a.GetDB()

	if err != nil {
		return nil
	}

	tx, err := db.Begin()

	if err != nil {
		return nil
	}

	err = func() error {

		rows, err := db.Query(fmt.Sprintf("SELECT scheduleid, COUNT(id) FROM %s%s WHERE orderid=? AND status=? GROUP BY scheduleid", a.GetPrefix(), a.GetTicketTable()), task.Order.Id, TicketStatusNone)

		if err != nil {
			return err
		}

		defer rows.Close()

		var id int64 = 0
		var count int64 = 0

		for rows.Next() {
			err = rows.Scan(&id, &count)
			if err != nil {
				return err
			}

			log.Println(id, count)

			_, err = db.Exec(fmt.Sprintf("UPDATE %s%s SET `count` = `count` - ? WHERE id=?", a.GetPrefix(), a.GetScheduleTable().Name), count, id)
			if err != nil {
				return err
			}
		}

		_, err = db.Exec(fmt.Sprintf("UPDATE %s%s SET status=? WHERE status=? AND orderid=?", a.GetPrefix(), a.GetTicketTable().Name), TicketStatusTimeout, TicketStatusNone, task.Order.Id)

		if err != nil {
			return nil
		}

		return nil
	}()

	if err == nil {
		err = tx.Commit()
	}

	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return err
		}
	}

	return nil
}

func (S *TicketService) HandleTriggerOrderPayDidTask(a ITicketApp, task *order.TriggerOrderPayDidTask) error {

	var db, err = a.GetDB()

	if err != nil {
		return nil
	}

	_, err = db.Exec(fmt.Sprintf("UPDATE %s%s SET status=? WHERE status=? AND orderid=?", a.GetPrefix(), a.GetTicketTable().Name), TicketStatusPay, TicketStatusNone, task.Order.Id)

	if err != nil {
		return nil
	}

	return nil
}

func (S *TicketService) HandleTriggerOrderCancelDidTask(a ITicketApp, task *order.TriggerOrderCancelDidTask) error {

	log.Println("TriggerOrderCancelDidTask")

	var db, err = a.GetDB()

	if err != nil {
		return nil
	}

	tx, err := db.Begin()

	if err != nil {
		return nil
	}

	err = func() error {

		rows, err := db.Query(fmt.Sprintf("SELECT scheduleid, COUNT(id) FROM %s%s WHERE orderid=? AND status=? GROUP BY scheduleid", a.GetPrefix(), a.GetTicketTable()), task.Order.Id, TicketStatusNone)

		if err != nil {
			return err
		}

		defer rows.Close()

		var id int64 = 0
		var count int64 = 0

		for rows.Next() {
			err = rows.Scan(&id, &count)
			if err != nil {
				return err
			}

			log.Println(id, count)

			_, err = db.Exec(fmt.Sprintf("UPDATE %s%s SET `count` = `count` - ? WHERE id=?", a.GetPrefix(), a.GetScheduleTable().Name), count, id)
			if err != nil {
				return err
			}
		}

		_, err = db.Exec(fmt.Sprintf("UPDATE %s%s SET status=? WHERE status=? AND orderid=?", a.GetPrefix(), a.GetTicketTable().Name), TicketStatusCancel, TicketStatusNone, task.Order.Id)

		if err != nil {
			return nil
		}

		return nil
	}()

	if err == nil {
		err = tx.Commit()
	}

	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return err
		}
	}

	return nil
}

func (S *TicketService) HandleTicketCalendarTask(a ITicketApp, task *TicketCalendarTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_TICKET
		task.Result.Errmsg = err.Error()
		return nil
	}

	var rows = Calendar{}

	var now = time.Now()

	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var min = time.Unix(0, 0)
	var max = time.Unix(0, 0)
	var days = map[time.Time]CalendarCell{}

	var v = CalendarCell{}
	var scanner = kk.NewDBScaner(&v)

	rs, err := kk.DBQuery(db, a.GetScheduleTable(), a.GetPrefix(),
		" WHERE lineid=? AND date>=? AND status=?", task.LineId, now.Unix(), schedule.ScheduleStatusIn)

	if err != nil {
		task.Result.Errno = ERROR_TICKET
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rs.Close()

	for rs.Next() {

		err = scanner.Scan(rs)

		if err != nil {
			task.Result.Errno = ERROR_TICKET
			task.Result.Errmsg = err.Error()
			return nil
		}

		t := time.Unix(v.Date, 0)

		v.Day.SetTime(t)

		if task.Uid != 0 {
			v.Ucount, err = kk.DBQueryCount(db, a.GetTicketTable(), a.GetPrefix(), " WHERE scheduleid=? AND uid=? AND status IN (?,?)", v.Id, task.Uid, TicketStatusNone, TicketStatusPay)
		}

		days[t] = v

		if min.Unix() == 0 {
			min = t
		} else if t.Unix() < min.Unix() {
			min = t
		}

		if min.Unix() == 0 {
			max = t
		} else if t.Unix() > max.Unix() {
			max = t
		}
	}

	if len(days) == 0 {

		row := CalendarRow{}

		for now.Weekday() != 0 {
			now = now.AddDate(0, 0, -1)
		}

		for i := 0; i < 7; i++ {
			vv := CalendarCell{}
			vv.Day.SetTime(now)
			row = append(row, vv)
			now = now.AddDate(0, 0, 1)
		}

		rows = append(rows, row)

	} else {

		b := min

		for b.Weekday() != 0 {
			b = b.AddDate(0, 0, -1)
		}

		row := CalendarRow{}

		for i := 0; i < 7; i++ {
			vv, ok := days[b]
			if !ok {
				vv.Day.SetTime(b)
			}
			row = append(row, vv)
			delete(days, b)
			b = b.AddDate(0, 0, 1)
		}

		rows = append(rows, row)

		for len(days) > 0 && b.Unix() <= max.Unix() {

			year := b.Year()
			month := b.Month()

			row := CalendarRow{}

			for i := 0; i < int(b.Weekday()); i++ {
				row = append(row, CalendarCell{})
			}

			for i := int(b.Weekday()); i < 7; i++ {

				if b.Year() == year && b.Month() == month {
					vv, ok := days[b]
					if !ok {
						vv.Day.SetTime(b)
					}
					row = append(row, vv)
					delete(days, b)
					b = b.AddDate(0, 0, 1)

				} else {
					row = append(row, CalendarCell{})
				}

			}

			rows = append(rows, row)
		}

	}

	task.Result.Calendar = rows

	return nil
}
