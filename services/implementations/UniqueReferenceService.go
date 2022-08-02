package implementations

import (
	"fmt"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"time"
)

type UniqueReferenceService struct {
	UniqueSessionNumber interfaces.IUniqueSessionNumber
}

func (self *UniqueReferenceService) Next(ref string) string {
	return fmt.Sprintf("%v.%v.%v",
		ref,
		time.Now().Format("20060102150405"),
		self.UniqueSessionNumber.Next())
}

func NewUniqueReferenceService(uniqueSessionNumber interfaces.IUniqueSessionNumber) interfaces.IUniqueReferenceService {
	return &UniqueReferenceService{UniqueSessionNumber: uniqueSessionNumber}
}
