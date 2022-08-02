package model

type ConnectionType uint8

const (
	ServerConnection ConnectionType = iota
	ClientConnection
)
