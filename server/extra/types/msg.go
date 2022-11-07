package types

import (
	"encoding/json"
	"strings"
)

type FmtMessage struct {
	At  int    `json:"at,omitempty"`
	Len int    `json:"len,omitempty"`
	Tp  string `json:"tp,omitempty"`
	Key int    `json:"key,omitempty"`
}

type EntMessage struct {
	Tp   string  `json:"tp,omitempty"`
	Data EntData `json:"data"`
}

type EntData struct {
	Mime   string      `json:"mime,omitempty"`
	Val    interface{} `json:"val,omitempty"`
	Url    string      `json:"url,omitempty"`
	Ref    string      `json:"ref,omitempty"`
	Width  int         `json:"width,omitempty"`
	Height int         `json:"height,omitempty"`
	Name   string      `json:"name,omitempty"`
	Size   int         `json:"size,omitempty"`
	Act    string      `json:"act,omitempty"`
}

type ChatMessage struct {
	Text        string       `json:"txt,omitempty"`
	Fmt         []FmtMessage `json:"fmt,omitempty"`
	Ent         []EntMessage `json:"ent,omitempty"`
	IsPlainText bool         `json:"-"`
	MessageType string       `json:"-"`
}

// GetFormattedText Get original text message, inlude original '\n'
func (c ChatMessage) GetFormattedText() string {
	if c.Text == "" {
		return ""
	}
	t := []byte(c.Text)
	for _, item := range c.Fmt {
		if item.Tp == "BR" {
			t[item.At] = '\n'
		}
	}
	return string(t)
}

// GetEntDatas get entity from chat message by entity type
func (c ChatMessage) GetEntDatas(tp string) []EntData {
	var ret []EntData
	for _, item := range c.Ent {
		if item.Tp == tp {
			ret = append(ret, item.Data)
		}
	}
	return ret
}

// GetMentions get mentioned users
func (c ChatMessage) GetMentions() []EntData {
	return c.GetEntDatas("MN")
}

// GetImages get images
func (c ChatMessage) GetImages() []EntData {
	return c.GetEntDatas("IM")
}

// GetHashTags get hashtags
func (c ChatMessage) GetHashTags() []EntData {
	return c.GetEntDatas("HT")
}

// GetLinks get links
func (c ChatMessage) GetLinks() []EntData {
	return c.GetEntDatas("LN")
}

// GetGenericAttachment get generic attachment
func (c ChatMessage) GetGenericAttachment() []EntData {
	return c.GetEntDatas("EX")
}

func (c ChatMessage) Content() (map[string]interface{}, interface{}) {
	if c.IsPlainText {
		return nil, c.Text
	}
	d, err := json.Marshal(c)
	if err != nil {
		return nil, ""
	}
	return map[string]interface{}{
		"mime": "text/x-drafty",
	}, json.RawMessage(d)
}

type MsgBuilder struct {
	Message ChatMessage
}

// AppendText Append text message to build message
func (m *MsgBuilder) AppendText(text string, opt TextOption) {
	baseLen := len(m.Message.Text)
	m.Message.Text += text
	if strings.Contains(text, "\n") {
		for i := 0; i < len(text); i++ {
			if text[i] == '\n' {
				fmt := FmtMessage{
					At:  baseLen + i,
					Tp:  "BR",
					Len: 1,
				}
				m.Message.Fmt = append(m.Message.Fmt, fmt)
			}
		}
	}

	leftLen := baseLen + (len(text) - len(strings.TrimLeft(text, "\t\n\v\f\r ")))
	subLen := len(text) - len(strings.TrimRight(text, "\t\n\v\f\r "))
	validLen := len(m.Message.Text) - leftLen - subLen

	if opt.IsBold {
		fmt := FmtMessage{
			Tp:  "ST",
			At:  leftLen,
			Len: validLen,
		}
		m.Message.Fmt = append(m.Message.Fmt, fmt)
	}
	if opt.IsItalic {
		fmt := FmtMessage{
			Tp:  "EM",
			At:  leftLen,
			Len: validLen,
		}
		m.Message.Fmt = append(m.Message.Fmt, fmt)
	}
	if opt.IsDeleted {
		fmt := FmtMessage{
			Tp:  "DL",
			At:  leftLen,
			Len: validLen,
		}
		m.Message.Fmt = append(m.Message.Fmt, fmt)
	}
	if opt.IsCode {
		fmt := FmtMessage{
			Tp:  "CO",
			At:  leftLen,
			Len: validLen,
		}
		m.Message.Fmt = append(m.Message.Fmt, fmt)
	}
	if opt.IsLink {
		fmt := FmtMessage{
			At:  leftLen,
			Len: validLen,
			Key: len(m.Message.Ent),
		}
		m.Message.Fmt = append(m.Message.Fmt, fmt)
		url := strings.ToLower(strings.TrimSpace(text))
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + strings.TrimSpace(text)
		}
		ent := EntMessage{
			Tp: "LN",
			Data: EntData{
				Url: url,
			},
		}
		m.Message.Ent = append(m.Message.Ent, ent)
	}
	if opt.IsMention {
		fmt := FmtMessage{
			At:  leftLen,
			Len: validLen,
			Key: len(m.Message.Ent),
		}
		m.Message.Fmt = append(m.Message.Fmt, fmt)
		mentionName := substr(strings.TrimSpace(text), 1, len(strings.TrimSpace(text))-1)
		ent := EntMessage{
			Tp: "MN",
			Data: EntData{
				Val: mentionName,
			},
		}
		m.Message.Ent = append(m.Message.Ent, ent)
	}
	if opt.IsHashTag {
		fmt := FmtMessage{
			At:  leftLen,
			Len: validLen,
			Key: len(m.Message.Ent),
		}
		m.Message.Fmt = append(m.Message.Fmt, fmt)
		hashTag := strings.TrimSpace(text)
		ent := EntMessage{
			Tp: "HT",
			Data: EntData{
				Val: hashTag,
			},
		}
		m.Message.Ent = append(m.Message.Ent, ent)
	}
	if opt.IsForm {
		fmt := FmtMessage{
			Tp:  "FM",
			At:  leftLen,
			Len: validLen,
		}
		m.Message.Fmt = append(m.Message.Fmt, fmt)
	}
	if opt.IsButton {
		var btnName = opt.ButtonDataName
		if btnName == "" {
			opt.ButtonDataName = strings.ToLower(strings.TrimSpace(text))
		}
		fmt := FmtMessage{
			At:  leftLen,
			Len: validLen,
			Key: len(m.Message.Ent),
		}
		m.Message.Fmt = append(m.Message.Fmt, fmt)
		//btnText := strings.TrimSpace(text)
		ent := EntMessage{
			Tp: "BN",
			Data: EntData{
				Name: opt.ButtonDataName,
				Act:  opt.ButtonDataAct,
				Val:  opt.ButtonDataVal,
			},
		}
		m.Message.Ent = append(m.Message.Ent, ent)
	}
}

// AppendTextLine Append text message and line break to build message
func (m *MsgBuilder) AppendTextLine(text string, opt TextOption) {
	m.AppendText(text+"\n", opt)
}

// AppendImage Append image to build message
func (m *MsgBuilder) AppendImage(imageName string, opt ImageOption) {
	m.Message.Fmt = append(m.Message.Fmt, FmtMessage{
		At:  len(m.Message.Text),
		Len: 1,
		Key: len(m.Message.Ent),
	})
	m.Message.Ent = append(m.Message.Ent, EntMessage{
		Tp: "IM",
		Data: EntData{
			Mime:   opt.Mime,
			Width:  opt.Width,
			Height: opt.Height,
			Name:   imageName,
			Val:    opt.ImageBase64,
		},
	})
}

// AppendFile Append file to build message
func (m *MsgBuilder) AppendFile(fileName string, opt FileOption) {
	m.Message.Fmt = append(m.Message.Fmt, FmtMessage{
		At:  len(m.Message.Text),
		Len: 0,
		Key: len(m.Message.Ent),
	})
	m.Message.Ent = append(m.Message.Ent, EntMessage{
		Tp: "EX",
		Data: EntData{
			Mime: opt.Mime,
			Name: fileName,
			Val:  opt.ContentBase64,
		},
	})
}

// AppendAttachment append a attachment file to chat message
func (m *MsgBuilder) AppendAttachment(fileName string, opt AttachmentOption) {
	m.Message.Fmt = append(m.Message.Fmt, FmtMessage{
		At:  len(m.Message.Text),
		Len: 1,
		Key: len(m.Message.Ent),
	})
	m.Message.Ent = append(m.Message.Ent, EntMessage{
		Tp: "EX",
		Data: EntData{
			Mime: opt.Mime,
			Name: fileName,
			Ref:  opt.RelativeUrl,
			Size: opt.Size,
		},
	})
}

// Parse a raw ServerData to friendly ChatMessage
func (m *MsgBuilder) Parse(message ServerData) (ChatMessage, error) {
	var chatMsg ChatMessage
	if strings.Contains(message.Head, "mime") {
		err := json.Unmarshal([]byte(message.Content), &chatMsg)
		if err != nil {
			return ChatMessage{}, err
		}
		chatMsg.IsPlainText = false
	} else {
		err := json.Unmarshal([]byte(message.Content), &chatMsg)
		if err != nil {
			return ChatMessage{}, err
		}
		chatMsg.IsPlainText = true
	}
	if strings.HasPrefix(message.Topic, "usr") {
		chatMsg.MessageType = "user"
	}
	if strings.HasPrefix(message.Topic, "grp") {
		chatMsg.MessageType = "group"
	}
	return chatMsg, nil
}

// BuildTextMessage build text chat message with formatted
func (m *MsgBuilder) BuildTextMessage(text string) ChatMessage {
	msg := ChatMessage{}
	msg.Text = text
	msg.Ent = []EntMessage{}
	msg.Fmt = []FmtMessage{}
	if strings.Contains(text, "\n") {
		for i := 0; i < len(text); i++ {
			if text[i] == '\n' {
				fmt := FmtMessage{
					At:  i,
					Tp:  "BR",
					Len: 1,
				}
				msg.Fmt = append(msg.Fmt, fmt)
			}
		}
	}
	return msg
}

// BuildImageMessage build a image chat message
func (m *MsgBuilder) BuildImageMessage(imageName string, text string, opt ImageOption) ChatMessage {
	msg := ChatMessage{}
	msg.Text = text
	msg.Ent = []EntMessage{}
	msg.Fmt = []FmtMessage{}
	msg.Ent = append(msg.Ent, EntMessage{
		Tp: "IM",
		Data: EntData{
			Mime:   opt.Mime,
			Width:  opt.Width,
			Height: opt.Height,
			Name:   imageName,
			Val:    opt.ImageBase64,
		},
	})
	msg.Fmt = append(msg.Fmt, FmtMessage{
		At:  len(text),
		Len: 1,
		Key: 0,
	})
	if strings.Contains(text, "\n") {
		for i := 0; i < len(text); i++ {
			if text[i] == '\n' {
				fmt := FmtMessage{
					At:  i,
					Tp:  "BR",
					Len: 1,
				}
				msg.Fmt = append(msg.Fmt, fmt)
			}
		}
	}
	return msg
}

// BuildFileMessage build a file chat message
func (m *MsgBuilder) BuildFileMessage(fileName string, text string, opt FileOption) ChatMessage {
	msg := ChatMessage{}
	msg.Text = text
	msg.Ent = []EntMessage{}
	msg.Fmt = []FmtMessage{}
	msg.Ent = append(msg.Ent, EntMessage{
		Tp: "EX",
		Data: EntData{
			Mime: opt.Mime,
			Name: fileName,
			Val:  opt.ContentBase64,
		},
	})
	msg.Fmt = append(msg.Fmt, FmtMessage{
		At:  len(text),
		Len: 0,
		Key: 0,
	})
	if strings.Contains(text, "\n") {
		for i := 0; i < len(text); i++ {
			if text[i] == '\n' {
				fmt := FmtMessage{
					At:  i,
					Tp:  "BR",
					Len: 1,
				}
				msg.Fmt = append(msg.Fmt, fmt)
			}
		}
	}
	return msg
}

// BuildAttachmentMessage build a attachment message
func (m *MsgBuilder) BuildAttachmentMessage(fileName string, text string, opt AttachmentOption) ChatMessage {
	msg := ChatMessage{}
	msg.Text = text
	msg.Ent = []EntMessage{}
	msg.Fmt = []FmtMessage{}
	msg.Ent = append(msg.Ent, EntMessage{
		Tp: "EX",
		Data: EntData{
			Mime: opt.Mime,
			Name: fileName,
			Ref:  opt.RelativeUrl,
			Size: opt.Size,
		},
	})
	msg.Fmt = append(msg.Fmt, FmtMessage{
		At:  len(text),
		Len: 1,
		Key: 0,
	})
	if strings.Contains(text, "\n") {
		for i := 0; i < len(text); i++ {
			if text[i] == '\n' {
				fmt := FmtMessage{
					At:  i,
					Tp:  "BR",
					Len: 1,
				}
				msg.Fmt = append(msg.Fmt, fmt)
			}
		}
	}
	return msg
}

type ServerData struct {
	Head    string
	Content string
	Topic   string
}

type TextOption struct {
	IsBold         bool
	IsItalic       bool
	IsDeleted      bool
	IsCode         bool
	IsLink         bool
	IsMention      bool
	IsHashTag      bool
	IsForm         bool
	IsButton       bool
	ButtonDataName string
	ButtonDataAct  string
	ButtonDataVal  string
}

type ImageOption struct {
	Mime        string
	Width       int
	Height      int
	ImageBase64 string
}

type FileOption struct {
	Mime          string
	ContentBase64 string
}

type AttachmentOption struct {
	Mime        string
	RelativeUrl string
	Size        int
}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}
