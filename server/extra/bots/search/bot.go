package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	serverStore "github.com/tinode/chat/server/store"
	serverTypes "github.com/tinode/chat/server/store/types"
	"gorm.io/gorm"
	"strings"
)

const Name = "search"

var handler bot

func init() {
	bots.Register(Name, &handler)
}

type bot struct {
	initialized bool
	bots.Base
}

type configType struct {
	Enabled bool `json:"enabled"`
}

func (bot) Init(jsonconf json.RawMessage) error {

	// Check if the handler is already initialized
	if handler.initialized {
		return errors.New("already initialized")
	}

	var config configType
	if err := json.Unmarshal(jsonconf, &config); err != nil {
		return errors.New("failed to parse config: " + err.Error())
	}

	if !config.Enabled {
		logs.Info.Printf("bot %s disabled", Name)
		return nil
	}

	handler.initialized = true

	return nil
}

func (bot) IsReady() bool {
	return handler.initialized
}

func (b bot) Command(ctx types.Context, content interface{}) (types.MsgPayload, error) {
	return bots.RunCommand(commandRules, ctx, content)
}

func (b bot) Input(ctx types.Context, _ types.KV, content interface{}) (types.MsgPayload, error) {
	filter := ""
	if s, ok := content.(string); ok {
		filter = s
	}
	if filter == "" {
		return types.TextMsg{Text: "filter error"}, nil
	}

	items, err := store.Chatbot.SearchMessages(ctx.AsUser, ctx.RcptTo, filter)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logs.Err.Println(err)
		return types.TextMsg{Text: "Empty"}, nil
	}

	// bots
	botList, err := store.Chatbot.GetBotUsers()
	if err != nil {
		logs.Err.Println(err)
		return types.TextMsg{Text: "Empty"}, nil
	}
	botsM := make(map[uint64]string)
	for _, item := range botList {
		botsM[uint64(serverStore.EncodeUid(int64(item.ID)))] = item.Fn
	}

	// group
	groupList, err := store.Chatbot.GetGroupTopics(ctx.AsUser)
	if err != nil {
		logs.Err.Println(err)
		return types.TextMsg{Text: "Empty"}, nil
	}
	groupsM := make(map[string]string)
	for _, item := range groupList {
		groupsM[item.Name] = item.Fn
	}

	var header []string
	var row [][]interface{}
	if len(items) > 0 {
		header = []string{"Topic", "SeqId", "Content", "CreatedAt"}
		for _, v := range items {
			topic := ""
			if strings.HasPrefix(v.Topic, "p2p") {
				uid1, uid2, _ := serverTypes.ParseP2P(v.Topic)
				if uid1.IsZero() || uid2.IsZero() {
					continue
				}
				if fn, ok := botsM[uint64(uid1)]; ok {
					topic = fn
				}
				if fn, ok := botsM[uint64(uid2)]; ok {
					topic = fn
				}
			}
			if strings.HasPrefix(v.Topic, "grp") {
				if fn, ok := groupsM[v.Topic]; ok {
					topic = fn
				}
			}

			detail := v.Txt
			if detail == "" {
				detail = string(v.Raw)
			}
			row = append(row, []interface{}{fmt.Sprintf("%s (%s)", topic, v.Topic), v.Seqid, detail, v.Createdat})
		}
	}
	if len(row) == 0 {
		return types.TextMsg{Text: "Empty"}, nil
	}

	title := fmt.Sprintf("Search \"%s\" result", filter)
	return bots.StorePage(ctx, model.PageTable, title, types.TableMsg{Title: title, Header: header, Row: row}), nil
}
