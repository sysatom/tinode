package linkit

import "github.com/tinode/chat/server/extra/types"

type Data struct {
	Action  Action   `json:"action"`
	Version int      `json:"version"`
	Content types.KV `json:"content"`
}

type Action string

const (
	Info  Action = "info"
	Pull  Action = "pull"
	Agent Action = "agent"
	Bots  Action = "bots"
	Help  Action = "help"
	Ack   Action = "ack"
)

// ClientComMessage is a wrapper for client messages.
type ClientComMessage struct {
	Data Data `json:"data"`
}

// ServerComMessage is a wrapper for server-side messages.
type ServerComMessage struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
