package line

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type LineTaskResult struct {
	app.Result
	Line *Line `json:"line,omitempty"`
}

type LineTask struct {
	app.Task
	Id     int64 `json:"id,string"`
	Result LineTaskResult
}

func (task *LineTask) GetResult() interface{} {
	return &task.Result
}

func (task *LineTask) GetInhertType() string {
	return "line"
}

func (task *LineTask) GetClientName() string {
	return "Line.Get"
}
