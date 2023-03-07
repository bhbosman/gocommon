package gocommon

import "github.com/reactivex/rxgo/v2"

type IObservable interface {
	Observe(opts ...rxgo.Option) <-chan rxgo.Item
	Map(apply rxgo.Func, opts ...rxgo.Option) rxgo.Observable
	FlatMap(apply rxgo.ItemToObservable, opts ...rxgo.Option) rxgo.Observable
	//rxgo.Observable
}
