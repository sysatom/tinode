package dev

import (
	_ "embed"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/extra/page/library"
	"github.com/tinode/chat/server/extra/page/uikit"
	"github.com/tinode/chat/server/extra/ruleset/page"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/types"
	"time"
)

const (
	devPageId = "dev"
)

//go:embed static/example.css
var exampleCss string

//go:embed static/example.js
var exampleJs string

var pageRules = []page.Rule{
	{
		Id: devPageId,
		UI: func(ctx types.Context, flag string) (*types.UI, error) {
			p, err := store.Chatbot.ParameterGet(flag)
			if err != nil {
				return nil, err
			}

			css := []app.UI{
				uikit.Style(library.GithubMarkdownCss),
				uikit.Css(exampleCss),
			}
			js := []app.HTMLScript{
				uikit.Script(library.VueJs),
				uikit.Script(library.AxiosJs),
				uikit.Js(exampleJs),
			}

			app := uikit.App(
				uikit.H1("{{ message }}").Class(uikit.TextCenterClass),
				uikit.Grid(
					uikit.Card("One", app.Div().Text("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.")),
					uikit.Card("Two", app.Div().Text("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.")),
				).Class(uikit.FlexClass, uikit.FlexCenterClass),
				uikit.Icon("home"),
				uikit.Div(
					uikit.Label("One"),
					uikit.Label("Two").Class(uikit.LabelSuccessClass),
					uikit.Label("Three").Class(uikit.LabelWarningClass),
					uikit.Label("Four").Class(uikit.LabelDangerClass),
				),
				uikit.Article("title", time.Now().Format(time.DateTime), uikit.Text("article......")),
				uikit.Image("https://images.unsplash.com/photo-1490822180406-880c226c150b?fit=crop&w=650&h=433&q=80"),
				uikit.DividerIcon(),
				uikit.ModalToggle("example_modal", "modal"),
				uikit.Modal("example_modal", "modal", uikit.Text("content......")),
				uikit.Placeholder("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
				uikit.Progress(10, 100),
				uikit.Button("click event").Attr("@click", "greet"),
				uikit.Table(
					uikit.THead(
						uikit.Tr(
							uikit.Th(uikit.Text("heading")),
							uikit.Th(uikit.Text("heading")),
							uikit.Th(uikit.Text("heading")),
						)),
					uikit.TBody(
						uikit.Tr(
							uikit.Td(uikit.Text("data")),
							uikit.Td(uikit.Text("data")),
							uikit.Td(uikit.Text("data")),
						),
					),
					uikit.TFoot(
						uikit.Tr(
							uikit.Td(uikit.Text("footer")),
							uikit.Td(uikit.Text("footer")),
							uikit.Td(uikit.Text("footer")),
						),
					),
				).Class(uikit.TableDividerClass, uikit.TableHoverClass),
				uikit.Countdown(p.ExpiredAt),
			)

			return &types.UI{
				App:    app,
				CSS:    css,
				JS:     js,
				Global: p.Params,
			}, nil
		},
	},
}
