package messageRouter

import (
	"github.com/bhbosman/goerrors"
	"github.com/reactivex/rxgo/v2"
	"reflect"
)

type MessageRouter struct {
	m       map[reflect.Type]reflect.Value
	unknown rxgo.NextFunc
}

func NewMessageRouter() *MessageRouter {
	m := make(map[reflect.Type]reflect.Value)
	return &MessageRouter{
		m: m,
	}
}

func (self *MessageRouter) Add(fn interface{}) error {
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
func (self *MessageRouter) RegisterUnknown(unknown rxgo.NextFunc) {
	self.unknown = unknown
}
func (self *MessageRouter) Route(i interface{}) {
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

func (self *MessageRouter) MultiRoute(messages ...interface{}) {
	for _, message := range messages {
		self.Route(message)
	}
}
