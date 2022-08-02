package implementations

import (
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
