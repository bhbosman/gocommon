package model

type PublishRxHandlerCounters struct {
	ConnectionId string
	Direction    StreamDirection
	Counters     []*RxHandlerCounter
	data         map[string]string
}

func NewPublishRxHandlerCounters(connectionId string, direction StreamDirection) *PublishRxHandlerCounters {
	return &PublishRxHandlerCounters{
		ConnectionId: connectionId,
		Direction:    direction,
		data:         make(map[string]string),
	}
}

func (self *PublishRxHandlerCounters) Add(counter *RxHandlerCounter) {
	self.Counters = append(self.Counters, counter)
}

func (self *PublishRxHandlerCounters) AddMapData(key string, value string) {
	self.data[key] = value
}
