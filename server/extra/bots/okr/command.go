package okr

import (
	"fmt"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/ruleset/command"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	"strconv"
)

var commandRules = []command.Rule{
	{
		Define: "version",
		Help:   `Version`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
		},
	},
	{
		Define: `obj list`,
		Help:   `List objectives`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			items, err := store.Chatbot.ListObjectives(ctx.AsUser, ctx.Original)
			if err != nil {
				logs.Err.Println(err)
				return nil
			}

			var header []string
			var row [][]interface{}
			if len(items) > 0 {
				header = []string{"Sequence", "Title", "Current Value", "Total Value"}
				for _, v := range items {
					row = append(row, []interface{}{strconv.Itoa(int(v.Sequence)), v.Title, strconv.Itoa(int(v.CurrentValue)), strconv.Itoa(int(v.TotalValue))})
				}
			}
			if len(row) == 0 {
				return types.TextMsg{Text: "Empty"}
			}

			return types.TableMsg{Title: "Objectives", Header: header, Row: row}
		},
	},
	{
		Define: `obj [number]`,
		Help:   `View objective`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			sequence, _ := tokens[1].Value.Int64()

			objective, err := store.Chatbot.GetObjectiveBySequence(ctx.AsUser, ctx.Original, sequence)
			if err != nil {
				logs.Err.Println(err)
				return nil
			}

			keyResult, err := store.Chatbot.ListKeyResultsByObjectiveId(objective.Id)
			if err != nil {
				logs.Err.Println(err)
				return nil
			}

			return bots.StoreOkr(ctx, types.OkrMsg{
				Title:     fmt.Sprintf("Objective #%d", objective.Sequence),
				Objective: objective,
				KeyResult: keyResult,
			})
		},
	},
	{
		Define: `obj del [number]`,
		Help:   `Delete objective`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			sequence, _ := tokens[2].Value.Int64()

			err := store.Chatbot.DeleteObjectiveBySequence(ctx.AsUser, ctx.Original, sequence)
			if err != nil {
				logs.Err.Println(err)
				return types.TextMsg{Text: "failed"}
			}

			return types.TextMsg{Text: "ok"}
		},
	},
	{
		Define: `obj update [number]`,
		Help:   `Update objective`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			sequence, _ := tokens[2].Value.Int64()

			item, err := store.Chatbot.GetObjectiveBySequence(ctx.AsUser, ctx.Original, sequence)
			if err != nil {
				return nil
			}

			return bots.StoreForm(ctx, types.FormMsg{
				ID:    UpdateObjectiveFormID,
				Title: fmt.Sprintf("Update Objective #%d", sequence),
				Field: []types.FormField{
					{
						Key:       "sequence",
						Type:      types.FormFieldNumber,
						ValueType: types.FormFieldValueInt64,
						Value:     item.Sequence,
						Label:     "Sequence",
					},
					{
						Key:       "title",
						Type:      types.FormFieldText,
						ValueType: types.FormFieldValueString,
						Value:     item.Title,
						Label:     "Title",
					},
					{
						Key:       "memo",
						Type:      types.FormFieldText,
						ValueType: types.FormFieldValueString,
						Value:     item.Memo,
						Label:     "Memo",
					},
					{
						Key:       "motive",
						Type:      types.FormFieldText,
						ValueType: types.FormFieldValueString,
						Value:     item.Motive,
						Label:     "Motive",
					},
					{
						Key:       "feasibility",
						Type:      types.FormFieldText,
						ValueType: types.FormFieldValueString,
						Value:     item.Feasibility,
						Label:     "Feasibility",
					},
				},
			})
		},
	},
	{
		Define: `obj create`,
		Help:   `Create Objective`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return bots.StoreForm(ctx, types.FormMsg{
				ID:    CreateObjectiveFormID,
				Title: "Create Objective",
				Field: []types.FormField{
					{
						Key:       "title",
						Type:      types.FormFieldText,
						ValueType: types.FormFieldValueString,
						Label:     "Title",
					},
					{
						Key:       "memo",
						Type:      types.FormFieldText,
						ValueType: types.FormFieldValueString,
						Label:     "Memo",
					},
					{
						Key:       "motive",
						Type:      types.FormFieldText,
						ValueType: types.FormFieldValueString,
						Label:     "Motive",
					},
					{
						Key:       "feasibility",
						Type:      types.FormFieldText,
						ValueType: types.FormFieldValueString,
						Label:     "Feasibility",
					},
					{
						Key:       "is_plan",
						Type:      types.FormFieldRadio,
						ValueType: types.FormFieldValueBool,
						Label:     "IsPlan",
						Option:    []string{"true", "false"},
					},
					{
						Key:       "plan_start",
						Type:      types.FormFieldText,
						ValueType: types.FormFieldValueString,
						Label:     "PlanStart",
					},
					{
						Key:       "plan_end",
						Type:      types.FormFieldText,
						ValueType: types.FormFieldValueString,
						Label:     "PlanEnd",
					},
				},
			})
		},
	},
	{
		Define: `kr list`,
		Help:   `List KeyResult`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			items, err := store.Chatbot.ListKeyResults(ctx.AsUser, ctx.Original)
			if err != nil {
				return nil
			}

			var header []string
			var row [][]interface{}
			if len(items) > 0 {
				header = []string{"Sequence", "Title", "Current Value", "Target Value"}
				for _, v := range items {
					row = append(row, []interface{}{strconv.Itoa(int(v.Sequence)), v.Title, strconv.Itoa(int(v.CurrentValue)), strconv.Itoa(int(v.TargetValue))})
				}
			}

			return types.TableMsg{
				Title:  "KeyResult",
				Header: header,
				Row:    row,
			}
		},
	},
	{
		Define: `kr create`,
		Help:   `Create KeyResult`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return bots.StoreForm(ctx, types.FormMsg{
				ID:    CreateKeyResultFormID,
				Title: "Create Key Result",
				Field: []types.FormField{
					{
						Key:       "objective_sequence",
						Type:      types.FormFieldNumber,
						ValueType: types.FormFieldValueInt64,
						Label:     "Objective Sequence",
					},
					{
						Key:       "title",
						Type:      types.FormFieldText,
						ValueType: types.FormFieldValueString,
						Label:     "Title",
					},
					{
						Key:       "memo",
						Type:      types.FormFieldText,
						ValueType: types.FormFieldValueString,
						Label:     "Memo",
					},
					{
						Key:       "initial_value",
						Type:      types.FormFieldNumber,
						ValueType: types.FormFieldValueInt64,
						Label:     "initial value",
					},
					{
						Key:       "target_value",
						Type:      types.FormFieldNumber,
						ValueType: types.FormFieldValueInt64,
						Label:     "target value",
					},
					{
						Key:       "value_mode",
						Type:      types.FormFieldSelect,
						ValueType: types.FormFieldValueString,
						Label:     "value mode",
						Option:    []string{"avg", "max", "sum", "last"},
					},
				},
			})
		},
	},
	{
		Define: `kr [number]`,
		Help:   `View KeyResult`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			sequence, _ := tokens[1].Value.Int64()

			item, err := store.Chatbot.GetKeyResultBySequence(ctx.AsUser, ctx.Original, sequence)
			if err != nil {
				return nil
			}

			return types.InfoMsg{
				Title: fmt.Sprintf("KeyResult #%d", sequence),
				Model: item,
			}
		},
	},
	{
		Define: `kr del [number]`,
		Help:   `Delete KeyResult`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			sequence, _ := tokens[2].Value.Int64()

			err := store.Chatbot.DeleteKeyResultBySequence(ctx.AsUser, ctx.Original, sequence)
			if err != nil {
				logs.Err.Println(err)
				return types.TextMsg{Text: "failed"}
			}

			return types.TextMsg{Text: "ok"}
		},
	},
	{
		Define: `kr update [number]`,
		Help:   `Update KeyResult`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			sequence, _ := tokens[2].Value.Int64()

			item, err := store.Chatbot.GetKeyResultBySequence(ctx.AsUser, ctx.Original, sequence)
			if err != nil {
				return nil
			}

			return bots.StoreForm(ctx, types.FormMsg{
				ID:    UpdateKeyResultFormID,
				Title: fmt.Sprintf("Update KeyResult #%d", sequence),
				Field: []types.FormField{
					{
						Key:       "sequence",
						Type:      types.FormFieldNumber,
						ValueType: types.FormFieldValueInt64,
						Label:     "Sequence",
						Value:     item.Sequence,
					},
					{
						Key:       "title",
						Type:      types.FormFieldText,
						ValueType: types.FormFieldValueString,
						Label:     "Title",
						Value:     item.Title,
					},
					{
						Key:       "memo",
						Type:      types.FormFieldText,
						ValueType: types.FormFieldValueString,
						Label:     "Memo",
						Value:     item.Memo,
					},
					{
						Key:       "target_value",
						Type:      types.FormFieldNumber,
						ValueType: types.FormFieldValueInt64,
						Label:     "target value",
						Value:     item.TargetValue,
					},
					{
						Key:       "value_mode",
						Type:      types.FormFieldSelect,
						ValueType: types.FormFieldValueString,
						Label:     "value mode",
						Option:    []string{"avg", "max", "sum", "last"},
						Value:     item.ValueMode,
					},
				},
			})
		},
	},
	{
		Define: `kr value`,
		Help:   `Create KeyResult value`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return bots.StoreForm(ctx, types.FormMsg{
				ID:    CreateKeyResultValueFormID,
				Title: "Create Key Result value",
				Field: []types.FormField{
					{
						Key:       "key_result_sequence",
						Type:      types.FormFieldNumber,
						ValueType: types.FormFieldValueInt64,
						Label:     "Key Result Sequence",
					},
					{
						Key:       "value",
						Type:      types.FormFieldNumber,
						ValueType: types.FormFieldValueInt64,
						Label:     "Value",
					},
				},
			})
		},
	},
	{
		Define: `kr value [number]`,
		Help:   `List KeyResult value`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			sequence, _ := tokens[2].Value.Int64()

			keyResult, err := store.Chatbot.GetKeyResultBySequence(ctx.AsUser, ctx.Original, sequence)
			if err != nil {
				return nil
			}

			items, err := store.Chatbot.GetKeyResultValues(keyResult.Id)
			if err != nil {
				return nil
			}

			var header []string
			var row [][]interface{}
			if len(items) > 0 {
				header = []string{"Value", "Datetime"}
				for _, v := range items {
					row = append(row, []interface{}{strconv.Itoa(int(v.Value)), v.CreatedAt})
				}
			}

			return types.TableMsg{
				Title:  fmt.Sprintf("KeyResult #%d Values", sequence),
				Header: header,
				Row:    row,
			}
		},
	},
	{
		Define: `todo list`,
		Help:   `List todo`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
		},
	},
	{
		Define: `todo create`,
		Help:   "Create Todo something",
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
		},
	},
	{
		Define: `todo update [number]`,
		Help:   "Update Todo something",
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
		},
	},
	{
		Define: `todo complete [number]`,
		Help:   "Complete Todo",
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
		},
	},
	{
		Define: `counters`,
		Help:   `List Counter`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
		},
	},
	{
		Define: "counter [string]",
		Help:   `Count things`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
		},
	},
	{
		Define: "increase [string]",
		Help:   `Increase Counter`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
		},
	},
	{
		Define: "decrease [string]",
		Help:   `Decrease Counter`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
		},
	},
	{
		Define: "reset [string]",
		Help:   `Reset Counter`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
		},
	},
	{
		Define: "tags",
		Help:   `List Tag`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
		},
	},
	{
		Define: "tag [string]",
		Help:   `Get Model tags`,
		Handler: func(ctx types.Context, tokens []*command.Token) types.MsgPayload {
			return types.TextMsg{Text: "V1"}
		},
	},
}
