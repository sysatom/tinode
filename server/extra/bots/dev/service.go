package dev

import (
	"github.com/emicklei/go-restful/v3"
)

const serviceVersion = "v1"

func example(_ *restful.Request, resp *restful.Response) {
	_ = resp.WriteAsJson(map[string]interface{}{
		"title": "example",
	})
}
