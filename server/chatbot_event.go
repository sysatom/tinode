package main

import (
	"errors"
	"github.com/tinode/chat/server/extra/pkg/event"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/types/linkit"
	"github.com/tinode/chat/server/store/types"
	"net/http"
)

// send message
func onSendEvent() {
	event.On(event.SendEvent, func(data extraTypes.KV) error {
		topic, ok := data.String("topic")
		if !ok {
			return errors.New("error param topic")
		}
		topicUid, ok := data.Int64("topic_uid")
		if !ok {
			return errors.New("error param topic_uid")
		}
		message, ok := data.String("message")
		if !ok {
			return errors.New("error param message")
		}
		botSend(topic, types.Uid(topicUid), extraTypes.TextMsg{Text: message})
		return nil
	})
}

// push instruct
func onPushInstruct() {
	event.On(event.InstructEvent, func(data extraTypes.KV) error {
		uidStr, ok := data.String("uid")
		if !ok {
			return errors.New("error param uid")
		}
		uid := types.ParseUserId(uidStr)
		if uid.IsZero() {
			return errors.New("error param uid")
		}

		sessionStore.Range(func(sid string, s *Session) bool {
			if s.uid == uid {

				s.queueOutExtra(&linkit.ServerComMessage{
					Code:    http.StatusOK,
					Message: "",
					Data:    data,
				})

				return false
			}
			return true
		})

		return nil
	})
}
