package page

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/tinode/chat/server/extra/page"
	"github.com/tinode/chat/server/extra/types"
	"net/http"
)

type Rule struct {
	Id string
	UI func(ctx types.Context, flag string) (app.UI, error)
}

type Ruleset []Rule

func (r Ruleset) ProcessPage(ctx types.Context, flag string) (string, error) {
	for _, rule := range r {
		if rule.Id == ctx.PageRuleId {
			ui, err := rule.UI(ctx, flag)
			if err != nil {
				return "", err
			}
			return page.Render(ui), nil
		}
	}
	return "", fmt.Errorf("%d not found", http.StatusNotFound)
}
