package utils

import (
	"embed"
	"github.com/emicklei/go-restful/v3"
	"io/fs"
	"net/http"
	"strings"
)

func ServeFile(req *restful.Request, resp *restful.Response, dist embed.FS, dir string) {
	s := fs.FS(dist)
	h, err := fs.Sub(s, dir)
	if err != nil {
		_ = resp.WriteError(http.StatusNotFound, err)
		return
	}

	subpath := req.PathParameter("subpath")
	if subpath == "" {
		subpath = "index.html"
	}

	if strings.HasSuffix(subpath, "html") {
		resp.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	if strings.HasSuffix(subpath, "css") {
		resp.ResponseWriter.Header().Set("Content-Type", "text/css; charset=utf-8")
	}
	if strings.HasSuffix(subpath, "js") {
		resp.ResponseWriter.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	}
	if strings.HasSuffix(subpath, "svg") {
		resp.ResponseWriter.Header().Set("Content-Type", "image/svg+xml")
	}

	content, err := fs.ReadFile(h, subpath)
	if err != nil {
		_ = resp.WriteError(http.StatusNotFound, err)
		return
	}

	_, _ = resp.Write(content)
}
