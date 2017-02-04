package ticket

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TicketRefundTaskResult struct {
	app.Result
	Ticket *Ticket `json:"ticket,omitempty"`
}

type TicketRefundTask struct {
	app.Task
	Id            int64  `json:"id"`
	RefundType    string `json:"refundType"`    //退款类型
	RefundTradeNo string `json:"refundTradeNo"` //退款订单号
	Result        TicketRefundTaskResult
}

func (task *TicketRefundTask) GetResult() interface{} {
	return &task.Result
}

func (task *TicketRefundTask) GetInhertType() string {
	return "ticket"
}

func (task *TicketRefundTask) GetClientName() string {
	return "Ticket.Refund"
}
