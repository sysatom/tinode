package component

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/extra/types"
)

type Form struct {
	app.Compo
	FormId string
	Uid    string
	Topic  string
	Schema types.FormMsg
}

func (c *Form) Render() app.UI {
	var fields []app.UI

	// hidden
	fields = append(fields, app.Input().Hidden(true).Type("text").Name("x-csrf-token").Value(types.Id()))
	fields = append(fields, app.Input().Hidden(true).Type("text").Name("x-form_id").Value(c.FormId))
	fields = append(fields, app.Input().Hidden(true).Type("text").Name("x-uid").Value(c.Uid))
	fields = append(fields, app.Input().Hidden(true).Type("text").Name("x-topic").Value(c.Topic))

	for _, field := range c.Schema.Field {
		switch field.Type {
		case types.FormFieldText, types.FormFieldPassword, types.FormFieldNumber, types.FormFieldColor,
			types.FormFieldFile, types.FormFieldMonth, types.FormFieldDate, types.FormFieldTime, types.FormFieldEmail,
			types.FormFieldUrl, types.FormFieldRange:
			// input
			fields = append(fields, app.Div().Class("uk-margin").Body(
				app.Label().Class("uk-form-label").Text(field.Label),
				app.Div().Class("uk-form-controls").Body(
					app.Input().
						Class("uk-input").
						Type(string(field.Type)).
						Name(field.Key).
						Placeholder(field.Placeholder).
						Required(field.Required),
				),
			))
		case types.FormFieldRadio, types.FormFieldCheckbox:
			var options []app.UI
			for _, option := range field.Option {
				options = append(options, app.Label().Body(
					app.Input().Class(fmt.Sprintf("uk-%s", field.Type)).
						Type(string(field.Type)).
						Name(field.Key).
						Value(option),
					app.Text(option)),
				)
			}
			fields = append(fields, app.Div().Class("uk-margin").Body(
				app.Label().Class("uk-form-label").Text(field.Label),
				app.Div().Class("uk-form-controls").Body(options...),
			))
		case types.FormFieldTextarea:
			// textarea
			fields = append(fields, app.Div().Class("uk-margin").Body(
				app.Label().Class("uk-form-label").Text(field.Label),
				app.Div().Class("uk-form-controls").Body(
					app.Textarea().
						Class("uk-textarea").
						Name(field.Key).
						Placeholder(field.Placeholder).
						Required(field.Required),
				),
			))
		case types.FormFieldSelect:
			// select
			var options []app.UI
			for _, option := range field.Option {
				options = append(options, app.Option().Value(option).Text(option))
			}
			fields = append(fields, app.Div().Class("uk-margin").Body(
				app.Label().Class("uk-form-label").Text(field.Label),
				app.Div().Class("uk-form-controls").Body(
					app.Select().Class("uk-select").
						Name(field.Key).
						Required(field.Required).Body(options...),
				),
			))
		}
	}

	// button
	fields = append(fields, app.Div().Class("uk-margin").Body(
		app.Button().Class("uk-button uk-button-primary").Text("Submit").Type("submit"),
		app.Button().Class("uk-button uk-button-default").Text("Cancel"),
	))

	return app.Div().Body(
		app.H1().Class(".uk-heading-small").Text(c.Schema.Title),
		app.Form().Class("uk-form-stacked").Method("POST").Action("/extra/form").
			Body(fields...),
	)
}
