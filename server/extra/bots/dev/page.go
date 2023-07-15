package dev

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/extra/page/uikit"
	"github.com/tinode/chat/server/extra/ruleset/page"
	"github.com/tinode/chat/server/extra/types"
)

const (
	devPageId = "dev"
)

var pageRules = []page.Rule{
	{
		Id: devPageId,
		UI: func(ctx types.Context, flag string) (app.UI, error) {
			return app.Div().Body(uikit.Grid(app.Div(), app.Div())), nil
		},
	},
}
