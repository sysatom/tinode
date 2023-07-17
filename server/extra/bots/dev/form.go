package dev

import (
	"fmt"
	"github.com/tinode/chat/server/extra/ruleset/form"
	"github.com/tinode/chat/server/extra/types"
)

const (
	devFormID = "dev_form"
)

var formRules = []form.Rule{
	{
		Id:    devFormID,
		Title: "Current Value: 1, add/reduce ?",
		Field: []types.FormField{
			{
				Key:         "text",
				Type:        types.FormFieldText,
				ValueType:   types.FormFieldValueString,
				Value:       "",
				Label:       "Text",
				Placeholder: "Input text",
			},
			{
				Key:         "password",
				Type:        types.FormFieldPassword,
				ValueType:   types.FormFieldValueString,
				Value:       "",
				Label:       "Password",
				Placeholder: "Input password",
			},
			{
				Key:         "number",
				Type:        types.FormFieldNumber,
				ValueType:   types.FormFieldValueInt64,
				Value:       "",
				Label:       "Number",
				Placeholder: "Input number",
			},
			{
				Key:         "bool",
				Type:        types.FormFieldRadio,
				ValueType:   types.FormFieldValueBool,
				Value:       "",
				Label:       "Bool",
				Placeholder: "Switch",
				Option:      []string{"true", "false"},
			},
			{
				Key:         "multi",
				Type:        types.FormFieldCheckbox,
				ValueType:   types.FormFieldValueStringSlice,
				Value:       "",
				Label:       "Multiple",
				Placeholder: "Select multiple",
				Option:      []string{"a", "b", "c"},
			},
			{
				Key:         "textarea",
				Type:        types.FormFieldTextarea,
				ValueType:   types.FormFieldValueString,
				Value:       "",
				Label:       "Textarea",
				Placeholder: "Input textarea",
			},
			{
				Key:         "select",
				Type:        types.FormFieldSelect,
				ValueType:   types.FormFieldValueFloat64,
				Value:       "",
				Label:       "Select",
				Placeholder: "Select float",
				Option:      []string{"1.01", "2.02", "3.03"},
			},
			{
				Key:         "range",
				Type:        types.FormFieldRange,
				ValueType:   types.FormFieldValueInt64,
				Value:       "",
				Label:       "Range",
				Placeholder: "range value",
			},
		},
		Handler: func(ctx types.Context, values map[string]interface{}) types.MsgPayload {
			fmt.Println(values)
			return types.TextMsg{Text: fmt.Sprintf("ok, form [%s]", ctx.FormId)}
		},
	},
}
