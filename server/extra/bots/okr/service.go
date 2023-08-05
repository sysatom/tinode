package okr

import (
	"embed"
	"github.com/emicklei/go-restful/v3"
	"github.com/tinode/chat/server/extra/bots"
)

const serviceVersion = "v1"

//go:embed webapp/build
var dist embed.FS

func webapp(req *restful.Request, resp *restful.Response) {
	bots.ServeFile(req, resp, dist, "webapp/build")
}
