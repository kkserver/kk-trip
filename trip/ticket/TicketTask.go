package ticket

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TicketTaskResult struct {
	app.Result
	Ticket *Ticket `json:"ticket,omitempty"`
}

type TicketTask struct {
	app.Task
	Id     int64       `json:"id"`
	Uid    interface{} `json:"uid"` //用户ID
	Result TicketTaskResult
}

func (task *TicketTask) GetResult() interface{} {
	return &task.Result
}

func (task *TicketTask) GetInhertType() string {
	return "ticket"
}

func (task *TicketTask) GetClientName() string {
	return "Ticket.Get"
}
