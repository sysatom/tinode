package meta

import "github.com/tinode/chat/server/extra/store/model"

type Step struct {
	DagUID       string
	NodeId       string
	DependNodeId []string
	State        model.StepState
}
