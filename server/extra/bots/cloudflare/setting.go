package cloudflare

import (
	"github.com/tinode/chat/server/extra/ruleset/setting"
	"github.com/tinode/chat/server/extra/types"
)

const (
	tokenSettingKey     = "token"
	zoneIdSettingKey    = "zone_id"
	accountIdSettingKey = "account_id"
)

var settingRules = setting.Rule([]setting.Row{
	{tokenSettingKey, types.FormFieldText, "Token", ""},
	{zoneIdSettingKey, types.FormFieldText, "Zone Id", ""},
	{accountIdSettingKey, types.FormFieldText, "Account Id", ""},
})
