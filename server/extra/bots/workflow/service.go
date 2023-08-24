package workflow

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/ruleset/workflow"
	"github.com/tinode/chat/server/extra/types"
)

const serviceVersion = "v1"

func actions(_ *restful.Request, resp *restful.Response) {
	var result []struct {
		Bot          string
		Id           string
		Title        string
		Desc         string
		InputSchema  []types.FormField
		OutputSchema []types.FormField
	}
	for name, handler := range bots.List() {
		for _, item := range handler.Rules() {
			switch v := item.(type) {
			case []workflow.Rule:
				for _, rule := range v {
					result = append(result, struct {
						Bot          string
						Id           string
						Title        string
						Desc         string
						InputSchema  []types.FormField
						OutputSchema []types.FormField
					}{
						Bot:          name,
						Id:           rule.Id,
						Title:        rule.Title,
						Desc:         rule.Desc,
						InputSchema:  rule.InputSchema,
						OutputSchema: rule.OutputSchema,
					})
				}
			}
		}
	}

	_ = resp.WriteAsJson(result)
}
