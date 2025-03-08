module github.com/bhbosman/gocomms

go 1.24.0

require (
	github.com/bhbosman/goCommsDefinitions v0.0.0-20250308144130-64993b60920c
	github.com/bhbosman/goConnectionManager v0.0.0-20250308194756-b3c21d488492
	github.com/bhbosman/gocommon v0.0.0-20250308194442-9c45d7859806
	github.com/bhbosman/goerrors v0.0.0-20250307194237-312d070c8e38
	github.com/bhbosman/gomessageblock v0.0.0-20250308073733-0b3daca12e3a
	github.com/bhbosman/goprotoextra v0.0.2
	github.com/golang/mock v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/reactivex/rxgo/v2 v2.5.0
	github.com/stretchr/testify v1.10.0
	go.uber.org/fx v1.23.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
	golang.org/x/net v0.37.0
)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/icza/gox v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/teivah/onecontext v1.3.0 // indirect
	go.uber.org/dig v1.18.1 // indirect
	golang.org/x/sys v0.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require github.com/cskr/pubsub v1.0.2

replace (
	github.com/cskr/pubsub => github.com/bhbosman/pubsub v1.0.3-0.20250308211253-e47ee222bb15
	github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20250308071159-4cf72f668c72
)
