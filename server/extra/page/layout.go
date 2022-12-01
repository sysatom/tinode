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

func RenderForm(page model.Page, form model.Form) app.UI {
	d, err := json.Marshal(page.Schema)
	if err != nil {
		return nil
	}
	var msg types.FormMsg
	err = json.Unmarshal(d, &msg)
	if err != nil {
		return nil
	}

	comp := &component.Form{
		Page:   page,
		Form:   form,
		Schema: msg,
	}
	return comp
}

func RenderTable(page model.Page) app.UI {
	d, err := json.Marshal(page.Schema)
	if err != nil {
		return nil
	}
	var msg types.TableMsg
	err = json.Unmarshal(d, &msg)
	if err != nil {
		return nil
	}

	comp := &component.Table{
		Page:   page,
		Schema: msg,
	}
	return comp
}

func RenderOkr(page model.Page) app.UI {
	d, err := json.Marshal(page.Schema)
	if err != nil {
		return nil
	}
	var msg types.OkrMsg
	err = json.Unmarshal(d, &msg)
	if err != nil {
		return nil
	}

	comp := &component.Okr{
		Page:   page,
		Schema: msg,
	}
	return comp
}

func RenderShare(page model.Page) app.UI {
	d, err := json.Marshal(page.Schema)
	if err != nil {
		return nil
	}
	var msg types.TextMsg
	err = json.Unmarshal(d, &msg)
	if err != nil {
		return nil
	}

	comp := &component.Share{
		Page:   page,
		Schema: msg,
	}
	return comp
}

func RenderJson(page model.Page) app.UI {
	d, err := json.Marshal(page.Schema)
	if err != nil {
		return nil
	}
	var msg types.TextMsg
	err = json.Unmarshal(d, &msg)
	if err != nil {
		return nil
	}

	comp := &component.Json{
		Page:   page,
		Schema: msg,
	}
	return comp
}

func RenderHtml(page model.Page) app.UI {
	d, err := json.Marshal(page.Schema)
	if err != nil {
		return nil
	}
	var msg types.HtmlMsg
	err = json.Unmarshal(d, &msg)
	if err != nil {
		return nil
	}

	comp := &component.Html{
		Page:   page,
		Schema: msg,
	}
	return comp
}
