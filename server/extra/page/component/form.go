package component

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
)

type Form struct {
	app.Compo
	Page   model.Page
	Form   model.Form
	Schema types.FormMsg
}

func (c *Form) Render() app.UI {
	var fields []app.UI

	var alert app.UI
	switch c.Page.State {
	case model.PageStateProcessedSuccess:
		alert = app.Div().Class("uk-alert-success").Body(
			app.P().Style("padding", "20px").Text(fmt.Sprintf("Form [%s] processed success, %s", c.Page.PageID, c.Page.UpdatedAt)))
	case model.PageStateProcessedFailed:
		alert = app.Div().Class("uk-alert-danger").Body(
			app.P().Style("padding", "20px").Text(fmt.Sprintf("Form [%s] processed failed, %s", c.Page.PageID, c.Page.UpdatedAt)))
	}

	// hidden
	fields = append(fields, app.Input().Hidden(true).Type("text").Name("x-csrf-token").Value(types.Id()))
	fields = append(fields, app.Input().Hidden(true).Type("text").Name("x-form_id").Value(c.Page.PageID))
	fields = append(fields, app.Input().Hidden(true).Type("text").Name("x-uid").Value(c.Page.UID))
	fields = append(fields, app.Input().Hidden(true).Type("text").Name("x-topic").Value(c.Page.Topic))
	for _, field := range c.Schema.Field {
		if field.Hidden {
			field.Value = fixInt64Value(field.ValueType, field.Value)
			fields = append(fields, app.Input().Hidden(true).Type("text").Name(field.Key).Value(field.Value))
		}
	}

	// fields
	for _, field := range c.Schema.Field {
		if field.Hidden {
			continue
		}
		field.Value = fixInt64Value(field.ValueType, field.Value)
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
						Value(field.Value).
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
						Checked(option == field.Value).
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
						Text(field.Value).
						Required(field.Required),
				),
			))
		case types.FormFieldSelect:
			// select
			var options []app.UI
			for _, option := range field.Option {
				options = append(options, app.Option().Selected(option == field.Value).Value(option).Text(option))
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
	if c.Page.State == model.PageStateCreated {
		fields = append(fields, app.Div().Class("uk-margin").Body(
			app.Button().Class("uk-button uk-button-primary").Text("Submit").Type("submit"),
		))
	}

	// record value
	if c.Page.State == model.PageStateProcessedSuccess || c.Page.State == model.PageStateProcessedFailed {
		fields = append(fields, app.Div().Class("").Body(
			app.H3().Text("Submit values"),
			app.Pre().Text(c.Form.Values),
		))
	}

	return app.Div().Body(
		alert,
		app.H1().Class(".uk-heading-small").Text(c.Schema.Title),
		app.Form().Class("uk-form-stacked").Method("POST").Action("/extra/form").
			Body(fields...),
	)
}

func fixInt64Value(t types.FormFieldValueType, v interface{}) interface{} {
	if t == types.FormFieldValueInt64 {
		switch v := v.(type) {
		case float64:
			return int64(v)
		}
	}
	return v
}
