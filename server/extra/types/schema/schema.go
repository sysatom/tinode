package schema

import (
	"github.com/tinode/chat/server/extra/types"
)

func Step(steps ...types.Step) []types.Step {
	return steps
}

func Action(id string) types.Step {
	return types.Step{
		Type: types.ActionStep,
		Flag: id,
	}
}

func Bot(name string) types.Bot {
	return types.Bot(name)
}

func Command(bot types.Bot, token ...string) types.Step {
	return types.Step{
		Type: types.CommandStep,
		Bot:  bot,
		Args: token,
	}
}

func Form(id string) types.Step {
	return types.Step{
		Type: types.FormStep,
		Flag: id,
	}
}

func Instruct(id string, args ...string) types.Step {
	return types.Step{
		Type: types.InstructStep,
		Flag: id,
		Args: args,
	}
}

func Session(id string, args ...string) types.Step {
	return types.Step{
		Type: types.SessionStep,
		Flag: id,
		Args: args,
	}
}

func CommandTrigger(define string) types.Trigger {
	return types.Trigger{
		Type:   types.TriggerCommandType,
		Define: define,
	}
}

func TriggerArg(flag interface{}) interface{} { // todo
	return nil
}
