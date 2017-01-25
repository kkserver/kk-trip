package line

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type LineRemoveTaskResult struct {
	app.Result
	Line *Line `json:"line,omitempty"`
}

type LineRemoveTask struct {
	app.Task
	Id     int64 `json:"id,string"`
	Result LineRemoveTaskResult
}

func (task *LineRemoveTask) GetResult() interface{} {
	return &task.Result
}

func (task *LineRemoveTask) GetInhertType() string {
	return "line"
}

func (task *LineRemoveTask) GetClientName() string {
	return "Line.Remove"
}
