package rollbar

import (
	"github.com/rollbar/rollbar-go"
	"github.com/tinode/chat/server/extra/pkg/flog"
	"github.com/tinode/chat/server/extra/vendors"
)

const (
	ID = "rollbar"

	EnableKey      = "enable"
	TokenKey       = "token"
	EnvironmentKey = "environment"
	ServerRootKey  = "server_root"
)

func Setup() error {
	enableVal, err := vendors.GetConfig(ID, EnableKey)
	if err != nil {
		return err
	}
	if !enableVal.Bool() {
		flog.Info("rollbar disable")
		return nil
	}

	tokenVal, err := vendors.GetConfig(ID, TokenKey)
	if err != nil {
		return err
	}
	envVal, err := vendors.GetConfig(ID, EnvironmentKey)
	if err != nil {
		return err
	}
	rootVal, err := vendors.GetConfig(ID, ServerRootKey)
	if err != nil {
		return err
	}
	rollbar.SetToken(tokenVal.String())
	rollbar.SetEnvironment(envVal.String())
	rollbar.SetServerRoot(rootVal.String())
	return nil
}
