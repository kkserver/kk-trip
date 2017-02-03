package car

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TriggerCarRemovingTask struct {
	app.Task
	Car *Car `json:"car,omitempty"`
}

func (task *TriggerCarRemovingTask) GetResult() interface{} {
	return nil
}
