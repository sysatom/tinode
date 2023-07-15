package types

import (
	"encoding/json"
	"github.com/tinode/chat/server/extra/store/model"
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
	// seq id
	SeqId int
	// form Rule id
	ActionRuleId string
	// condition
	Condition string
	// agent
	AgentId string
	// agent
	AgentVersion int
	// session Rule id
	SessionRuleId string
	// session init values
	SessionInitValues model.JSON
	// session last values
	SessionLastValues model.JSON
	// group event
	GroupEvent GroupEvent
	// workflow flag id
	WorkflowFlag string
	// workflow rule id
	WorkflowRuleId string
	// workflow version
	WorkflowVersion int
	// workflow step index
	WorkflowStepIndex int
	// page rule id
	PageRuleId string
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

type QueuePayload struct {
	RcptTo string          `json:"rcpt_to"`
	Uid    string          `json:"uid"`
	Type   string          `json:"type"`
	Msg    json.RawMessage `json:"msg"`
}

func ConvertQueuePayload(rcptTo string, uid string, msg MsgPayload) (QueuePayload, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return QueuePayload{}, err
	}
	return QueuePayload{
		RcptTo: rcptTo,
		Uid:    uid,
		Type:   Tye(msg),
		Msg:    data,
	}, nil
}

type DataFilter struct {
	Prefix       *string
	CreatedStart *time.Time
	CreatedEnd   *time.Time
}

type SendFunc func(rcptTo string, uid types.Uid, out MsgPayload, option ...interface{})

func WithContext(ctx Context) Context {
	return ctx
}
