package notion

import (
	"github.com/tinode/chat/server/extra/ruleset/setting"
	"github.com/tinode/chat/server/extra/types"
)

const (
	tokenSettingKey        = "token"
	importPageIdSettingKey = "import_page_id"
)

var settingRules = setting.Rule([]setting.Row{
	{tokenSettingKey, types.FormFieldText, "Internal Integration Token", ""},
	{importPageIdSettingKey, types.FormFieldText, "MindCache page id", ""},
})
