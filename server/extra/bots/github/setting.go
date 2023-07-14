package github

import (
	"github.com/tinode/chat/server/extra/ruleset/setting"
	"github.com/tinode/chat/server/extra/types"
)

const (
	repoSettingKey = "repo"
)

var settingRules = setting.Rule([]setting.Row{
	{Key: repoSettingKey, Type: types.FormFieldText, Title: "Repo"},
})
