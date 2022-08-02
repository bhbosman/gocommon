package interfaces

type IUniqueSessionNumber interface {
	Next() uint64
	NextUin32() uint32
}
