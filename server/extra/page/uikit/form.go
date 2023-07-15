package uikit

import "github.com/maxence-charriere/go-app/v9/pkg/app"

func Form(elems ...app.UI) app.HTMLForm {
	return app.Form().Body(elems...)
}

func Fieldset(elems ...app.UI) app.HTMLFieldSet {
	return app.FieldSet().Class("uk-fieldset").Body(elems...)
}

func Select(elems ...app.UI) app.HTMLSelect {
	return app.Select().Class("uk-select").Body(elems...)
}

func Textarea(elems ...app.UI) app.HTMLTextarea {
	return app.Textarea().Class("uk-textarea").Body(elems...)
}

func Radio() app.HTMLInput {
	return app.Input().Class("uk-radio").Type("radio")
}

func Checkbox() app.HTMLInput {
	return app.Input().Class("uk-checkbox").Type("checkbox")
}

func Range() app.HTMLInput {
	return app.Input().Class("uk-range").Type("range")
}
