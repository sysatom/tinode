package page

import (
	"encoding/json"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/extra/page/component"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
)

const Layout = `
<!DOCTYPE html>
<html>
    <head>
        <title>Page</title>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
     	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/uikit@3.15.12/dist/css/uikit.min.css" />
		<script src="https://cdn.jsdelivr.net/npm/uikit@3.15.12/dist/js/uikit.min.js"></script>
		<script src="https://cdn.jsdelivr.net/npm/uikit@3.15.12/dist/js/uikit-icons.min.js"></script>
    </head>

    <body>
        <div id="app" style="padding: 20px">%s</div>
    </body>
</html>
`

func RenderForm(page model.Page) app.UI {
	d, err := json.Marshal(page.Schema)
	if err != nil {
		return nil
	}
	var formMsg types.FormMsg
	err = json.Unmarshal(d, &formMsg)
	if err != nil {
		return nil
	}

	comp := &component.Form{
		Page:   page,
		Schema: formMsg,
	}
	return comp
}

func RenderChart(page model.Page) app.UI {
	return nil
}

func RenderTable(page model.Page) app.UI {
	return nil
}
