package model

type LineDataEntry struct {
	Name             string
	OtherMsgCountIn  int64
	RwsMsgCountIn    int64
	OtherMsgCountOut int64
	RwsMsgCountOut   int64
	RwsBytesIn       int64
	RwsBytesOut      int64
}

type LineData struct {
	InValue  LineDataEntry
	OutValue LineDataEntry
}

type KeyValue struct {
	Key   string
	Value string
}
