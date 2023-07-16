package types

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type UI struct {
	App    app.UI
	CSS    []app.UI
	JS     []app.HTMLScript
	Global map[string]interface{}
}
