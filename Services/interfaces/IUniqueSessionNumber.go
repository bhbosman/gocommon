package interfaces

import "time"

type IUniqueSessionNumber interface {
	Next() uint64
	NextUin32() uint32
}

type IUniqueReferenceService interface {
	Next(ref string) string
}

type IDateTimeService interface {
	Now() time.Time
}
