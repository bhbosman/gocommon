package model

type PublishRxHandlerCounters struct {
	ConnectionId string
	Direction    StreamDirection
	Counters     []*RxHandlerCounter
}

func NewPublishRxHandlerCounters(connectionId string, direction StreamDirection) *PublishRxHandlerCounters {
	return &PublishRxHandlerCounters{
		ConnectionId: connectionId,
		Direction:    direction,
	}
}

func (self *PublishRxHandlerCounters) Add(counter *RxHandlerCounter) {
	self.Counters = append(self.Counters, counter)
}
