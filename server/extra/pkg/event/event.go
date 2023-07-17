package event

import (
	"github.com/gookit/event"
	"github.com/tinode/chat/server/extra/types"
)

type ListenerFunc func(data types.KV) error

func eventName(name string) string {
	return name
}

func On(name string, listener ListenerFunc) {
	event.Std().On(eventName(name), event.ListenerFunc(func(e event.Event) error {
		return listener(e.Data())
	}))
}

func Emit(name string, params types.KV) error {
	err, _ := event.Std().Fire(eventName(name), params)
	return err
}

func AsyncEmit(name string, params types.KV) {
	event.Std().FireC(eventName(name), params)
}
