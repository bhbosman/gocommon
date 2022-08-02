package IDateTimeService

import "time"

type IDateTimeService interface {
	Now() time.Time
}
