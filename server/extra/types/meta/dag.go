package meta

import (
	"github.com/tinode/chat/server/extra/store/model"
	"time"
)

type Step struct {
	Name              string
	UID               string
	WorkerUID         string
	ResourceVersion   string
	Generation        int
	Finalizers        interface{}
	DeletionTimestamp *time.Time

	DagUID       string
	NodeId       string
	DependNodeId []string
	State        model.StepState
}
