package uikit

import "github.com/maxence-charriere/go-app/v9/pkg/app"

func Grid(elems ...app.UI) app.UI {
	return app.Div().Attr("uk-grid", "").Body(elems...)
}
