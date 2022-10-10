package router

import (
	"fmt"
	"github.com/tinode/chat/server/extra/vendors"
	"net/http"
	"regexp"
	"strings"
)

var rOAuth = regexp.MustCompile(`/oauth/\w+`)
var rOAuthRedirect = regexp.MustCompile(`/oauth/\w+/redirect`)

func ServeExtra(rw http.ResponseWriter, req *http.Request) {
	switch {
	case rOAuth.MatchString(req.URL.Path):
		oauth(rw, req)
	case rOAuthRedirect.MatchString(req.URL.Path):
		oauthRedirect(rw, req)
	default:
		rw.Write([]byte("Unknown Pattern"))
	}
}

func oauthRedirect(rw http.ResponseWriter, req *http.Request) {
	category := strings.ReplaceAll(req.URL.Path, "/extra/oauth/", "")
	category = strings.ReplaceAll(req.URL.Path, "/redirect", "")
	category = strings.ToLower(category)
	provider := vendors.NewOAuthProvider(category, "")
	url, err := provider.Redirect(req)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("oauth redirect error"))
		return
	}
	rw.Header().Set("Location", url)
	rw.WriteHeader(http.StatusFound)
}

func oauth(rw http.ResponseWriter, req *http.Request) {
	category := strings.ToLower(strings.ReplaceAll(req.URL.Path, "/extra/oauth/", ""))
	provider := vendors.NewOAuthProvider(category, "")
	tk, err := provider.StoreAccessToken(req)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("oauth error"))
		return
	}
	// todo store
	fmt.Println(tk)
	rw.Write([]byte("ok"))
}
