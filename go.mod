module github.com/bhbosman/gocomms

go 1.18

require (
	github.com/bhbosman/goCommsDefinitions v0.0.0-20220707044904-ceb8c2737904
	github.com/bhbosman/goConnectionManager v0.0.0-20220705103338-3f5a18784e60
	github.com/bhbosman/gocommon v0.0.0-20220707045107-1b6a40e49fd5
	github.com/bhbosman/goerrors v0.0.0-20220623084908-4d7bbcd178cf
	github.com/bhbosman/gomessageblock v0.0.0-20220617132215-32f430d7de62
	github.com/bhbosman/goprotoextra v0.0.2-0.20210817141206-117becbef7c7
	github.com/golang/mock v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/reactivex/rxgo/v2 v2.5.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/fx v1.17.1
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.21.0
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f
)

require (
	github.com/cenkalti/backoff/v4 v4.0.0 // indirect
	github.com/cskr/pubsub v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/icza/gox v0.0.0-20220321141217-e2d488ab2fbc // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.1.0 // indirect
	github.com/teivah/onecontext v0.0.0-20200513185103-40f981bfd775 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/dig v1.14.0 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20220617134815-f277ff266f47

replace github.com/bhbosman/goCommsDefinitions => ../goCommsDefinitions

//replace github.com/bhbosman/goerrors => ../goerrors

replace github.com/bhbosman/goConnectionManager => ../goConnectionManager

replace github.com/bhbosman/goprotoextra => ../goprotoextra

replace github.com/bhbosman/gocommon => ../gocommon
