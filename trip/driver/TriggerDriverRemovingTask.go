package driver

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TriggerDriverRemovingTask struct {
	app.Task
	Driver *Driver `json:"driver,omitempty"`
}

func (task *TriggerDriverRemovingTask) GetResult() interface{} {
	return nil
}
