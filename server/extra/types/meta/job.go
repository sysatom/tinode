package meta

import (
	"github.com/looplab/fsm"
	"github.com/tinode/chat/server/extra/store/model"
)

type JobInfo struct {
	Job *model.Job
	FSM *fsm.FSM
}
