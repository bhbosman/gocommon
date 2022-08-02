package model

import "context"

type IRegisterToConnectionManager interface {
	RegisterConnection(id string, function context.CancelFunc, CancelContext context.Context) error
	DeregisterConnection(id string) error
	NameConnection(id string, name string) error
}
