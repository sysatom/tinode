package clipboard

import "github.com/tinode/chat/server/extra/ruleset/instruct"

const (
	ShareInstruct = "clipboard_share"
)

var instructRules = []instruct.Rule{
	{
		Id:   ShareInstruct,
		Args: []string{"txt"},
	},
}
