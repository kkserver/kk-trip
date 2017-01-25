package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TriggerRouteRemovingTask struct {
	app.Task
	Route *Route `json:"route,omitempty"`
}

func (task *TriggerRouteRemovingTask) GetResult() interface{} {
	return nil
}
