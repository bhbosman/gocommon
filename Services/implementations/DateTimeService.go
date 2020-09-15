package implementations

import (
	"github.com/bhbosman/gocommon/Services/interfaces"
	"go.uber.org/fx"
	"time"
)

type DateTimeService struct {
}

func (self *DateTimeService) Now() time.Time {
	return time.Now()
}

func NewDateTimeService() *DateTimeService {
	return &DateTimeService{}
}

func ProvideDateTimeService() fx.Option {
	return fx.Provide(
		func() (*DateTimeService, interfaces.IDateTimeService) {
			v := NewDateTimeService()
			return v, v
		})
}
