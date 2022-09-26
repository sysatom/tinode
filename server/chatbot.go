package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tinode/chat/server/auth"
	"github.com/tinode/chat/server/extra/bots"
	"github.com/tinode/chat/server/extra/channels"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	"strings"
	"time"

	// bots
	_ "github.com/tinode/chat/server/extra/bots/demo"

	// channels
	_ "github.com/tinode/chat/server/extra/channels/demo"
)

// init
func botsInit(configString json.RawMessage) {
	// init bots
	err := bots.Init(string(configString))
	if err != nil {
		logs.Err.Fatal("Failed to initialize bot:", err)
	}

	// bot father
	err = initializeBotFather()
	if err != nil {
		logs.Err.Fatal("Failed to create or update bot father:", err)
	}

	// bot users
	err = initializeBotUsers()
	if err != nil {
		logs.Err.Fatal("Failed to create or update bot users:", err)
	}

	// stats register
	statsRegisterInt("BotTotal")
	statsRegisterInt("BotRunTotal")

	statsSet("BotTotal", int64(len(bots.List())))
}

func channelsInit(configString json.RawMessage) {
	err := channels.Init(string(configString))
	if err != nil {
		logs.Err.Fatal("Failed to initialize channel:", err)
	}

	err = initializeChannels()
	if err != nil {
		logs.Err.Fatal("Failed to create or update channels:", err)
	}

	// stats register
	statsRegisterInt("ChannelTotal")
	statsRegisterInt("ChannelPublishTotal")

	statsSet("ChannelTotal", int64(len(channels.List())))
}

// hook
func handleBotIncomingMessage(t *Topic, msg *ClientComMessage) {
	subs, err := store.Topics.GetUsers(msg.Pub.Topic, nil)
	if err != nil {
		logs.Err.Println(err)
		return
	}

	for _, sub := range subs {
		if !isBot(sub) {
			continue
		}

		// bot name
		name := botName(sub)
		handle, ok := bots.List()[name]
		if !ok {
			continue
		}
		heads, contents, err := handle.Run(msg.Pub.Head, msg.Pub.Content)
		if err != nil {
			logs.Warn.Printf("topic[%s]: failed to run bot: %v", t.name, err)
			continue
		}
		// multiple messages
		for i, content := range contents {
			head := heads[i]
			if content == nil {
				continue
			}

			// stats
			statsInc("BotRunTotal", 1)

			now := types.TimeNow()
			if err := store.Messages.Save(
				&types.Message{
					ObjHeader: types.ObjHeader{CreatedAt: now},
					SeqId:     t.lastID + 1,
					Topic:     t.name,
					From:      sub.User,
					Head:      head,
					Content:   content,
				}, nil, true); err != nil {
				logs.Warn.Printf("topic[%s]: failed to save bot message: %v", t.name, err)
				continue
			}

			t.lastID++
			t.touched = now

			data := &ServerComMessage{
				Data: &MsgServerData{
					Topic:     msg.Original,
					From:      sub.User,
					Timestamp: now,
					SeqId:     t.lastID,
					Head:      head,
					Content:   content,
				},
				// Internal-only values.
				Id:        msg.Id,
				RcptTo:    msg.RcptTo,
				AsUser:    sub.User,
				Timestamp: now,
				sess:      msg.sess,
			}

			t.broadcastToSessions(data)

			asUid := types.ParseUid(sub.User)

			// sendPush will update unread message count and send push notification.
			if pushRcpt := t.pushForData(asUid, data.Data); pushRcpt != nil {
				sendPush(pushRcpt)
			}
		}
	}
}

const BotFather = "BotFather"

// init bot father
func initializeBotFather() error {
	msg := &ClientComMessage{
		Acc: &MsgClientAcc{
			User:      "new",
			State:     "ok",
			AuthLevel: "",
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

// init bots
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
			Id: "1",
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
		authLvl: auth.LevelRoot,
		subs:    make(map[string]*Subscription),
		send:    make(chan interface{}, sendQueueLimit+32),
		stop:    make(chan interface{}, 1),
		detach:  make(chan string, 64),
	}

	for name, channel := range channels.List() {
		var msg = &ClientComMessage{
			Sub: &MsgClientSub{
				Topic: channel.Id(),
				Set: &MsgSetQuery{
					Desc: &MsgSetDesc{
						Public: map[string]interface{}{
							"fn":   fmt.Sprintf("%s%s", name, channels.ChannelNameSuffix),
							"note": fmt.Sprintf("%s channel", name),
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

			Original:  fmt.Sprintf("nch%s", channel.Id()),
			RcptTo:    fmt.Sprintf("grp%s", channel.Id()),
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
