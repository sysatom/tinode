package dev

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/tinode/chat/server/extra/types"
)

const serviceVersion = "v1"

func example(_ *restful.Request, resp *restful.Response) {
	_ = resp.WriteAsJson(types.KV{
		"title": "example",
	})
}
