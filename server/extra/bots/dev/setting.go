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
	{secretSettingKey, types.FormFieldText, "Key", ""},
	{numberSettingKey, types.FormFieldNumber, "Number", ""},
})
