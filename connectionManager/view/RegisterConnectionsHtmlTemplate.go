package view

import (
	"go.uber.org/fx"
	"html/template"
)

func RegisterConnectionsHtmlTemplate() fx.Option {
	const connectionsTemplateName = "Connections.html"
	return fx.Options(
		fx.Provide(fx.Annotated{
			Name: connectionsTemplateName,
			Target: func() (*template.Template, error) {
				bytes, err := Asset(connectionsTemplateName)
				if err != nil {
					return nil, err
				}
				return template.New(connectionsTemplateName).Parse(string(bytes))
			},
		}))
}
