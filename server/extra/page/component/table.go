package component

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
)

type Table struct {
	app.Compo
	Page   model.Page
	Schema types.TableMsg
}

func (c *Table) Render() app.UI {
	var alert app.UI
	switch c.Page.State {
	case model.PageStateProcessedSuccess:
		alert = app.Div().Class("uk-alert-success").Body(
			app.P().Style("padding", "20px").Text(fmt.Sprintf("Table [%s] processed success, %s", c.Page.PageId, c.Page.UpdatedAt)))
	case model.PageStateProcessedFailed:
		alert = app.Div().Class("uk-alert-danger").Body(
			app.P().Style("padding", "20px").Text(fmt.Sprintf("Table [%s] processed failed, %s", c.Page.PageId, c.Page.UpdatedAt)))
	}

	return app.Div().Body(
		alert,
		app.H1().Class(".uk-heading-small").Text(c.Schema.Title),
		app.Table().Class("uk-table uk-table-striped").Body(
			app.THead().Body(
				app.Tr().Body(
					app.Range(c.Schema.Header).Slice(func(i int) app.UI {
						return app.Th().Text(c.Schema.Header[i])
					}),
				),
			),
			app.TBody().Body(
				app.Range(c.Schema.Row).Slice(func(i int) app.UI {
					return app.Tr().Body(
						app.Range(c.Schema.Row[i]).Slice(func(j int) app.UI {
							return app.Td().Text(c.Schema.Row[i][j])
						}),
					)
				}),
			),
		),
	)
}
