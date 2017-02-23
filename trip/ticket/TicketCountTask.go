package ticket

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TicketCountTaskResult struct {
	app.Result
	Count int `json:"count"`
}

type TicketCountTask struct {
	app.Task
	Uid             int64  `json:"uid"`     //用户ID
	OrderId         int64  `json:"orderId"` //订单ID
	LineId          int64  `json:"lineId"`  //班次ID
	Status          string `json:"status"`
	StartDate       int64  `json:"startDate"`
	EndDate         int64  `json:"endDate"`
	StartRefundTime int64  `json:"startRefundTime"`
	EndRefundTime   int64  `json:"endRefundTime"`
	Result          TicketCountTaskResult
}

func (task *TicketCountTask) GetResult() interface{} {
	return &task.Result
}

func (task *TicketCountTask) GetInhertType() string {
	return "ticket"
}

func (task *TicketCountTask) GetClientName() string {
	return "Ticket.Count"
}
