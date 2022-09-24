package types

type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeImage    MessageType = "image"
	MessageTypeFile     MessageType = "file"
	MessageTypeLocation MessageType = "location"
	MessageTypeVideo    MessageType = "video"
	MessageTypeLink     MessageType = "link"
	MessageTypeScript   MessageType = "script"
	MessageTypeAction   MessageType = "action"
	MessageTypeForm     MessageType = "form"
	MessageTypeTable    MessageType = "table"
	MessageTypeDigit    MessageType = "digit"
	MessageTypeOkr      MessageType = "okr"
	MessageTypeInfo     MessageType = "info"
	MessageTypeTodo     MessageType = "todo"
)

type MsgPayload interface {
	Type() MessageType
}

type TextMsg struct {
	Text string `json:"text"`
}

func (t TextMsg) Type() MessageType {
	return MessageTypeText
}

type ImageMsg struct {
	Src    string `json:"src"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Alt    string `json:"alt"`
}

func (i ImageMsg) Type() MessageType {
	return MessageTypeImage
}

type FileMsg struct {
	Src string `json:"src"`
	Alt string `json:"alt"`
}

func (i FileMsg) Type() MessageType {
	return MessageTypeFile
}

type VideoMsg struct {
	Src      string  `json:"src"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
	Alt      string  `json:"alt"`
	Duration float64 `json:"duration"`
}

func (i VideoMsg) Type() MessageType {
	return MessageTypeVideo
}

type AudioMsg struct {
	Src      string  `json:"src"`
	Alt      string  `json:"alt"`
	Duration float64 `json:"duration"`
}

func (i AudioMsg) Type() MessageType {
	return MessageTypeAudio
}

type ScriptMsg struct {
	Kind string `json:"kind"`
	Code string `json:"code"`
}

func (a ScriptMsg) Type() MessageType {
	return MessageTypeScript
}

type ActionMsg struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Option []string `json:"option"`
	Value  string   `json:"value"`
}

func (a ActionMsg) Type() MessageType {
	return MessageTypeAction
}

type FormMsg struct {
	ID    string      `json:"id"`
	Title string      `json:"title"`
	Field []FormField `json:"field"`
}

func (a FormMsg) Type() MessageType {
	return MessageTypeForm
}

type FormField struct {
	Key      string      `json:"key"`
	Type     string      `json:"type"`
	Required bool        `json:"required"`
	Value    interface{} `json:"value"`
	Default  interface{} `json:"default"`
	Intro    string      `json:"intro"`
}

type LinkMsg struct {
	Title string `json:"title"`
	Cover string `json:"cover"`
	Url   string `json:"url"`
}

func (a LinkMsg) Type() MessageType {
	return MessageTypeLink
}

type LocationMsg struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Address   string  `json:"address"`
}

func (a LocationMsg) Type() MessageType {
	return MessageTypeLocation
}

type TableMsg struct {
	Title  string          `json:"title"`
	Header []string        `json:"header"`
	Row    [][]interface{} `json:"row"`
}

func (t TableMsg) Type() MessageType {
	return MessageTypeTable
}

type DigitMsg struct {
	Title string `json:"title"`
	Digit int    `json:"digit"`
}

func (a DigitMsg) Type() MessageType {
	return MessageTypeDigit
}

type OkrMsg struct {
	Title     string        `json:"title"`
	Objective interface{}   `json:"objective"`
	KeyResult []interface{} `json:"key_result"`
}

func (o OkrMsg) Type() MessageType {
	return MessageTypeOkr
}

type InfoMsg struct {
	Title string      `json:"title"`
	Model interface{} `json:"model"`
}

func (i InfoMsg) Type() MessageType {
	return MessageTypeInfo
}

type TodoMsg struct {
	Title string        `json:"title"`
	Todo  []interface{} `json:"todo"`
}

func (t TodoMsg) Type() MessageType {
	return MessageTypeTodo
}

func Convert(payloads []MsgPayload) ([]map[string]interface{}, []interface{}) {
	var heads []map[string]interface{}
	var contents []interface{}
	for _, item := range payloads {
		switch v := item.(type) {
		case TextMsg:
			heads = append(heads, nil)
			contents = append(contents, v.Text)
		}
	}
	return heads, contents
}
