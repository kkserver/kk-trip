package stop

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TriggerStopRemovingTask struct {
	app.Task
	Stop *Stop `json:"stop,omitempty"`
}

func (task *TriggerStopRemovingTask) GetResult() interface{} {
	return nil
}
