package vendors

import (
	"fmt"
	"github.com/tinode/chat/server/extra/vendors/dropbox"
	"github.com/tinode/chat/server/extra/vendors/email"
	"github.com/tinode/chat/server/extra/vendors/github"
	"github.com/tinode/chat/server/extra/vendors/pocket"
	"net/http"
)

var CredentialOptions = map[string]interface{}{
	// OAuth
	github.ID: map[string]string{
		github.ClientIdKey:     "Client ID",
		github.ClientSecretKey: "Client secrets",
	},
	pocket.ID: map[string]string{
		pocket.ClientIdKey: "Consumer Key",
	},
	dropbox.ID: map[string]string{
		dropbox.ClientIdKey:     "App key",
		dropbox.ClientSecretKey: "App secret",
	},

	// Service
	email.ID: map[string]string{
		email.Host:     "SMTP Host",
		email.Port:     "SMTP Port",
		email.Username: "Username Mail",
		email.Password: "Password",
	},
}

type OAuthProvider interface {
	AuthorizeURL() string
	GetAccessToken(code string) (interface{}, error)
	Redirect(req *http.Request) (string, error)
	StoreAccessToken(req *http.Request) (interface{}, error)
}

func NewOAuthProvider(category, url string) OAuthProvider {
	redirectURI := fmt.Sprintf("%s/extra/oauth/%s", url, category)
	var provider OAuthProvider

	switch category {
	case pocket.ID:
		p := pocket.NewPocket("", "", redirectURI, "")
		provider = p
	case github.ID:
		provider = github.NewGithub("", "", redirectURI, "")
	case dropbox.ID:
		provider = dropbox.NewDropbox("", "", redirectURI, "")
	default:
		return nil
	}

	return provider
}

func RedirectURI(category string) string {
	url := "http://127.0.0.1:6060" // todo
	return fmt.Sprintf("%s/extra/oauth/%s", url, category)
}
