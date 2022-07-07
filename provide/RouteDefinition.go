package provide

import "net/http"

type RouteDefinition struct {
	Path     string
	Function func(http.ResponseWriter, *http.Request)
}

func NewRouteDefinition(path string, function func(http.ResponseWriter, *http.Request)) *RouteDefinition {
	return &RouteDefinition{Path: path, Function: function}
}
