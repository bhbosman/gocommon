package implementations

import (
	"github.com/bhbosman/gocommon/Services/interfaces"
	"go.uber.org/fx"
	"sync/atomic"
)

type UniqueSessionNumber struct {
	count uint64
}

func (n *UniqueSessionNumber) NextUin32() uint32 {
	return uint32(n.Next())
}

func (n *UniqueSessionNumber) Next() uint64 {
	return atomic.AddUint64(&n.count, 1)
}

func NewUniqueSessionNumber() *UniqueSessionNumber {
	return &UniqueSessionNumber{}
}

func ProvideUniqueSessionNumber() fx.Option {
	return fx.Provide(
		func() (*UniqueSessionNumber, interfaces.IUniqueSessionNumber) {
			v := NewUniqueSessionNumber()
			return v, v
		})
}
