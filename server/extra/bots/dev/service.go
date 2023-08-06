package dev

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"io"
)

const serviceVersion = "v1"

func example(req *restful.Request, resp *restful.Response) {
	fmt.Println(io.ReadAll(req.Request.Body))
	_ = resp.WriteAsJson(map[string]interface{}{
		"title": "example",
	})
}
