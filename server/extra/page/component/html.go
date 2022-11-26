package component

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
)

type Html struct {
	app.Compo
	Page   model.Page
	Schema types.HtmlMsg
}

func (c *Html) Render() app.UI {
	return app.Raw(c.Schema.Raw)
}
