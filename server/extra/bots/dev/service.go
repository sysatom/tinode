package dev

import (
	"embed"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/tinode/chat/server/extra/bots"
	"io"
)

const serviceVersion = "v1"

func example(req *restful.Request, resp *restful.Response) {
	fmt.Println(io.ReadAll(req.Request.Body))
	_ = resp.WriteAsJson(map[string]interface{}{
		"title": "example",
	})
}

//go:embed webapp/build
var dist embed.FS

func webapp(req *restful.Request, resp *restful.Response) {
	bots.ServeFile(req, resp, dist, "webapp/build")
}
