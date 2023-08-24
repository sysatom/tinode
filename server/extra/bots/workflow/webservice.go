package workflow

import (
	"github.com/tinode/chat/server/extra/pkg/route"
	"github.com/tinode/chat/server/extra/ruleset/webservice"
)

var webserviceRules = []webservice.Rule{
	{
		Method:        "GET",
		Path:          "/actions",
		Function:      actions,
		Documentation: "get bot actions",
		Option:        []route.Option{},
	},
}
