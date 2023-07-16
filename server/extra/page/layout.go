package page

import (
	"encoding/json"
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/extra/page/component"
	"github.com/tinode/chat/server/extra/page/library"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"html"
	"strings"
)

const Layout = `
<!DOCTYPE html>
<html>
    <head>
        <title>Page</title>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
     	<link rel="stylesheet" href="%s" />
		<script src="%s"></script>
		<script src="%s"></script>
		%s
    </head>
    <body>
        %s
		%s
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

func RenderMarkdown(page model.Page) app.UI {
	d, err := json.Marshal(page.Schema)
	if err != nil {
		return nil
	}
	var msg types.MarkdownMsg
	err = json.Unmarshal(d, &msg)
	if err != nil {
		return nil
	}

	comp := &component.Markdown{
		Page:   page,
		Schema: msg,
	}
	return comp
}

func Render(comp types.UI) string {
	stylesStr := strings.Builder{}
	for _, style := range comp.CSS {
		stylesStr.WriteString(app.HTMLString(style))
	}
	scriptsStr := strings.Builder{}
	for _, script := range comp.JS {
		scriptsStr.WriteString(html.UnescapeString(app.HTMLString(script)))
	}
	return fmt.Sprintf(Layout,
		library.UIKitCss, library.UIKitJs, library.UIKitIconJs,
		stylesStr.String(), app.HTMLString(comp.App), scriptsStr.String())
}
