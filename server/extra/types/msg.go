package types

import (
	"encoding/json"
	"fmt"
	"github.com/tinode/chat/server/extra/store/model"
	"sort"
	"time"
)

type FormFieldType string

const (
	FormFieldText     FormFieldType = "text"
	FormFieldPassword FormFieldType = "password"
	FormFieldNumber   FormFieldType = "number"
	FormFieldColor    FormFieldType = "color"
	FormFieldFile     FormFieldType = "file"
	FormFieldMonth    FormFieldType = "month"
	FormFieldDate     FormFieldType = "date"
	FormFieldTime     FormFieldType = "time"
	FormFieldEmail    FormFieldType = "email"
	FormFieldUrl      FormFieldType = "url"
	FormFieldRadio    FormFieldType = "radio"
	FormFieldCheckbox FormFieldType = "checkbox"
	FormFieldRange    FormFieldType = "range"
	FormFieldSelect   FormFieldType = "select"
	FormFieldTextarea FormFieldType = "textarea"
)

type FormFieldValueType string

const (
	FormFieldValueString       FormFieldValueType = "string"
	FormFieldValueBool         FormFieldValueType = "bool"
	FormFieldValueInt64        FormFieldValueType = "int64"
	FormFieldValueFloat64      FormFieldValueType = "float64"
	FormFieldValueStringSlice  FormFieldValueType = "string_slice"
	FormFieldValueInt64Slice   FormFieldValueType = "int64_slice"
	FormFieldValueFloat64Slice FormFieldValueType = "float64_slice"
)

type TextMsg struct {
	Text string `json:"text"`
}

func (t TextMsg) Convert() (map[string]interface{}, interface{}) {
	return nil, t.Text
}

type FormMsg struct {
	ID    string      `json:"id"`
	Title string      `json:"title"`
	Field []FormField `json:"field"`
}

func (a FormMsg) Convert() (map[string]interface{}, interface{}) {
	return nil, nil
}

type FormField struct {
	Type        FormFieldType      `json:"type"`
	Key         string             `json:"key"`
	Value       interface{}        `json:"value"`
	ValueType   FormFieldValueType `json:"value_type"`
	Required    bool               `json:"required"`
	Label       string             `json:"label"`
	Placeholder string             `json:"placeholder"`
	Option      []string           `json:"option"`
	Hidden      bool               `json:"hidden"`
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

type CardMsg struct {
	Name  string
	Image string
	URI   string
	Text  string
}

func (m CardMsg) Convert() (map[string]interface{}, interface{}) {
	builder := MsgBuilder{Payload: m}
	builder.AppendText(m.Name, TextOption{IsBold: true})
	builder.AppendText(" ", TextOption{})
	builder.AppendText(m.URI, TextOption{IsLink: true})
	return builder.Content()
}

type CardListMsg struct {
	Cards []CardMsg
}

func (m CardListMsg) Convert() (map[string]interface{}, interface{}) {
	builder := MsgBuilder{Payload: m}
	for _, card := range m.Cards {
		builder.AppendText(card.Name, TextOption{IsBold: true})
		builder.AppendText(" ", TextOption{})
		builder.AppendTextLine(card.URI, TextOption{IsLink: true})
	}
	return builder.Content()
}

type HtmlMsg struct {
	Raw string
}

func (m HtmlMsg) Convert() (map[string]interface{}, interface{}) {
	return nil, nil
}
