module github.com/bhbosman/gocommon

go 1.15

require github.com/bhbosman/goerrors v0.0.0-20200918064252-e47717b09c4f
require github.com/bhbosman/gologging v0.0.0-20200921180328-d29fc55c00bc
require github.com/bhbosman/gomessageblock v0.0.0-20200921180725-7cd29a998aa3
require github.com/bhbosman/goprotoextra v0.0.2-0.20210414124526-a342e2a9e82f
require github.com/cskr/pubsub v1.0.2
require go.uber.org/fx v1.13.1
require google.golang.org/protobuf v1.25.0

replace github.com/bhbosman/gologging => ../gologging
replace github.com/bhbosman/gomessageblock => ../gomessageblock
replace github.com/bhbosman/goerrors => ../goerrors
