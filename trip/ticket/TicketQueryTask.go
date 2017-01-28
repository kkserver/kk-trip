package ticket

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TicketQueryCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
}

type TicketQueryTaskResult struct {
	app.Result
	Counter *TicketQueryCounter `json:"counter,omitempty"`
	Tickets []Ticket            `json:"tickets,omitempty"`
}

type TicketQueryTask struct {
	app.Task
	Id        int64  `json:"id"`
	Uid       int64  `json:"uid"`     //用户ID
	OrderId   int64  `json:"orderId"` //订单ID
	Status    string `json:"status"`
	StartDate int64  `json:"startDate"`
	EndDate   int64  `json:"endDate"`
	OrderBy   string `json:"orderBy"` // desc, asc , date
	PageIndex int    `json:"p"`
	PageSize  int    `json:"size"`
	Counter   bool   `json:"counter"`
	Result    TicketQueryTaskResult
}

func (task *TicketQueryTask) GetResult() interface{} {
	return &task.Result
}

func (task *TicketQueryTask) GetInhertType() string {
	return "ticket"
}

func (task *TicketQueryTask) GetClientName() string {
	return "Ticket.Query"
}
