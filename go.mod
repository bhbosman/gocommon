module github.com/bhbosman/gocommon

go 1.15

require (
	github.com/bhbosman/goerrors v0.0.0-20200918064252-e47717b09c4f
	github.com/bhbosman/gologging v0.0.0-20200921180328-d29fc55c00bc
	github.com/bhbosman/gomessageblock v0.0.0-20200921180725-7cd29a998aa3
	github.com/bhbosman/goprotoextra v0.0.1
	github.com/cskr/pubsub v1.0.2
	go.uber.org/fx v1.13.1
	google.golang.org/protobuf v1.25.0
)

//replace github.com/reactivex/rxgo/v2 v2.1.0 => github.com/bhbosman/rxgo/v2 v2.1.1-0.20200922152528-6aef42e76e00

//replace github.com/bhbosman/gocomms => ../gocomms
