package model

type RxHandlerCounter struct {
	Data LineDataEntry
}

func NewRxHandlerCounter(
	stackName string,
	otherMessageCountIn int64,
	rwsMessageCountIn int64,
	otherMessageCountOut int64,
	rwsMessageCountOut int64,
	rwsBytesIn int64,
	rwsBytesOut int64,
) *RxHandlerCounter {
	return &RxHandlerCounter{
		Data: LineDataEntry{
			Name:             stackName,
			OtherMsgCountIn:  otherMessageCountIn,
			RwsMsgCountIn:    rwsMessageCountIn,
			OtherMsgCountOut: otherMessageCountOut,
			RwsMsgCountOut:   rwsMessageCountOut,
			RwsBytesIn:       rwsBytesIn,
			RwsBytesOut:      rwsBytesOut,
		},
	}
}
