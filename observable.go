package gocommon

import "github.com/reactivex/rxgo/v2"

type IObservable interface {
	Map(apply rxgo.Func, opts ...rxgo.Option) rxgo.Observable
	rxgo.Observable
}
