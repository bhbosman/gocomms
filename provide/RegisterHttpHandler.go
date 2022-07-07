package provide

import (
	"context"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocomms/common"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"net/http"
	url2 "net/url"
	"time"
)

func RegisterHttpHandler(url string) fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: func(
					params struct {
						fx.In
						RouteDefinition []*RouteDefinition `group:"RouteDefinition"`
						Logger          *zap.Logger
					},
				) *fx.App {
					fxApp := fx.New(
						fx.StartTimeout(time.Hour),
						fx.StopTimeout(time.Hour),
						fx.Provide(func() *zap.Logger {
							return params.Logger.Named("HttpHandler")
						}),
						fx.WithLogger(
							func(logger *zap.Logger) fxevent.Logger {
								return &fxevent.ZapLogger{Logger: logger}
							},
						),
						fx.Supply(params.RouteDefinition),
						fx.Provide(fx.Annotated{Target: common.CreateUrl(url)}),
						fx.Provide(fx.Annotated{Target: createHttpHandler}),
						fx.Provide(fx.Annotated{Target: createHttpServer}),
						fx.Invoke(
							func(
								params struct {
									fx.In
									Lifecycle         fx.Lifecycle
									Server            *http.Server
									GoFunctionCounter GoFunctionCounter.IService
								},
							) {
								params.Lifecycle.Append(fx.Hook{
									OnStart: func(ctx context.Context) error {

										// this function is part of the GoFunctionCounter count
										go func() {
											functionName := params.GoFunctionCounter.CreateFunctionName("RegisterHttpHandler.OnStart")
											defer func(GoFunctionCounter GoFunctionCounter.IService, name string) {
												_ = GoFunctionCounter.Remove(name)
											}(params.GoFunctionCounter, functionName)
											_ = params.GoFunctionCounter.Add(functionName)

											//
											err := params.Server.ListenAndServe()
											println(err.Error())
										}()
										return nil
									},
									OnStop: func(ctx context.Context) error {
										return params.Server.Close()
									},
								})
							}))
					return fxApp
				},
			}),
	)
}

func createHttpHandler(
	params struct {
		fx.In
		RouteDefinition []*RouteDefinition
	}) http.Handler {
	router := mux.NewRouter()
	for _, def := range params.RouteDefinition {
		router.HandleFunc(def.Path, def.Function)
	}
	return router
}

func createHttpServer(params struct {
	fx.In
	Url     *url2.URL
	Handler http.Handler
}) *http.Server {
	server := &http.Server{
		Addr:    params.Url.Host,
		Handler: params.Handler,
	}
	return server
}
