package dev

import (
	"github.com/tinode/chat/server/extra/ruleset/setting"
	"github.com/tinode/chat/server/extra/types"
)

const (
	secretSettingKey = "secret"
	numberSettingKey = "number"
)

var settingRules = setting.Rule([]setting.Row{
	{Key: secretSettingKey, Type: types.FormFieldText, Title: "Key"},
	{Key: numberSettingKey, Type: types.FormFieldNumber, Title: "Number"},
})
