package uikit

import "github.com/maxence-charriere/go-app/v9/pkg/app"

const (
	HiddenClass             = "uk-hidden"
	InvisibleClass          = "uk-invisible"
	PanelClass              = "uk-panel"
	FloatLeftClass          = "uk-float-left"
	FloatRightClass         = "uk-float-right"
	ClearfixClass           = "uk-clearfix"
	OverflowHiddenClass     = "uk-overflow-hidden"
	OverflowAutoClass       = "uk-overflow-auto"
	ResizeClass             = "uk-resize"
	ResizeVerticalClass     = "uk-resize-vertical"
	DisplayBlockClass       = "uk-display-block"
	DisplayInlineClass      = "uk-display-inline"
	DisplayInlineBlockClass = "uk-display-inline-block"
	InlineClass             = "uk-inline"
	BorderRoundedClass      = "uk-border-rounded"
	BorderCircleClass       = "uk-border-circle"
	BorderPillClass         = "uk-border-pill"
	BoxShadowSmallClass     = "uk-box-shadow-small"
	BoxShadowMediumClass    = "uk-box-shadow-medium"
	BoxShadowLargeClass     = "uk-box-shadow-large"
	BoxShadowXLargeClass    = "uk-box-shadow-xlarge"
	DisabledClass           = "uk-disabled"
)

func Div(elems ...app.UI) app.HTMLDiv {
	return app.Div().Body(elems...)
}

func Text(v interface{}) app.HTMLDiv {
	return app.Div().Text(v)
}

func Style(url string) app.HTMLLink {
	return app.Link().Rel("stylesheet").Href(url)
}

func Css(css string) app.HTMLStyle {
	return app.Style().Text(css)
}

func Script(url string) app.HTMLScript {
	return app.Script().Src(url)
}

func Js(js string) app.HTMLScript {
	return app.Script().Text(js)
}
