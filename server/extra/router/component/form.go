package component

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Form struct {
	app.Compo
}

func (h *Form) Render() app.UI {
	return app.Div().Body(
		app.H1().Class("uk-heading-small uk-heading-bullet").Text("Form"),
		app.Form().Class("uk-form-stacked").Method("POST").Action("/extra/form").
			Body(
				app.Div().Class("uk-margin").Body(
					app.Label().Class("uk-form-label").Text("title"),
					app.Div().Class("uk-form-controls").Body(
						app.Input().
							Class("uk-input").
							Type("text").
							Placeholder("Please input name").Required(true),
					),
				),
				app.Div().Class("uk-margin").Body(
					app.Select().Class("uk-select").Name("D").Body(
						app.Option().Text("F"),
						app.Option().Text("M"),
					),
				),
				app.Div().Class("uk-margin").Body(
					app.Textarea().Class("uk-textarea"),
				),
				app.Div().Class("uk-margin").Body(
					app.Input().
						Class("uk-radio").
						Type("radio"),
				),
				app.Div().Class("uk-margin").Body(
					app.Input().
						Class("uk-checkbox").
						Type("checkbox"),
				),
				app.Div().Class("uk-margin").Body(
					app.Input().
						Class("uk-range").
						Type("range"),
				),
				app.Div().Class("uk-margin").Body(
					app.Button().Class("uk-button uk-button-primary").Text("Submit").Type("submit"),
					app.Button().Class("uk-button uk-button-default").Text("Cancel"),
				),
			),
	)
}
