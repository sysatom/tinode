package event

import (
	"github.com/gookit/event"
	"github.com/tinode/chat/server/extra/store/model"
)

type ListenerFunc func(data model.JSON) error

func eventName(name string) string {
	return name
}

func On(name string, listener ListenerFunc) {
	event.Std().On(eventName(name), event.ListenerFunc(func(e event.Event) error {
		return listener(e.Data())
	}))
}

func Emit(name string, params model.JSON) error {
	err, _ := event.Std().Fire(eventName(name), params)
	return err
}

func AsyncEmit(name string, params model.JSON) {
	event.Std().FireC(eventName(name), params)
}
