package types

import (
	"github.com/tinode/chat/server/extra/utils"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store/types"
	"os"
	"time"
)

type MsgPayload interface {
	Convert() (map[string]interface{}, interface{})
}

type Context struct {
	// Message ID denormalized
	Id string
	// Un-routable (original) topic name denormalized from XXX.Topic.
	Original string
	// Routable (expanded) topic name.
	RcptTo string
	// Sender's UserId as string.
	AsUser types.Uid
	// Sender's authentication level.
	AuthLvl int
	// Denormalized 'what' field of meta messages (set, get, del).
	MetaWhat int
	// Timestamp when this message was received by the server.
	Timestamp time.Time
	// OAuth token
	Token string
	// form id
	FormId string
	// form Rule id
	FormRuleId string
	// condition
	Condition string
	// agent
	AgentId string
	// agent
	AgentVersion int
}

func Id() types.Uid {
	key, err := utils.GenerateRandomString(16)
	if err != nil {
		logs.Err.Println("bot command id", err)
		return 0
	}

	uGen := types.UidGenerator{}
	err = uGen.Init(1, []byte(key))
	if err != nil {
		logs.Err.Println("bot command id", err)
		return 0
	}

	return uGen.Get()
}

func AppUrl() string {
	return os.Getenv("TINODE_URL")
}
