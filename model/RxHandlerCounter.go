package model

type RxHandlerCounter struct {
	Data LineDataEntry
}

func NewRxHandlerCounter(
	stackName string,
	otherMessageCount int64,
	rwsMessageCount int64,
	rwsBytesIn int64,
	rwsBytesOut int64,
) *RxHandlerCounter {
	return &RxHandlerCounter{
		Data: LineDataEntry{
			Name:          stackName,
			OtherMsgCount: otherMessageCount,
			RwsMsgCount:   rwsMessageCount,
			RwsBytesIn:    rwsBytesIn,
			RwsBytesOut:   rwsBytesOut,
		},
	}
}
