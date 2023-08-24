package vendors

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"net/http"
)

type OAuthProvider interface {
	GetAuthorizeURL() string
	GetAccessToken(req *http.Request) (extraTypes.KV, error)
}

func RedirectURI(name string, flag string) string {
	return fmt.Sprintf("%s/extra/oauth/%s/%s", extraTypes.AppUrl(), name, flag)
}

var Configs json.RawMessage

func GetConfig(name, key string) (gjson.Result, error) {
	if len(Configs) == 0 {
		return gjson.Result{}, errors.New("error configs")
	}
	value := gjson.GetBytes(Configs, fmt.Sprintf("%s.%s", name, key))
	return value, nil
}
