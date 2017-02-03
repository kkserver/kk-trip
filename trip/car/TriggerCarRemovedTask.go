package car

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TriggerCarRemovedTask struct {
	app.Task
	Car *Car `json:"car,omitempty"`
}

func (task *TriggerCarRemovedTask) GetResult() interface{} {
	return nil
}
