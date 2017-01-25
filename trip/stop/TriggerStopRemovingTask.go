package stop

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TriggerStopRemovedTask struct {
	app.Task
	Stop *Stop `json:"stop,omitempty"`
}

func (task *TriggerStopRemovedTask) GetResult() interface{} {
	return nil
}
