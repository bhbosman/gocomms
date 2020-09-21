package endpoints

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/provide"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"html/template"
	"net/http"
)

func RegisterConnectionManagerEndpoint() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Group: "RouteDefinition",
				Target: func(param struct {
					fx.In
					ConnectionManager connectionManager.IObtainConnectionManagerInformation
					Template          *template.Template `name:"Connections.html"`
				}) (*provide.RouteDefinition, error) {
					return provide.NewRouteDefinition("/connections", GetConnections(param.Template, param.ConnectionManager)), nil
				},
			}),
		fx.Provide(
			fx.Annotated{
				Group: "RouteDefinition",
				Target: func(param struct {
					fx.In
					ConnectionManager connectionManager.IConnectionManager
					Template          *template.Template `name:"Connections.html"`
				}) (*provide.RouteDefinition, error) {
					return provide.NewRouteDefinition(
						"/connections/close/id/{id}",
						func(writer http.ResponseWriter, request *http.Request) {
							vars := mux.Vars(request)
							param.ConnectionManager.CloseConnection(vars["id"])
							http.Redirect(writer, request, "/connections", http.StatusSeeOther)
						}), nil
				},
			}),

		fx.Provide(
			fx.Annotated{
				Group: "RouteDefinition",
				Target: func(param struct {
					fx.In
					ConnectionManager connectionManager.ICommandsToConnectionManager
				}) (*provide.RouteDefinition, error) {
					return provide.NewRouteDefinition(
						"/connections/closeAll",
						func(writer http.ResponseWriter, request *http.Request) {
							err := param.ConnectionManager.CloseAllConnections(context.TODO())
							if err != nil {
								writer.WriteHeader(http.StatusInternalServerError)
								return
							}
							http.Redirect(writer, request, "/connections", http.StatusSeeOther)
						}), nil
				},
			}),
	)
}
