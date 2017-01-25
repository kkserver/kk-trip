package route

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TriggerRouteRemovedTask struct {
	app.Task
	Route *Route `json:"route,omitempty"`
}

func (task *TriggerRouteRemovedTask) GetResult() interface{} {
	return nil
}
