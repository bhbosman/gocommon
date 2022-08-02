package model

type LineDataEntry struct {
	Name          string
	OtherMsgCount int64
	RwsMsgCount   int64
	RwsBytesIn    int64
	RwsBytesOut   int64
}

type LineData struct {
	InValue  LineDataEntry
	OutValue LineDataEntry
}
