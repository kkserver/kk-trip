package suggest

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-trip/trip/line"
	"github.com/kkserver/kk-trip/trip/route"
	"github.com/kkserver/kk-trip/trip/ticket"
	"strings"
)

type SuggestService struct {
	app.Service

	Line *SuggestLineTask
}

func (S *SuggestService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *SuggestService) HandleSuggestLineTask(a ISuggestApp, task *SuggestLineTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = line.ERROR_LINE
		task.Result.Errmsg = err.Error()
		return nil
	}

	var lineIds = map[int64]bool{}

	var lines = []line.Line{}

	var sql = bytes.NewBuffer(nil)

	var args = []interface{}{}

	sql.WriteString("1")

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

	if task.Direction != "" {
		vs := strings.Split(task.Direction, ",")
		sql.WriteString(" AND direction IN (")
		for i, v := range vs {
			if i != 0 {
				sql.WriteString(",")
			}
			sql.WriteString("?")
			args = append(args, v)
		}
		sql.WriteString(")")
	}

	args = append(args, task.Limit)

	if len(lines) < task.Limit {

		if task.Uid != 0 {

			t := ticket.TicketQueryTask{}
			t.Uid = task.Uid
			t.PageSize = task.Limit
			t.OrderBy = "desc"

			app.Handle(a, &t)

			if t.Result.Tickets != nil {

				idx := 0
				ids := bytes.NewBuffer(nil)

				for _, ticket := range t.Result.Tickets {
					if idx != 0 {
						ids.WriteString(",")
					}
					ids.WriteString(fmt.Sprintf("%d", ticket.LineId))
					idx = idx + 1
				}

				rows, err := kk.DBQuery(db, a.GetLineTable(), a.GetPrefix(), " WHERE id IN ("+ids.String()+") AND "+sql.String()+" ORDER BY FIELD(id,"+ids.String()+") LIMIT ?", args...)

				if err == nil {

					v := line.Line{}
					scanner := kk.NewDBScaner(&v)

					for rows.Next() {
						err = scanner.Scan(rows)
						if err != nil {
							break
						}
						_, ok := lineIds[v.Id]
						if !ok {
							lines = append(lines, v)
							lineIds[v.Id] = true
						}

					}

					rows.Close()
				}
			}
		}
	}

	if len(lines) < task.Limit {

		t := route.RouteNearbyTask{}
		t.Phone = task.Phone
		t.Latitude = task.Latitude
		t.Longitude = task.Longitude
		t.Distance = task.Distance
		t.PageSize = task.Limit - len(lines)

		app.Handle(a, &t)

		if t.Result.Routes != nil {

			idx := 0
			ids := bytes.NewBuffer(nil)

			for _, route := range t.Result.Routes {
				if idx != 0 {
					ids.WriteString(",")
				}
				ids.WriteString(fmt.Sprintf("%d", route.Id))
				idx = idx + 1
			}

			rows, err := kk.DBQuery(db, a.GetLineTable(), a.GetPrefix(), " WHERE routeid IN ("+ids.String()+") AND "+sql.String()+" ORDER BY FIELD(routeid,"+ids.String()+"),id DESC LIMIT ?", args...)

			if err == nil {

				v := line.Line{}
				scanner := kk.NewDBScaner(&v)

				for rows.Next() {
					err = scanner.Scan(rows)
					if err != nil {
						break
					}
					_, ok := lineIds[v.Id]
					if !ok {
						lines = append(lines, v)
						lineIds[v.Id] = true
					}

				}

				rows.Close()
			}
		}
	}

	if len(lines) < task.Limit {

		t := line.LineQueryTask{}
		t.Phone = task.Phone
		t.PageSize = task.Limit - len(lines)
		t.Status = "1"

		app.Handle(a, &t)

		if t.Result.Lines != nil {

			for _, v := range t.Result.Lines {
				_, ok := lineIds[v.Id]
				if !ok {
					lines = append(lines, v)
					lineIds[v.Id] = true
				}
			}
		}
	}

	task.Result.Lines = lines

	return nil
}
