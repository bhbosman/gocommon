package implementations

import (
	"fmt"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"go.uber.org/fx"
	"time"
)

type UniqueReferenceService struct {
	UniqueSessionNumber interfaces.IUniqueSessionNumber
}

func (self *UniqueReferenceService) Next(ref string) string {
	return fmt.Sprintf("%v.%v.%-120d",
		ref,
		time.Now().Format("20060102150405"),
		self.UniqueSessionNumber.Next())
}

func NewUniqueReferenceService(uniqueSessionNumber interfaces.IUniqueSessionNumber) *UniqueReferenceService {
	return &UniqueReferenceService{UniqueSessionNumber: uniqueSessionNumber}
}

func ProvideNewUniqueReferenceService() fx.Option {
	return fx.Provide(
		func(uniqueSessionNumber interfaces.IUniqueSessionNumber) (*UniqueReferenceService, interfaces.IUniqueReferenceService) {
			v := NewUniqueReferenceService(uniqueSessionNumber)
			return v, v
		})
}
