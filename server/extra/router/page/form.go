package page

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Form struct {
	app.Compo
}

func (h *Form) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Form"),
		app.Form().Method("POST").Action("").
			Body(
				app.Input().
					Class("uk-input").
					Type("text"),
				app.Select().Class("uk-select").Name("D").Body(
					app.Option().Text("F"),
					app.Option().Text("M"),
				),
				app.Textarea().Class("uk-textarea"),
				app.Input().
					Class("uk-radio").
					Type("radio"),
				app.Input().
					Class("uk-checkbox").
					Type("checkbox"),
				app.Input().
					Class("uk-range").
					Type("range"),

				app.Button().Class("uk-button uk-button-primary").Text("Submit").Type("submit"),
				app.Button().Class("uk-button uk-button-default").Text("Cancel"),
			),
	)
}
