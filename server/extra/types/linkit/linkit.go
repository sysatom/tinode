package linkit

type Data struct {
	Action  Action      `json:"action"`
	Version int         `json:"version"`
	Content interface{} `json:"content"`
}

type Action string

const (
	Info  Action = "info"
	Pull  Action = "pull"
	Agent Action = "agent"
	Bots  Action = "bots"
	Help  Action = "help"
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
