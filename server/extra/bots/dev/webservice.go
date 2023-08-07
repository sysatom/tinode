package dev

import (
	"github.com/tinode/chat/server/extra/pkg/route"
	"github.com/tinode/chat/server/extra/ruleset/webservice"
	"github.com/tinode/chat/server/extra/store/model"
)

var webserviceRules = []webservice.Rule{
	{
		Method:        "GET",
		Path:          "/example",
		Function:      example,
		Documentation: "get example data",
		Option: []route.Option{
			route.WithReturns(model.Message{}),
			route.WithWrites(model.Message{}),
		},
	},
}
