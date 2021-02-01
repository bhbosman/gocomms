module github.com/bhbosman/gocomms

go 1.15

require (
	github.com/bhbosman/goMessages v0.0.0-20200922081308-8c8f88094624
	github.com/bhbosman/gocommon v0.0.0-20200921215456-bfddd9bb050e
	github.com/bhbosman/goerrors v0.0.0-20200918064252-e47717b09c4f
	github.com/bhbosman/gologging v0.0.0-20200921180328-d29fc55c00bc
	github.com/bhbosman/gomessageblock v0.0.0-20200921180725-7cd29a998aa3
	github.com/bhbosman/goprotoextra v0.0.1
	github.com/bhbosman/gorxextra v0.0.0-20200918070301-48dbd8b934dc
	github.com/cskr/pubsub v1.0.2
	github.com/gobwas/httphead v0.0.0-20200921212729-da3d93bc3c58 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.0.4
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.8.0
	github.com/icza/gox v0.0.0-20200702115100-7dc3510ae515
	github.com/reactivex/rxgo/v2 v2.1.0
	github.com/stretchr/testify v1.6.1
	go.uber.org/fx v1.13.1
	go.uber.org/multierr v1.6.0
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	google.golang.org/protobuf v1.25.0

)
replace github.com/reactivex/rxgo/v2 v2.1.0 => github.com/bhbosman/rxgo/v2 v2.1.1-0.20200922152528-6aef42e76e00


replace (
		github.com/bhbosman/goMessages => ../goMessages
    	github.com/bhbosman/gocommon => ../gocommon
    	github.com/bhbosman/goerrors => ../goerrors
    	github.com/bhbosman/gologging => ../gologging
    	github.com/bhbosman/gomessageblock => ../gomessageblock
    	github.com/bhbosman/goprotoextra => ../goprotoextra
    	github.com/bhbosman/gorxextra => ../gorxextra

)