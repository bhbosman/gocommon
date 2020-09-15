package messageRouter

import (
	"github.com/bhbosman/gocommon/constants"
	"reflect"
)

type MessageRouter struct {
	m map[reflect.Type]reflect.Value
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
		return constants.InvalidParam
	}
	if typeOf.NumIn() != 1 {
		return constants.InvalidParam
	}
	paramType := typeOf.In(0)
	self.m[paramType] = reflect.ValueOf(fn)
	return nil
}

func (self *MessageRouter) Route(i interface{}) (bool, error) {
	typeof := reflect.TypeOf(i)
	if function, ok := self.m[typeof]; ok {
		param := reflect.ValueOf(i)
		returnValues := function.Call([]reflect.Value{param})
		for _, returnValue := range returnValues {
			if err, _ := returnValue.Interface().(error); err != nil {
				return true, err
			}
		}
		return true, nil
	}
	return false, nil
}
