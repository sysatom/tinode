package helper

type Data struct {
	Action  Action      `json:"action"`
	Version int         `json:"version"`
	Content interface{} `json:"content"`
}

type Action string

const (
	Pull  Action = "pull"
	Agent Action = "agent"
)
