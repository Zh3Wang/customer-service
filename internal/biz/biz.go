package biz

type MessageType int

const (
	Common MessageType = iota
	Heartbeat
	Broadcast
	Online
	Offline
)