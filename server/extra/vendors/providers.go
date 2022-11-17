package vendors

import (
	"fmt"
	extraTypes "github.com/tinode/chat/server/extra/types"
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
	return fmt.Sprintf("%s/extra/oauth/%s/%d/%d", extraTypes.AppUrl(), category, uid1, uid2)
}
