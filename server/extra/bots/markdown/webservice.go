package markdown

import (
	"github.com/tinode/chat/server/extra/pkg/route"
	"github.com/tinode/chat/server/extra/ruleset/webservice"
)

var webserviceRules = []webservice.Rule{
	{
		Method:        "GET",
		Path:          "/editor/{flag}",
		Function:      editor,
		Documentation: "get markdown editor",
		Option: []route.Option{
			route.WithPathParam("flag", "flag param", "string"),
		},
	},
	{
		Method:        "POST",
		Path:          "/markdown",
		Function:      saveMarkdown,
		Documentation: "create markdown page",
	},
}
