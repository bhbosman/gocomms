package impl

import (
	"go.uber.org/fx"
)

func RegisterAllConnectionRelatedServices() fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotated{Target: newConnectionReactorFactories}),
	)
}
