package component

import (
	"bytes"
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/utils"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Markdown struct {
	app.Compo
	Page   model.Page
	Schema types.MarkdownMsg
}

func (c *Markdown) Render() app.UI {
	var buf bytes.Buffer
	source := utils.StringToBytes(c.Schema.Raw)

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	err := md.Convert(source, &buf)
	if err != nil {
		buf.WriteString("error markdown")
	}

	return app.Raw(fmt.Sprintf(`
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.2.0/github-markdown.min.css" />
<div class="markdown-body">%s</div>`, buf.String()))
}
