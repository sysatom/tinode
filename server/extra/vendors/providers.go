package vendors

import (
	"fmt"
	"github.com/tinode/chat/server/store/types"
	"net/http"
)

type OAuthProvider interface {
	AuthorizeURL() string
	GetAccessToken(code string) (interface{}, error)
	Redirect(req *http.Request) (string, error)
	StoreAccessToken(req *http.Request) (map[string]interface{}, error)
}

func RedirectURI(category string, uid1, uid2 types.Uid) string {
	url := "http://127.0.0.1:6060" // todo
	return fmt.Sprintf("%s/extra/oauth/%s/%d/%d", url, category, uid1, uid2)
}
