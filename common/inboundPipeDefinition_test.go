package common_test

//func TestCurrentBehaviorWithNoDefinedStacks(t *testing.T) {
//	mockController := gomock.NewController(t)
//	defer mockController.Finish()
//
//	bottomStack := common.NewMockIInboundData(mockController)
//	bottomStack.EXPECT().Name().Return("BottomStack").AnyTimes()
//	def := common.NewInboundPipeDefinition(
//		[]common.IInboundData{
//			bottomStack,
//		},
//	)
//	items := make(chan rxgo.Item)
//	defer close(items)
//	m := make(map[string]*common.StackDataContainer)
//
//	obs, err := def.BuildIncomingObs(items, m, context.Background())
//	if !assert.NoError(t, err) {
//		return
//	}
//	obs.Observe()
//}
