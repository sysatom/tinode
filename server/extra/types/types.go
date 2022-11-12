package types

import (
	"crypto/rand"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store/types"
	"math/big"
	"strconv"
	"time"
)

var commonHead = map[string]interface{}{
	"mime": "text/x-drafty",
}

type MsgPayload interface {
	Convert() (map[string]interface{}, interface{})
}

type TextMsg struct {
	Text string `json:"text"`
}

func (t TextMsg) Convert() (map[string]interface{}, interface{}) {
	return nil, t.Text
}

type ImageMsg struct {
	Src    string `json:"src"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Alt    string `json:"alt"`
}

func (i ImageMsg) Convert() (map[string]interface{}, interface{}) {
	return commonHead, nil //todo
}

type FileMsg struct {
	Src string `json:"src"`
	Alt string `json:"alt"`
}

func (i FileMsg) Convert() (map[string]interface{}, interface{}) {
	return commonHead, nil //todo
}

type VideoMsg struct {
	Src      string  `json:"src"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
	Alt      string  `json:"alt"`
	Duration float64 `json:"duration"`
}

func (i VideoMsg) Convert() (map[string]interface{}, interface{}) {
	return commonHead, nil //todo
}

type AudioMsg struct {
	Src      string  `json:"src"`
	Alt      string  `json:"alt"`
	Duration float64 `json:"duration"`
}

func (i AudioMsg) Convert() (map[string]interface{}, interface{}) {
	return commonHead, nil //todo
}

type ScriptMsg struct {
	Kind string `json:"kind"`
	Code string `json:"code"`
}

func (a ScriptMsg) Convert() (map[string]interface{}, interface{}) {
	return commonHead, nil //todo
}

type ActionMsg struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Option []string `json:"option"`
	Value  string   `json:"value"`
}

func (a ActionMsg) Convert() (map[string]interface{}, interface{}) {
	return commonHead, nil //todo
}

type LinkMsg struct {
	Title string `json:"title"`
	Cover string `json:"cover"`
	Url   string `json:"url"`
}

func (a LinkMsg) Convert() (map[string]interface{}, interface{}) {
	builder := MsgBuilder{}
	if a.Title != "" {
		builder.AppendText(a.Title, TextOption{IsButton: true, ButtonDataAct: "url", ButtonDataRef: a.Url})
	} else {
		builder.AppendText(a.Url, TextOption{IsLink: true})
	}

	return builder.Message.Content()
}

type LocationMsg struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Address   string  `json:"address"`
}

func (a LocationMsg) Convert() (map[string]interface{}, interface{}) {
	return commonHead, nil //todo
}

type TableMsg struct {
	Title  string          `json:"title"`
	Header []string        `json:"header"`
	Row    [][]interface{} `json:"row"`
}

func (t TableMsg) Convert() (map[string]interface{}, interface{}) {
	builder := MsgBuilder{}
	// title
	builder.AppendTextLine(t.Title, TextOption{})
	// header
	builder.AppendText(" | ", TextOption{})
	for _, header := range t.Header {
		builder.AppendText(header, TextOption{IsBold: true})
		builder.AppendText(" | ", TextOption{})
	}
	builder.AppendText("\n", TextOption{})
	// row
	for _, row := range t.Row {
		builder.AppendText(" | ", TextOption{})
		for _, item := range row {
			switch t := item.(type) {
			case string:
				builder.AppendText(t, TextOption{})
			case int:
				builder.AppendText(strconv.Itoa(t), TextOption{})
			}
			builder.AppendText(" | ", TextOption{})
		}
		builder.AppendText("\n", TextOption{})
	}

	return builder.Message.Content()
}

type DigitMsg struct {
	Title string `json:"title"`
	Digit int    `json:"digit"`
}

func (a DigitMsg) Convert() (map[string]interface{}, interface{}) {
	return commonHead, nil //todo
}

type OkrMsg struct {
	Title     string        `json:"title"`
	Objective interface{}   `json:"objective"`
	KeyResult []interface{} `json:"key_result"`
}

func (o OkrMsg) Convert() (map[string]interface{}, interface{}) {
	return commonHead, nil //todo
}

type InfoMsg struct {
	Title string      `json:"title"`
	Model interface{} `json:"model"`
}

func (i InfoMsg) Convert() (map[string]interface{}, interface{}) {
	return commonHead, nil //todo
}

type TodoMsg struct {
	Title string        `json:"title"`
	Todo  []interface{} `json:"todo"`
}

func (t TodoMsg) Convert() (map[string]interface{}, interface{}) {
	return commonHead, nil //todo
}

func Convert(payloads []MsgPayload) ([]map[string]interface{}, []interface{}) {
	var heads []map[string]interface{}
	var contents []interface{}
	for _, item := range payloads {
		head, content := item.Convert()
		heads = append(heads, head)
		contents = append(contents, content)
	}
	return heads, contents
}

type Context struct {
	// Message ID denormalized
	Id string `json:"-"`
	// Un-routable (original) topic name denormalized from XXX.Topic.
	Original string `json:"-"`
	// Routable (expanded) topic name.
	RcptTo string `json:"-"`
	// Sender's UserId as string.
	AsUser types.Uid `json:"-"`
	// Sender's authentication level.
	AuthLvl int `json:"-"`
	// Denormalized 'what' field of meta messages (set, get, del).
	MetaWhat int `json:"-"`
	// Timestamp when this message was received by the server.
	Timestamp time.Time `json:"-"`
	// OAuth token
	Token string `json:"-"`
	// form id
	FormId string `json:"-"`
	// seq
	SeqId int `json:"-"`
}

func Id() string {
	key, err := generateRandomString(16)
	if err != nil {
		logs.Err.Println("bot command id", err)
		return ""
	}

	uGen := types.UidGenerator{}
	err = uGen.Init(1, []byte(key))
	if err != nil {
		logs.Err.Println("bot command id", err)
		return ""
	}

	return uGen.GetStr()
}

func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
