package ticket

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TicketPreCreateTaskResult struct {
	app.Result
	Value    int64          `json:"value,omitempty"`    //订单金额
	PayValue int64          `json:"payValue,omitempty"` //支付金额
	Values   []*TicketValue `json:"values,omitempty"`
}

type TicketPreCreateTask struct {
	app.Task
	Uid      int64        `json:"uid,string"`
	Text     string       `json:"text"` //车票 lineId:date:seatno;lineId:date:seatno;
	Items    []TicketItem `json:"items"`
	CouponId int64        `json:"couponId"`
	Result   TicketPreCreateTaskResult
}

func (task *TicketPreCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *TicketPreCreateTask) GetInhertType() string {
	return "ticket"
}

func (task *TicketPreCreateTask) GetClientName() string {
	return "Ticket.PreCreate"
}
