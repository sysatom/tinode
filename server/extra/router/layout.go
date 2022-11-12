package router

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/extra/router/component"
	"github.com/tinode/chat/server/extra/store/model"
)

const layout = `
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

func renderForm(page model.Page) app.UI {
	return &component.Form{}
}

func renderChart(page model.Page) app.UI {
	return nil
}

func renderTable(page model.Page) app.UI {
	return nil
}
