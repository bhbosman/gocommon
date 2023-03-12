package GoFunctionCounter

import (
	"github.com/bhbosman/gocommon/messageRouter"
	"github.com/bhbosman/gocommon/messages"
	"sort"
)

type data struct {
	emptyQueueReceived bool
	dirty              bool
	m                  map[string]bool
	messageRouter      messageRouter.IMessageRouter
	cb                 func(names []string)
}

func (self *data) SetConnectionListChange(cb func(names []string)) {
	self.cb = cb
}

func (self *data) Remove(name string) error {
	delete(self.m, name)
	self.dirty = true
	return nil
}

func (self *data) Add(name string) error {
	if _, ok := self.m[name]; !ok {
		self.m[name] = true
		self.dirty = true
		return nil
	}
	panic("make up a unique number")
}

func (self *data) Send(message interface{}) error {
	self.messageRouter.Route(message)
	return nil
}

func (self *data) ShutDown() error {
	return nil
}

func (self *data) handleTimerReceived(_ *timerReceived) {
	if self.emptyQueueReceived {
		self.emptyQueueReceived = false
		if self.dirty {
			self.dirty = false
			if self.cb != nil {
				ss := make([]string, 0, len(self.m))
				for k := range self.m {
					ss = append(ss, k)
				}
				sort.Strings(ss)
				self.cb(ss)
			}
		}
	}
}

func (self *data) handleEmptyQueue(*messages.EmptyQueue) {
	self.emptyQueueReceived = true
}
func newData() (IData, error) {
	result := &data{
		messageRouter: messageRouter.NewMessageRouter(),
		m:             make(map[string]bool),
	}

	_ = result.messageRouter.Add(result.handleTimerReceived)
	_ = result.messageRouter.Add(result.handleEmptyQueue)
	return result, nil
}
