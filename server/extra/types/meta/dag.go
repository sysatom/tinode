package meta

import (
	"github.com/tinode/chat/server/extra/store/model"
	"time"
)

type Step struct {
	Name              string      `json:"name,omitempty"`
	UID               string      `json:"uid,omitempty"`
	WorkerUID         string      `json:"worker_uid,omitempty"`
	ResourceVersion   string      `json:"resource_version,omitempty"`
	Generation        int         `json:"generation,omitempty"`
	Finalizers        interface{} `json:"finalizers,omitempty"`
	DeletionTimestamp *time.Time  `json:"deletion_timestamp,omitempty"`

	JobId  int32           `json:"job_id,omitempty"`
	NodeId string          `json:"node_id,omitempty"`
	Depend []string        `json:"depend,omitempty"`
	State  model.StepState `json:"state,omitempty"`
}
