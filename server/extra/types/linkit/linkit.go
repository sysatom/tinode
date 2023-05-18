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
