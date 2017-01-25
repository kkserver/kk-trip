package line

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TriggerLineRemovedTask struct {
	app.Task
	Line *Line `json:"line,omitempty"`
}

func (task *TriggerLineRemovedTask) GetResult() interface{} {
	return nil
}
