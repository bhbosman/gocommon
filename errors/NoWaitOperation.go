package errors

type noWaitOperation struct {
}

func (self *noWaitOperation) Error() string {
	return "No wait operation"
}

var NoWaitOperationError error = &noWaitOperation{}
