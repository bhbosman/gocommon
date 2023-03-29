package gocommon

import (
	"go.uber.org/multierr"
	"sync"
)

func RegisterConnectionShutdown(connectionId string, callback func(), cancellationContext ...ICancellationContext) error {
	mutex := sync.Mutex{}
	cancelCalled := false
	cb := func() {
		mutex.Lock()
		b := cancelCalled
		cancelCalled = true
		mutex.Unlock()
		if !b {
			callback()
		}
		for _, instance := range cancellationContext {
			_ = instance.Remove(connectionId)
		}
	}
	var result error
	for _, ctx := range cancellationContext {
		b, err := ctx.Add(connectionId, cb)
		result = multierr.Append(result, err)
		if !b {
			cb()
			return result
		}
	}
	return result
}
