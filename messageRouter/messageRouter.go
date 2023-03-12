package messageRouter

import (
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"reflect"
)

type IMessageRouter interface {
	Add(fn interface{}) error
	Route(i interface{})
	MultiRoute(messages ...interface{})
	RegisterUnknown(unknown rxgo.NextFunc)
}

type messageRouter struct {
	m       map[reflect.Type]reflect.Value
	unknown rxgo.NextFunc
}

func NewMessageRouter() IMessageRouter {
	m := make(map[reflect.Type]reflect.Value)
	return &messageRouter{
		m: m,
	}
}

func (self *messageRouter) Add(fn interface{}) error {
	typeOf := reflect.TypeOf(fn)
	if typeOf.Kind() != reflect.Func {
		return goerrors.NewInvalidParamError("fn", "The parameter must be a function")
	}
	if typeOf.NumIn() != 1 {
		return goerrors.NewInvalidParamError("fn", "The incoming function must have 1 parameter")
	}
	paramType := typeOf.In(0)
	self.m[paramType] = reflect.ValueOf(fn)
	return nil
}
func (self *messageRouter) RegisterUnknown(unknown rxgo.NextFunc) {
	self.unknown = unknown
}
func (self *messageRouter) Route(i interface{}) {
	typeof := reflect.TypeOf(i)
	if function, ok := self.m[typeof]; ok {
		param := reflect.ValueOf(i)
		returnValues := function.Call([]reflect.Value{param})
		for _, returnValue := range returnValues {
			if err, _ := returnValue.Interface().(error); err != nil {
				return
			}
		}
		return
	}
	if self.unknown != nil {
		self.unknown(i)
	}
	return
}

func (self *messageRouter) MultiRoute(messages ...interface{}) {
	for _, message := range messages {
		self.Route(message)
	}
}
