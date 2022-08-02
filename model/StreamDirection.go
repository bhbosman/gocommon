package model

type StreamDirection uint8

func (self *StreamDirection) String() string {
	switch *self {
	case StreamDirectionInbound:
		return "StreamDirectionInbound"
	case StreamDirectionOutbound:
		return "StreamDirectionOutbound"
	default:
		return "StreamDirectionUnknown"

	}
}

const (
	StreamDirectionInbound StreamDirection = iota
	StreamDirectionOutbound
	StreamDirectionUnknown
)
