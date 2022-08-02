package GoFunctionCounter

import (
	"github.com/bhbosman/gocommon/Services/IDataShutDown"
	"github.com/bhbosman/gocommon/Services/IFxService"
)

type IGoFunctionCounter interface {
	Remove(name string) error
	Add(name string) error
}

type IUi interface {
	SetConnectionListChange(cb func(names []string))
}

type IData interface {
	IDataShutDown.IDataShutDown
	IGoFunctionCounter
	IUi
}
type IService interface {
	IGoFunctionCounter
	IFxService.IFxServices
	IUi

	// CreateFunctionName(s string) string

	GoRun(s string, cb func()) error
}
