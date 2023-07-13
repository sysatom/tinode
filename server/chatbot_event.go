package main

import (
	"errors"
	"github.com/tinode/chat/server/extra/pkg/event"
	"github.com/tinode/chat/server/extra/store/model"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/store/types"
)

func onSendEvent() {
	event.On(event.SendEvent, func(data model.JSON) error {
		topic, ok := data.String("topic")
		if !ok {
			return errors.New("error param rcpt_to")
		}
		topicUid, ok := data.Int64("topic_uid")
		if !ok {
			return errors.New("error param uid")
		}
		message, ok := data.String("message")
		if !ok {
			return errors.New("error param message")
		}
		botSend(topic, types.Uid(topicUid), extraTypes.TextMsg{Text: message})
		return nil
	})
}
