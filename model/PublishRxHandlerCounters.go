package model

type PublishRxHandlerCounters struct {
	ConnectionId string
	Direction    StreamDirection
	Counters     []*RxHandlerCounter
	Data         map[string]string
}

func NewPublishRxHandlerCounters(connectionId string, direction StreamDirection) *PublishRxHandlerCounters {
	return &PublishRxHandlerCounters{
		ConnectionId: connectionId,
		Direction:    direction,
		Data:         make(map[string]string),
	}
}

func (self *PublishRxHandlerCounters) Add(counter *RxHandlerCounter) {
	self.Counters = append(self.Counters, counter)
}

func (self *PublishRxHandlerCounters) AddMapData(key string, value string) {
	self.Data[key] = value
}
