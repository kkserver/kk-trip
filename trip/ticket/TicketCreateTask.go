package ticket

import (
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-order/order"
)

type TicketCreateTaskResult struct {
	app.Result
	Order   *order.Order `json:"order,omitempty"`
	Tickets []Ticket     `json:"tickets,omitempty"`
}

type TicketCreateItem struct {
	LineId int64  `json:"lineId"`
	Date   int64  `json:"date"`
	SeatNo string `json:"seatNo"`
}

type TicketCreateTask struct {
	app.Task
	Uid     int64        `json:"uid,string"`
	Text    string       `json:"text"` //车票 lineId:date:seatno;lineId:date:seatno;
	Items   []TicketItem `json:"items"`
	Expires int64        `json:"expires"`
	Result  TicketCreateTaskResult
}

func (task *TicketCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *TicketCreateTask) GetInhertType() string {
	return "ticket"
}

func (task *TicketCreateTask) GetClientName() string {
	return "Ticket.Create"
}
