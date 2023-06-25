package component

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
)

type Okr struct {
	app.Compo
	Page   model.Page
	Schema types.OkrMsg
}

func (c *Okr) Render() app.UI {
	var alert app.UI
	switch c.Page.State {
	case model.PageStateProcessedSuccess:
		alert = app.Div().Class("uk-alert-success").Body(
			app.P().Style("padding", "20px").Text(fmt.Sprintf("Okr [%s] processed success, %s", c.Page.PageID, c.Page.UpdatedAt)))
	case model.PageStateProcessedFailed:
		alert = app.Div().Class("uk-alert-danger").Body(
			app.P().Style("padding", "20px").Text(fmt.Sprintf("Okr [%s] processed failed, %s", c.Page.PageID, c.Page.UpdatedAt)))
	}

	ratio := float64(0)
	if c.Schema.Objective.TotalValue != 0 {
		ratio = float64(c.Schema.Objective.CurrentValue) / float64(c.Schema.Objective.TotalValue) * 100
	}

	return app.Div().Body(
		app.Raw(`
<style>
.okr-title {
    font-size: 1rem;
    font-weight: 500;
    margin-top: 10px;
    margin-bottom: 10px;
}

.okr-objective {
    font-size: 1rem;
    font-weight: 500;
}

.okr-item-title {
    font-size: 0.9rem;
    font-weight: 300;
    margin-top: 10px;
    margin-bottom: 10px;
}

.okr-progress {
    display: flex;
    flex-direction: row;
    padding: 5px;
    background-color: #fff;
    align-items: center;
    border-radius: 5px;
    width: 300px;
}

.okr-progress .ratio {
    font-size: 0.7rem;
    font-weight: 300;
    width: 20%;
    text-align: center;
}

.okr-memo {
    font-size: 0.9rem;
    font-weight: 300;
    padding: 5px;
    background-color: #fff;
    border-radius: 5px;
}

.okr-motive {
    font-size: 0.9rem;
    font-weight: 300;
    padding: 5px;
    background-color: #fff;
    border-radius: 5px;
}

.okr-feasibility {
    font-size: 0.9rem;
    font-weight: 300;
    padding: 5px;
    background-color: #fff;
    border-radius: 5px;
}

.okr-keyresult {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-between;
    padding: 10px;
    background-color: #fff;
    border-radius: 5px;
}

.okr-keyresult-item {
    width: 40%;
    background-color: #7232dd;
    color: #fff;
    padding: 10px;
    border-radius: 5px;
    line-height: 30px;
    margin-bottom: 10px;
}

.okr-keyresult-item .title {
    font-size: 0.9rem;
    font-weight: 400;
}

.okr-keyresult-item .value {
    font-size: 0.8rem;
    font-weight: 300;
}

.progress-bg {
    margin: 0 auto;
    width: 100%;
    height: 10px;
    border-radius: 4px;
    background-color: rgb(222, 228, 247);
    display: flex;
    align-items: center;
}
.progress-bg-line {
    background-color: rgb(222, 228, 247);
}

.progress-inner {
    height: 100%;
    border-radius: 4px;
    transition: all 0.5s cubic-bezier(0, 0.64, 0.36, 1);
}

.progress-line {
    background-color: blue;
}
</style>
`),
		alert,
		app.Div().Class("okr").Body(
			app.Div().Class("okr-title").Text(c.Schema.Title),
			app.Div().Body(
				app.Div().Text(c.Schema.Objective.Title),
				app.Div().Class("okr-item-title").Text("Progress"),
				app.Div().Class("okr-progress").Body(
					app.Div().Class("progress-bg-line progress-bg").Body(
						app.Div().Class("progress-line progress-inner").Style("width", fmt.Sprintf("%d%%", ratio)),
					),
					app.Div().Class("ratio").Text(ratio),
				),

				app.Div().Class("okr-item-title").Text("Memo"),
				app.Div().Class("okr-memo").Text(c.Schema.Objective.Memo),

				app.Div().Class("okr-item-title").Text("Motive"),
				app.Div().Class("okr-motive").Text(c.Schema.Objective.Motive),

				app.Div().Class("okr-item-title").Text("Feasibility"),
				app.Div().Class("okr-feasibility").Text(c.Schema.Objective.Feasibility),

				app.Div().Class("okr-item-title").Text("KeyResult"),
				app.Div().Class("okr-keyresult").Body(
					app.Range(c.Schema.KeyResult).Slice(func(i int) app.UI {
						return app.Div().Class("okr-keyresult-item").Body(
							app.Div().Class("title").Text(fmt.Sprintf("#%d %s", c.Schema.KeyResult[i].Sequence, c.Schema.KeyResult[i].Title)),
							app.Div().Class("value").Text(fmt.Sprintf("%d -> %d", c.Schema.KeyResult[i].CurrentValue, c.Schema.KeyResult[i].TargetValue)),
						)
					}),
				),
			),
		),
	)
}
