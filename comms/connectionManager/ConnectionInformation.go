package connectionManager

import (
	"context"
	rxgo "github.com/ReactiveX/RxGo"
	"sync"
	"time"
)

type StackPropertyKey struct {
	Index     int
	Name      string
	Direction rxgo.StreamDirection
}

type StackPropertyValue struct {
	MsgCount  int
	ByteCount int
}

type ConnectionInformation struct {
	Id              string
	CancelFunc      context.CancelFunc
	Name            string
	Status          string
	ConnectionTime  time.Time
	StackProperties map[StackPropertyKey]StackPropertyValue
	mutex           sync.Mutex
}

func NewConnectionInformation(id string, function context.CancelFunc) *ConnectionInformation {

	return &ConnectionInformation{
		Id:              id,
		CancelFunc:      function,
		StackProperties: make(map[StackPropertyKey]StackPropertyValue),
		ConnectionTime:  time.Now().Truncate(time.Second),
	}
}
