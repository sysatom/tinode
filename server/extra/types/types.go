package types

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store/types"
	"math/big"
	"os"
	"reflect"
	"sort"
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
	builder := MsgBuilder{Payload: a}
	if a.Title != "" {
		builder.AppendText(a.Title, TextOption{IsButton: true, ButtonDataAct: "url", ButtonDataRef: a.Url})
	} else {
		builder.AppendText(a.Url, TextOption{IsLink: true})
	}

	return builder.Content()
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
	return nil, nil
}

type DigitMsg struct {
	Title string `json:"title"`
	Digit int    `json:"digit"`
}

func (a DigitMsg) Convert() (map[string]interface{}, interface{}) {
	builder := MsgBuilder{Payload: a}
	builder.AppendText(fmt.Sprintf("Counter %s : %d", a.Title, a.Digit), TextOption{})
	return builder.Content()
}

type OkrMsg struct {
	Title     string             `json:"title"`
	Objective *model.Objective   `json:"objective"`
	KeyResult []*model.KeyResult `json:"key_result"`
}

func (o OkrMsg) Convert() (map[string]interface{}, interface{}) {
	return nil, nil
}

type InfoMsg struct {
	Title string      `json:"title"`
	Model interface{} `json:"model"`
}

func (i InfoMsg) Convert() (map[string]interface{}, interface{}) {
	builder := MsgBuilder{Payload: i}
	// title
	builder.AppendTextLine(i.Title, TextOption{})
	// model
	var m map[string]interface{}
	switch v := i.Model.(type) {
	case map[string]interface{}:
		m = v
	default:
		d, _ := json.Marshal(i.Model)
		_ = json.Unmarshal(d, &m)
	}

	// sort keys
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		builder.AppendText(fmt.Sprintf("%s: ", k), TextOption{IsBold: true})
		builder.AppendText(toString(m[k]), TextOption{})
		builder.AppendText("\n", TextOption{})
	}

	return builder.Content()
}

type TodoMsg struct {
	Title string        `json:"title"`
	Todo  []*model.Todo `json:"todo"`
}

func (t TodoMsg) Convert() (map[string]interface{}, interface{}) {
	if len(t.Todo) == 0 {
		return nil, "empty"
	}
	builder := MsgBuilder{Payload: t}
	builder.AppendTextLine("Todo", TextOption{IsBold: true})
	for i, todo := range t.Todo {
		builder.AppendTextLine(fmt.Sprintf("%d: %s", i+1, todo.Content), TextOption{})
	}
	return builder.Content()
}

type ChartMsg struct {
	Title    string    `json:"title"`
	SubTitle string    `json:"sub_title"`
	XAxis    []string  `json:"x_axis"`
	Series   []float64 `json:"series"`
}

func (t ChartMsg) Convert() (map[string]interface{}, interface{}) {
	return nil, nil
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

type RepoMsg struct {
	ID               *int64     `json:"id,omitempty"`
	NodeID           *string    `json:"node_id,omitempty"`
	Name             *string    `json:"name,omitempty"`
	FullName         *string    `json:"full_name,omitempty"`
	Description      *string    `json:"description,omitempty"`
	Homepage         *string    `json:"homepage,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	PushedAt         *time.Time `json:"pushed_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	HTMLURL          *string    `json:"html_url,omitempty"`
	Language         *string    `json:"language,omitempty"`
	Fork             *bool      `json:"fork,omitempty"`
	ForksCount       *int       `json:"forks_count,omitempty"`
	NetworkCount     *int       `json:"network_count,omitempty"`
	OpenIssuesCount  *int       `json:"open_issues_count,omitempty"`
	StargazersCount  *int       `json:"stargazers_count,omitempty"`
	SubscribersCount *int       `json:"subscribers_count,omitempty"`
	WatchersCount    *int       `json:"watchers_count,omitempty"`
	Size             *int       `json:"size,omitempty"`
	Topics           []string   `json:"topics,omitempty"`
	Archived         *bool      `json:"archived,omitempty"`
	Disabled         *bool      `json:"disabled,omitempty"`
}

func (i RepoMsg) Convert() (map[string]interface{}, interface{}) {
	builder := MsgBuilder{Payload: i}
	// title
	builder.AppendTextLine(*i.FullName, TextOption{})

	var m map[string]interface{}
	d, _ := json.Marshal(i)
	_ = json.Unmarshal(d, &m)

	// sort keys
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		builder.AppendText(fmt.Sprintf("%s: ", k), TextOption{IsBold: true})
		builder.AppendText(toString(m[k]), TextOption{})
		builder.AppendText("\n", TextOption{})
	}

	return builder.Content()
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
	key, err := generateRandomString(16)
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

func toString(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v

	case []byte:
		return string(v)

	case int:
		return strconv.Itoa(v)

	case float64:
		return strconv.FormatFloat(v, 'f', 4, 64)

	case bool:
		return strconv.FormatBool(v)

	case nil:
		return ""

	default:
		return fmt.Sprint(v)
	}
}

func AppUrl() string {
	return os.Getenv("TINODE_URL")
}

func tye(payload MsgPayload) string {
	t := reflect.TypeOf(payload)
	return t.Name()
}

func ToPayload(typ string, src []byte) MsgPayload {
	switch typ {
	case "TextMsg":
		var r TextMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "ImageMsg":
		var r ImageMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "FileMsg":
		var r FileMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "VideoMsg":
		var r VideoMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "AudioMsg":
		var r AudioMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "ScriptMsg":
		var r ScriptMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "ActionMsg":
		var r ActionMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "LinkMsg":
		var r LinkMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "LocationMsg":
		var r LocationMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "TableMsg":
		var r TableMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "DigitMsg":
		var r DigitMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "OkrMsg":
		var r OkrMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "InfoMsg":
		var r InfoMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "TodoMsg":
		var r TodoMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "ChartMsg":
		var r ChartMsg
		_ = json.Unmarshal(src, &r)
		return r
	case "RepoMsg":
		var r RepoMsg
		_ = json.Unmarshal(src, &r)
		return r
	}
	return nil
}
