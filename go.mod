module github.com/bhbosman/gocommon

go 1.18

require (
	github.com/bhbosman/goCommsDefinitions v0.0.0-20220801175552-c5aa68065af3
	github.com/bhbosman/goerrors v0.0.0-20220623084908-4d7bbcd178cf
	github.com/bhbosman/gomessageblock v0.0.0-20220617132215-32f430d7de62
	github.com/bhbosman/goprotoextra v0.0.0-20220707073310-8b1052b95990
	github.com/cskr/pubsub v1.0.2
	github.com/golang/mock v1.6.0
	github.com/icza/gox v0.0.0-20220321141217-e2d488ab2fbc
	github.com/reactivex/rxgo/v2 v2.5.0
	go.uber.org/fx v1.18.2
	go.uber.org/zap v1.21.0
	golang.org/x/net v0.7.0
	google.golang.org/protobuf v1.25.0
)

require (
	github.com/cenkalti/backoff/v4 v4.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.1.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	github.com/teivah/onecontext v0.0.0-20200513185103-40f981bfd775 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/dig v1.15.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20230302060806-d02c40b7514e

replace github.com/bhbosman/goCommsDefinitions => ../goCommsDefinitions

//replace github.com/bhbosman/goerrors => ../goerrors

//replace github.com/bhbosman/goprotoextra => ../goprotoextra
replace github.com/cskr/pubsub => ../pubsub
