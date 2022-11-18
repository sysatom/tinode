package main

import (
	"errors"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/tinode/chat/server/auth"
	"github.com/tinode/chat/server/extra/bots"
	botGithub "github.com/tinode/chat/server/extra/bots/github"
	botPocket "github.com/tinode/chat/server/extra/bots/pocket"
	"github.com/tinode/chat/server/extra/channels"
	"github.com/tinode/chat/server/extra/channels/crawler"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/vendors"
	"github.com/tinode/chat/server/extra/vendors/dropbox"
	"github.com/tinode/chat/server/extra/vendors/github"
	"github.com/tinode/chat/server/extra/vendors/pocket"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	"math/rand"
	"strings"
	"time"
)

const BotFather = "BotFather"

func initializeBotFather() error {
	msg := &ClientComMessage{
		Acc: &MsgClientAcc{
			User:      "new",
			State:     "ok",
			AuthLevel: "auth",
			Token:     nil,
			Scheme:    "basic",
			Secret:    []byte(fmt.Sprintf("%s:170953280278461931", BotFather)),
			Login:     false,
			Tags:      []string{"bot"},
			Desc: &MsgSetDesc{
				DefaultAcs: nil,
				Public: map[string]interface{}{
					"fn": BotFather,
				},
				Trusted: map[string]interface{}{
					"staff": true,
				},
				Private: nil,
			},
		},
		AuthLvl: int(auth.LevelRoot),
	}

	authhdl := store.Store.GetLogicalAuthHandler("basic")

	// Check if login is unique.
	if ok, _ := authhdl.IsUnique(msg.Acc.Secret, ""); !ok {
		return nil
	}

	var user types.User
	var private interface{}

	if msg.Acc.State != "" {
		state, err := types.NewObjState(msg.Acc.State)
		if err != nil {
			return err
		}
		user.State = state
	}

	// Ensure tags are unique and not restricted.
	if tags := normalizeTags(msg.Acc.Tags); tags != nil {
		if !restrictedTagsEqual(tags, nil, globals.immutableTagNS) {
			return errors.New("create user: attempt to directly assign restricted tags")
		}
		user.Tags = tags
	}

	// Assign default access values in case the acc creator has not provided them
	user.Access.Auth = getDefaultAccess(types.TopicCatP2P, true, false) |
		getDefaultAccess(types.TopicCatGrp, true, false)
	user.Access.Anon = getDefaultAccess(types.TopicCatP2P, false, false) |
		getDefaultAccess(types.TopicCatGrp, false, false)

	// Assign actual access values, public and private.
	if msg.Acc.Desc != nil {
		if !isNullValue(msg.Acc.Desc.Public) {
			user.Public = msg.Acc.Desc.Public
		}
		if !isNullValue(msg.Acc.Desc.Trusted) {
			user.Trusted = msg.Acc.Desc.Trusted
		}
		if !isNullValue(msg.Acc.Desc.Private) {
			private = msg.Acc.Desc.Private
		}
	}

	// Create user record in the database.
	if _, err := store.Users.Create(&user, private); err != nil {
		return fmt.Errorf("create bot user: failed to create bot user, %s", err)
	}

	// Add authentication record. The authhdl.AddRecord may change tags.
	_, err := authhdl.AddRecord(&auth.Rec{Uid: user.Uid(), Tags: user.Tags}, msg.Acc.Secret, "")
	if err != nil {
		return fmt.Errorf("create bot user: add auth record failed, %s", err)
	}

	return nil
}

func initializeBotUsers() error {
	var msgs []*ClientComMessage

	for name := range bots.List() {
		msgs = append(msgs, &ClientComMessage{
			Acc: &MsgClientAcc{
				User:      "new",
				State:     "ok",
				AuthLevel: "",
				Token:     nil,
				Scheme:    "basic",
				Secret:    []byte(fmt.Sprintf("%s%s:170953280278461931", name, bots.BotNameSuffix)),
				Login:     false,
				Tags:      []string{"bot"},
				Desc: &MsgSetDesc{
					DefaultAcs: nil,
					Public: map[string]interface{}{
						"fn": fmt.Sprintf("%s%s", name, bots.BotNameSuffix),
					},
					Trusted: map[string]interface{}{
						"verified": true,
					},
					Private: nil,
				},
			},
			AuthLvl: int(auth.LevelRoot),
		})
	}

	authhdl := store.Store.GetLogicalAuthHandler("basic")

	for _, msg := range msgs {
		// Check if login is unique.
		if ok, _ := authhdl.IsUnique(msg.Acc.Secret, ""); !ok {
			continue
		}

		var user types.User
		var private interface{}

		if msg.Acc.State != "" {
			state, err := types.NewObjState(msg.Acc.State)
			if err != nil {
				return err
			}
			user.State = state
		}

		// Ensure tags are unique and not restricted.
		if tags := normalizeTags(msg.Acc.Tags); tags != nil {
			if !restrictedTagsEqual(tags, nil, globals.immutableTagNS) {
				return errors.New("create user: attempt to directly assign restricted tags")
			}
			user.Tags = tags
		}

		// Assign default access values in case the acc creator has not provided them
		user.Access.Auth = getDefaultAccess(types.TopicCatP2P, true, false) |
			getDefaultAccess(types.TopicCatGrp, true, false)
		user.Access.Anon = getDefaultAccess(types.TopicCatP2P, false, false) |
			getDefaultAccess(types.TopicCatGrp, false, false)

		// Assign actual access values, public and private.
		if msg.Acc.Desc != nil {
			if !isNullValue(msg.Acc.Desc.Public) {
				user.Public = msg.Acc.Desc.Public
			}
			if !isNullValue(msg.Acc.Desc.Trusted) {
				user.Trusted = msg.Acc.Desc.Trusted
			}
			if !isNullValue(msg.Acc.Desc.Private) {
				private = msg.Acc.Desc.Private
			}
		}

		// Create user record in the database.
		if _, err := store.Users.Create(&user, private); err != nil {
			return fmt.Errorf("create bot user: failed to create bot user, %s", err)
		}

		// Add authentication record. The authhdl.AddRecord may change tags.
		_, err := authhdl.AddRecord(&auth.Rec{Uid: user.Uid(), Tags: user.Tags}, msg.Acc.Secret, "")
		if err != nil {
			return fmt.Errorf("create bot user: add auth record failed, %s", err)
		}
	}
	return nil
}

func isBot(subs types.Subscription) bool {
	// normal bot user
	if subs.GetState() != types.StateOK {
		return false
	}
	// verified
	trusted := subs.GetTrusted()
	if trusted == nil {
		return false
	}
	if !isVerified(trusted) {
		return false
	}
	// check name
	public := subs.GetPublic()
	if public == nil {
		return false
	}
	name := fn(public)
	if !strings.HasSuffix(name, bots.BotNameSuffix) {
		return false
	}

	return true
}

func isVerified(trusted interface{}) bool {
	if v, ok := trusted.(map[string]interface{}); ok {
		if b, ok := v["verified"]; ok {
			if vv, ok := b.(bool); ok {
				return vv
			}
		}
	}
	return false
}

func fn(public interface{}) string {
	if v, ok := public.(map[string]interface{}); ok {
		if s, ok := v["fn"]; ok {
			if ss, ok := s.(string); ok {
				return ss
			}
		}
	}
	return ""
}

func botName(subs types.Subscription) string {
	public := subs.GetPublic()
	if public == nil {
		return ""
	}
	name := fn(public)
	name = strings.ReplaceAll(name, bots.BotNameSuffix, "")
	return name
}

// init channels
func initializeChannels() error {
	// bind to BotFather
	uid, _, _, _, err := store.Users.GetAuthUniqueRecord("basic", "botfather")
	if err != nil {
		return err
	}
	sess := &Session{
		uid:     uid,
		authLvl: auth.LevelAuth,
		subs:    make(map[string]*Subscription),
		send:    make(chan interface{}, sendQueueLimit+32),
		stop:    make(chan interface{}, 1),
		detach:  make(chan string, 64),
	}

	for _, channel := range channels.List() {
		topic, _ := store.Topics.Get(fmt.Sprintf("grp%s", channel.Id))
		if topic != nil && topic.Id != "" {
			logs.Info.Printf("channel %s registered", channel.Id)
			continue
		}

		var msg = &ClientComMessage{
			Sub: &MsgClientSub{
				Topic: channel.Name,
				Set: &MsgSetQuery{
					Desc: &MsgSetDesc{
						Public: map[string]interface{}{
							"fn":   fmt.Sprintf("%s%s", channel.Name, channels.ChannelNameSuffix),
							"note": fmt.Sprintf("%s channel", channel.Name),
						},
						Trusted: map[string]interface{}{
							"verified": true,
						},
					},
					Tags: []string{"channel"},
				},
				Created: false,
				Newsub:  false,
			},

			Original:  fmt.Sprintf("nch%s", channel.Id),
			RcptTo:    fmt.Sprintf("grp%s", channel.Id),
			AsUser:    uid.UserId(),
			AuthLvl:   int(auth.LevelRoot),
			Timestamp: time.Now(),
			init:      true,
			sess:      sess,
		}

		globals.hub.join <- msg

		statsInc("LiveTopics", 1)
		statsInc("TotalTopics", 1)
	}

	return nil
}

// init crawler
func initializeCrawler() error {
	uid, _, _, _, err := store.Users.GetAuthUniqueRecord("basic", "botfather")
	if err != nil {
		return err
	}

	c := crawler.New()
	c.Send = func(id, name string, out [][]byte) {
		if len(out) == 0 {
			return
		}
		topic := fmt.Sprintf("grp%s", id)
		dst, err := store.Topics.Get(topic)
		if err != nil {
			logs.Err.Println("init crawler", err)
			return
		}
		if dst == nil {
			return
		}
		builder := extraTypes.MsgBuilder{}
		for _, i := range out {
			builder.AppendTextLine(string(i), extraTypes.TextOption{})
		}
		head, content := builder.Message.Content()

		msg := &ClientComMessage{
			Pub: &MsgClientPub{
				Topic:   topic,
				Head:    head,
				Content: content,
			},
			AsUser:    uid.UserId(),
			Timestamp: types.TimeNow(),
		}

		t := &Topic{
			name:   topic,
			cat:    types.TopicCatGrp,
			status: topicStatusLoaded,
			lastID: dst.SeqId,
			perUser: map[types.Uid]perUserData{
				uid: {
					modeGiven: types.ModeCFull,
					modeWant:  types.ModeCFull,
					private:   nil,
				},
			},
		}
		t.handleClientMsg(msg)
	}

	var rules []crawler.Rule
	for _, publisher := range channels.List() {
		rules = append(rules, *publisher)
	}

	err = c.Init(rules...)
	if err != nil {
		return err
	}
	c.Run()
	return nil
}

func botSend(userUid, topicUid types.Uid, out extraTypes.MsgPayload) {
	if out == nil {
		return
	}

	topic := userUid.P2PName(topicUid)

	t := globals.hub.topicGet(topic)
	if t == nil {
		sess := &Session{
			uid:     topicUid,
			authLvl: auth.LevelAuth,
			subs:    make(map[string]*Subscription),
			send:    make(chan interface{}, sendQueueLimit+32),
			stop:    make(chan interface{}, 1),
			detach:  make(chan string, 64),
		}
		msg := &ClientComMessage{
			Sub: &MsgClientSub{
				Topic:   topicUid.UserId(),
				Get:     &MsgGetQuery{},
				Created: false,
				Newsub:  false,
			},
			Original:  topicUid.UserId(),
			RcptTo:    topicUid.P2PName(userUid),
			AsUser:    userUid.UserId(),
			AuthLvl:   int(auth.LevelAuth),
			MetaWhat:  0,
			Timestamp: time.Now(),
			sess:      sess,
			init:      true,
		}
		globals.hub.join <- msg
		// wait sometime
		time.Sleep(200 * time.Millisecond)

		t = globals.hub.topicGet(topic)
	}

	if t == nil {
		logs.Err.Printf("topic %s error, Failed to send", topic)
		return
	}

	heads, contents := extraTypes.Convert([]extraTypes.MsgPayload{out})
	if !(len(heads) > 0 && len(contents) > 0) {
		logs.Err.Printf("topic %s convert error, Failed to send", topic)
		return
	}
	head, content := heads[0], contents[0]
	msg := &ClientComMessage{
		Pub: &MsgClientPub{
			Topic:   topic,
			Head:    head,
			Content: content,
		},
		AsUser:    topicUid.UserId(),
		Timestamp: types.TimeNow(),
	}
	t.handleClientMsg(msg)
}

func newProvider(category string) vendors.OAuthProvider {
	var provider vendors.OAuthProvider

	switch category {
	case pocket.ID:
		provider = pocket.NewPocket(botPocket.Config.ConsumerKey, "", "", "")
	case github.ID:
		provider = github.NewGithub(botGithub.Config.ID, botGithub.Config.Secret, "", "")
	case dropbox.ID:
		provider = dropbox.NewDropbox("", "", "", "")
	default:
		return nil
	}

	return provider
}

// generate random data for bar chart
func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}
