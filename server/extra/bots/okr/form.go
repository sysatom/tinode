package okr

import (
	"fmt"
	"github.com/tinode/chat/server/extra/ruleset/form"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
)

const (
	CreateObjectiveFormID      = "create_objective"
	UpdateObjectiveFormID      = "update_objective"
	CreateKeyResultFormID      = "create_key_result"
	UpdateKeyResultFormID      = "Update_key_result"
	CreateKeyResultValueFormID = "create_key_result_value"
)

var formRules = []form.Rule{
	{
		Id: CreateObjectiveFormID,
		Handler: func(ctx types.Context, values map[string]interface{}) types.MsgPayload {
			var objective model.Objective
			for key, value := range values {
				switch key {
				case "title":
					objective.Title = value.(string)
				case "memo":
					objective.Memo = value.(string)
				case "motive":
					objective.Motive = value.(string)
				case "feasibility":
					objective.Feasibility = value.(string)
				}
			}

			_, err := store.Chatbot.CreateObjective(&objective)
			if err != nil {
				logs.Err.Println(err)
				return types.TextMsg{Text: fmt.Sprintf("failed, form [%s]", ctx.FormId)}
			}

			return types.TextMsg{Text: fmt.Sprintf("ok, form [%s]", ctx.FormId)}
		},
	},
	{
		Id: UpdateObjectiveFormID,
		Handler: func(ctx types.Context, values map[string]interface{}) types.MsgPayload {
			var objective model.Objective
			for key, value := range values {
				switch key {
				case "sequence":
					objective.Sequence = value.(int64)
				case "title":
					objective.Title = value.(string)
				case "memo":
					objective.Memo = value.(string)
				case "motive":
					objective.Motive = value.(string)
				case "feasibility":
					objective.Feasibility = value.(string)
				}
			}

			err := store.Chatbot.UpdateObjective(&objective)
			if err != nil {
				logs.Err.Println(err)
				return types.TextMsg{Text: fmt.Sprintf("failed, form [%s]", ctx.FormId)}
			}

			return types.TextMsg{Text: fmt.Sprintf("ok, form [%s]", ctx.FormId)}
		},
	},
	{
		Id: CreateKeyResultFormID,
		Handler: func(ctx types.Context, values map[string]interface{}) types.MsgPayload {
			objectiveSequence := int64(0)
			var keyResult model.KeyResult
			for key, value := range values {
				switch key {
				case "objective_sequence":
					objectiveSequence = value.(int64)
				case "title":
					keyResult.Title = value.(string)
				case "memo":
					keyResult.Memo = value.(string)
				case "initial_value":
					keyResult.InitialValue = int32(value.(int64))
				case "target_value":
					keyResult.TargetValue = int32(value.(int64))
				case "value_mode":
					keyResult.ValueMode = model.ValueModeType(value.(string))
				}
			}

			objective, err := store.Chatbot.GetObjectiveBySequence(1, objectiveSequence) // todo
			if err != nil {
				return nil
			}

			// check
			if keyResult.TargetValue <= 0 {
				return nil
			}
			if keyResult.ValueMode != model.ValueSumMode &&
				keyResult.ValueMode != model.ValueLastMode &&
				keyResult.ValueMode != model.ValueAvgMode &&
				keyResult.ValueMode != model.ValueMaxMode {
				return nil
			}

			// store
			if keyResult.InitialValue > 0 {
				keyResult.CurrentValue = keyResult.InitialValue
			}
			keyResult.ObjectiveId = objective.Id
			_, err = store.Chatbot.CreateKeyResult(&keyResult)
			if err != nil {
				return nil
			}

			// aggregate
			err = store.Chatbot.AggregateObjectiveValue(objective.Id)
			if err != nil {
				return nil
			}

			return types.TextMsg{Text: fmt.Sprintf("ok, form [%s]", ctx.FormId)}
		},
	},
	{
		Id: UpdateKeyResultFormID,
		Handler: func(ctx types.Context, values map[string]interface{}) types.MsgPayload {
			var keyResult model.KeyResult
			for key, value := range values {
				switch key {
				case "sequence":
					keyResult.Sequence = value.(int64)
				case "title":
					keyResult.Title = value.(string)
				case "memo":
					keyResult.Memo = value.(string)
				case "initial_value":
					keyResult.InitialValue = int32(value.(int64))
				case "target_value":
					keyResult.TargetValue = int32(value.(int64))
				case "value_mode":
					keyResult.ValueMode = model.ValueModeType(value.(string))
				}
			}

			// check
			if keyResult.TargetValue <= 0 {
				return nil
			}
			if keyResult.ValueMode != model.ValueSumMode &&
				keyResult.ValueMode != model.ValueLastMode &&
				keyResult.ValueMode != model.ValueAvgMode &&
				keyResult.ValueMode != model.ValueMaxMode {
				return nil
			}

			keyResult.UserId = 1 // todo
			err := store.Chatbot.UpdateKeyResult(&keyResult)
			if err != nil {
				return nil
			}

			// update value
			reply, err := store.Chatbot.GetKeyResultBySequence(1, keyResult.Sequence)
			if err != nil {
				return nil
			}
			err = store.Chatbot.AggregateKeyResultValue(reply.Id)
			if err != nil {
				return nil
			}

			return types.TextMsg{Text: fmt.Sprintf("ok, form [%s]", ctx.FormId)}
		},
	},
	{
		Id: CreateKeyResultValueFormID,
		Handler: func(ctx types.Context, values map[string]interface{}) types.MsgPayload {
			keyResultSequence := values["key_result_sequence"].(int64)
			value := int32(values["value"].(int64))

			keyResult, err := store.Chatbot.GetKeyResultBySequence(1, keyResultSequence) // todo
			if err != nil {
				return nil
			}
			_, err = store.Chatbot.CreateKeyResultValue(&model.KeyResultValue{Value: value, KeyResultId: keyResult.Id})
			if err != nil {
				return nil
			}
			err = store.Chatbot.AggregateKeyResultValue(keyResult.Id)
			if err != nil {
				return nil
			}
			err = store.Chatbot.AggregateObjectiveValue(keyResult.ObjectiveId)
			if err != nil {
				return nil
			}

			return types.TextMsg{Text: fmt.Sprintf("ok, form [%s]", ctx.FormId)}
		},
	},
}
