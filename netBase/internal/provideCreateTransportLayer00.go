package internal

//func ProvideCreateTransportLayer00() fx.Option {
//	return fx.Provide(
//		fx.Annotated{
//			Target: func(
//				params struct {
//					fx.In
//					TwoWayPipeDefinition common.ITwoWayPipeDefinition
//				},
//			) ([]common.IStackState, error) {
//				return params.TwoWayPipeDefinition.BuildStackState()
//			},
//		},
//	)
//}
