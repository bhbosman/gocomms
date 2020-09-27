package provide

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"net/http"
	url2 "net/url"
)

func RegisterHttpHandler(url string) fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "Apps",
				Target: func(params struct {
					fx.In
					RouteDefinition []*RouteDefinition `group:"RouteDefinition"`
				}) *fx.App {
					fxApp := fx.New(
						fx.Supply(params.RouteDefinition),
						fx.Provide(fx.Annotated{Target: internal.CreateUrl(url)}),
						fx.Provide(fx.Annotated{Target: createHttpHandler}),
						fx.Provide(fx.Annotated{Target: createHttpServer}),
						fx.Invoke(func(params struct {
							fx.In
							Lifecycle fx.Lifecycle
							Server    *http.Server
						}) {
							params.Lifecycle.Append(fx.Hook{
								OnStart: func(ctx context.Context) error {
									go func() {
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
