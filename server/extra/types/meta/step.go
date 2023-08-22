package meta

import (
	"github.com/looplab/fsm"
	"github.com/tinode/chat/server/extra/store/model"
)

type StepInfo struct {
	Step *model.Step
	FSM  *fsm.FSM
}
