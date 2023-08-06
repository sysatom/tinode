package webhook

import (
	"github.com/tinode/chat/server/extra/pkg/route"
	"github.com/tinode/chat/server/extra/ruleset/webservice"
)

var webserviceRules = []webservice.Rule{
	{
		Method:        "POST",
		Path:          "/webhook/{flag}",
		Function:      webhook,
		Documentation: "trigger webhook",
		Option: []route.Option{
			route.WithPathParam("flag", "flag param", "string"),
		},
	},
}
