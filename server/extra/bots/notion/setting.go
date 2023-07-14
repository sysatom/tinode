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
	{Key: tokenSettingKey, Type: types.FormFieldText, Title: "Internal Integration Token"},
	{Key: importPageIdSettingKey, Type: types.FormFieldText, Title: "MindCache page id"},
})
