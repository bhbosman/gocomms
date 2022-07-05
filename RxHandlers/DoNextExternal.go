package RxHandlers

//
//func DoNextExternal(
//	ctx context.Context,
//	cancelFunc context.CancelFunc,
//	obs rxgo.Observable,
//	nextFunc intf.NextExternalFunc,
//	logger *zap.Logger,
//	goFunctionCounter GoFunctionCounter.IService,
//	opts ...rxgo.Option,
//) rxgo.Disposed {
//	dispose := make(chan struct{})
//	handler := func(ctx context.Context, src <-chan rxgo.Item) {
//		defer close(dispose)
//		count := 0
//		for {
//			select {
//			case <-ctx.Done():
//				return
//			case i, ok := <-src:
//				if !ok {
//					continue
//				}
//				if i.Error() {
//					cancelFunc()
//					logger.Error("Error propagated through stack", zap.Error(i.E))
//					continue
//				}
//				if nextExternal, ok := i.V.(*NextExternal); ok {
//					nextFunc(nextExternal.External, nextExternal.Data)
//					count++
//					if len(src) == 0 {
//						nextFunc(false, &messages.EmptyQueue{
//							Count: count,
//						})
//						count = 0
//					}
//					continue
//				}
//			}
//		}
//	}
//
//	// this function is part of the GoFunctionCounter count
//	go func() {
//		functionName := goFunctionCounter.CreateFunctionName("DoNextExternal")
//		defer func(GoFunctionCounter GoFunctionCounter.IService, name string) {
//			_ = GoFunctionCounter.Remove(name)
//		}(goFunctionCounter, functionName)
//		_ = goFunctionCounter.Add(functionName)
//
//		//
//		handler(ctx, obs.Observe(opts...))
//	}()
//
//	return dispose
//}
