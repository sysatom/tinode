package types

import (
	"github.com/looplab/fsm"
	"github.com/tinode/chat/server/extra/store/model"
)

type JobInfo struct {
	Job *model.Job
	FSM *fsm.FSM
}

type StepInfo struct {
	Step *model.Step
	FSM  *fsm.FSM
}
