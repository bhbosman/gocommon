package IFxService

import (
	"context"
	"fmt"
)

type State int

const (
	NotInitialized State = 0
	Started        State = 1
	Stopped        State = 2
)

type IFxServices interface {
	OnStart(ctx context.Context) error
	OnStop(ctx context.Context) error
	State() State
	ServiceName() string
}

type ServiceStateError struct {
	ServiceName string
	Message     string
	Wanted      State
	Actual      State
}

func NewServiceStateError(serviceName string, message string, wanted State, actual State) *ServiceStateError {
	return &ServiceStateError{ServiceName: serviceName, Message: message, Wanted: wanted, Actual: actual}
}

func (self *ServiceStateError) Error() string {
	return fmt.Sprintf("Service %v state is not correct. Wanted: %v, Actual: %v. Additional message: %v",
		self.ServiceName, self.Wanted, self.Actual, self.Message)
}
