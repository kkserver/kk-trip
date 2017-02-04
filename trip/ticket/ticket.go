package ticket

import (
	"database/sql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
)

const TicketStatusNone = 0      // 未支付
const TicketStatusPay = 200     //已支付
const TicketStatusCancel = 300  //已取消
const TicketStatusTimeout = 400 //已超时
const TicketStatusRefund = 500  //已退款

/**
 * 车票
 */
type Ticket struct {
	Id            int64  `json:"id,string"`
	OrderId       int64  `json:"orderId,string"`    //订单ID
	ScheduleId    int64  `json:"scheduleId,string"` //排期Id
	LineId        int64  `json:"lineId,string"`     //车次ID
	Uid           int64  `json:"uid,string"`        //用户ID
	Date          int64  `json:"date"`              //日期
	Status        int    `json:"status"`            //状态
	SeatNo        string `json:"seatNo"`            //座位
	PayValue      int64  `json:"payValue"`          //支付金额
	RefundValue   int64  `json:"refundValue"`       //退款金额
	Value         int64  `json:"value"`             //金额
	RefundType    string `json:"refundType"`        //退款类型
	RefundTradeNo string `json:"refundTradeNo"`     //退款订单号
	Ctime         int64  `json:"ctime"`
}

type TicketItem struct {
	LineId int64  `json:"lineId"`
	Date   int64  `json:"date"`
	SeatNo string `json:"seatNo"`
}

type TicketValue struct {
	PayValue    int64 `json:"payValue"`    //支付金额
	RefundValue int64 `json:"refundValue"` //退款金额
	Value       int64 `json:"value"`       //金额
}

type ITicketApp interface {
	app.IApp
	GetDB() (*sql.DB, error)
	GetPrefix() string
	GetTicketTable() *kk.DBTable
	GetScheduleTable() *kk.DBTable
	GetLineTable() *kk.DBTable
}
