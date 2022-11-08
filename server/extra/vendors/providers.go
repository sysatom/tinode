package vendors

import (
	"fmt"
	"net/http"
)

type OAuthProvider interface {
	AuthorizeURL() string
	GetAccessToken(code string) (interface{}, error)
	Redirect(req *http.Request) (string, error)
	StoreAccessToken(req *http.Request) (map[string]interface{}, error)
}

func RedirectURI(category string) string {
	url := "http://127.0.0.1:6060" // todo
	return fmt.Sprintf("%s/extra/oauth/%s/usr1ScPwXm5MJg/usrpc-w_TB9ma4", url, category)
}
