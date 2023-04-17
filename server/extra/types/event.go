package types

type GroupEvent int

const (
	GroupEventUnknown GroupEvent = iota
	GroupEventJoin
	GroupEventExit
	GroupEventOnline
	GroupEventOffline
	GroupEventReceive
)
