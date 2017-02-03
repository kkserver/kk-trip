package driver

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TriggerDriverRemovedTask struct {
	app.Task
	Driver *Driver `json:"driver,omitempty"`
}

func (task *TriggerDriverRemovedTask) GetResult() interface{} {
	return nil
}
