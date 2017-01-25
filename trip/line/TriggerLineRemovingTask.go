package line

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TriggerLineRemovingTask struct {
	app.Task
	Line *Line `json:"line,omitempty"`
}

func (task *TriggerLineRemovingTask) GetResult() interface{} {
	return nil
}
