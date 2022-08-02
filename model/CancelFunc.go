package model

type ConnectionCancelFunc func(context string, inbound bool, err error)
