package workflow

import (
	"embed"
	"github.com/tinode/chat/server/extra/bots"
	"net/http"
)

const serviceVersion = "v1"

//go:embed webapp/build
var dist embed.FS

func webapp(rw http.ResponseWriter, req *http.Request) {
	bots.ServeFile(rw, req, dist, "webapp/build")
}
