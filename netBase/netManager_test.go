package netBase_test

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/common"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/netBase"
	"github.com/bhbosman/goprotoextra"
	"github.com/golang/mock/gomock"
	"github.com/reactivex/rxgo/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io"
	"net/url"
	"testing"
)

func TestNetManager_NewConnectionInstance(t *testing.T) {
	createSut := func(
		controller *gomock.Controller,
		crf intf.IConnectionReactorFactory,
		stackName string,
		connectionManager goConnectionManager.IService,
	) (netBase.ConnNetManager, error) {
		urlPath, _ := url.Parse("http://localhost:8080")
		cancelContext, cancelFunc := context.WithCancel(context.Background())

		uniqueSessionNumber := interfaces.NewMockIUniqueReferenceService(controller)
		return netBase.NewConnNetManager(
			"",
			"",
			false, nil,
			urlPath,
			cancelContext,
			cancelFunc,
			stackName,
			connectionManager,
			nil,
			zap.NewNop(),
			uniqueSessionNumber,
			func() fx.Option {
				return fx.Options()
			}, nil)
	}

	t.Run("Test Connection with no connection stack given", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		connectionManagerStarted := goConnectionManager.NewMockIService(controller)
		connectionManagerStarted.EXPECT().OnStateReturn(IFxService.Started)

		connectionReactor := intf.NewMockIConnectionReactor(controller)
		connectionReactor.EXPECT().OnInitDoAndReturnDefault(
			func(arg0 goprotoextra.ToConnectionFunc, arg1 goprotoextra.ToReactorFunc, arg2, arg3 rxgo.NextFunc) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
				return func(i interface{}) {

					},
					func(err error) {

					},
					func() {

					},
					nil
			},
		).AnyTimes()
		connectionReactor.EXPECT().OnOpenReturn(nil).AnyTimes()
		connectionReactor.EXPECT().OnCloseReturn(nil).AnyTimes()
		crf := intf.NewMockIConnectionReactorFactory(controller)
		crf.EXPECT().
			OnCreateReturnDefault(connectionReactor, nil).
			AnyTimes()

		sut, err := createSut(
			controller,
			crf,
			"", /* testing this */
			connectionManagerStarted /* not testing this */)
		assert.NoError(t, err)
		addr := common.NewMockAddr(controller)
		addr.EXPECT().OnStringReturn("ABC").AnyTimes()

		conn := common.NewMockConn(controller)
		conn.EXPECT().OnLocalAddrReturn(addr).AnyTimes()
		conn.EXPECT().OnRemoteAddrReturn(addr).AnyTimes()

		_, ctx, _ := sut.NewConnectionInstance("ABC", nil, model.ClientConnection, conn)
		assert.Error(t, ctx.Err())
	})
	t.Run("Check if connection manager is setup", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		connectionManagerNotStarted := goConnectionManager.NewMockIService(controller)
		connectionManagerNotStarted.EXPECT().OnStateReturn(IFxService.NotInitialized).AnyTimes()
		connectionManagerNotStarted.EXPECT().OnServiceNameReturn("Some service")
		connectionReactor := intf.NewMockIConnectionReactor(controller)
		connectionReactor.EXPECT().OnOpenReturn(nil).AnyTimes()
		connectionReactor.EXPECT().OnInitDoAndReturnDefault(
			func(arg0 goprotoextra.ToConnectionFunc, arg1 goprotoextra.ToReactorFunc, arg2, arg3 rxgo.NextFunc) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
				return func(i interface{}) {

					},
					func(err error) {

					},
					func() {

					}, nil
			},
		).AnyTimes()
		connectionReactor.EXPECT().OnCloseReturn(nil).AnyTimes()
		crf := intf.NewMockIConnectionReactorFactory(controller)
		crf.EXPECT().
			OnCreateReturnDefault(connectionReactor, nil).
			AnyTimes()

		sut, err := createSut(
			controller,
			crf,
			"dd", //goCommsDefinitions.TransportFactoryCompressedTlsName, /* not testing this */
			connectionManagerNotStarted /* testing this */)
		assert.NoError(t, err)
		addr := common.NewMockAddr(controller)
		addr.EXPECT().OnStringReturn("ABC").AnyTimes()

		conn := common.NewMockConn(controller)
		conn.EXPECT().OnLocalAddrReturn(addr).AnyTimes()
		conn.EXPECT().OnRemoteAddrReturn(addr).AnyTimes()

		instance, ctx, cancelFunc := sut.NewConnectionInstance("ABC", nil, model.ClientConnection, conn)
		assert.NotNil(t, ctx)
		assert.NotNil(t, cancelFunc)
		assert.NotNil(t, instance)
		assert.Error(t, instance.Err())

	})

	t.Run("Connection manager exist. Check registration and deRegistration, connection closes at first read", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		connectionManagerStarted := goConnectionManager.NewMockIService(controller)
		connectionManagerStarted.EXPECT().OnStateReturn(IFxService.Started).AnyTimes()
		connectionManagerStarted.EXPECT().OnServiceNameReturn("Some service").AnyTimes()
		connectionManagerStarted.EXPECT().OnRegisterConnectionReturn(gomock.Any(), gomock.Any(), gomock.Any(), nil)
		connectionManagerStarted.EXPECT().OnDeregisterConnectionReturn(gomock.Any(), nil)
		connectionManagerStarted.EXPECT().OnConnectionInformationReceivedReturn(gomock.Any(), nil).AnyTimes()
		connectionManagerStarted.EXPECT().OnConnectionInformationReceivedReturn(gomock.Any(), nil).AnyTimes()

		connectionReactor := intf.NewMockIConnectionReactor(controller)
		connectionReactor.EXPECT().OnInitDoAndReturnDefault(
			func(arg0 goprotoextra.ToConnectionFunc, arg1 goprotoextra.ToReactorFunc, arg2, arg3 rxgo.NextFunc) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
				return func(i interface{}) {

					},
					func(err error) {

					},
					func() {

					}, nil
			},
		).AnyTimes()
		connectionReactor.EXPECT().OnOpenReturn(nil).AnyTimes()
		connectionReactor.EXPECT().OnCloseReturn(nil).AnyTimes()
		crf := intf.NewMockIConnectionReactorFactory(controller)
		crf.EXPECT().
			OnCreateReturnDefault(connectionReactor, nil).
			AnyTimes()

		sut, err := createSut(
			controller,
			crf,
			goCommsDefinitions.EmptyStackName, /* not testing this */
			connectionManagerStarted /* testing this */)
		assert.NoError(t, err)
		addr := common.NewMockAddr(controller)
		addr.EXPECT().OnStringReturn("ABC").AnyTimes()

		dataChannel := make(chan goprotoextra.ReadWriterSize)
		conn := common.NewMockConn(controller)
		conn.EXPECT().OnLocalAddrReturn(addr).AnyTimes()
		conn.EXPECT().OnRemoteAddrReturn(addr).AnyTimes()
		conn.EXPECT().OnReadDoAndReturn(
			gomock.Any(),
			func(arg0 []byte) (int, error) {
				data, ok := <-dataChannel
				if ok {
					return data.Read(arg0)
				}
				return 0, io.EOF
			})
		conn.EXPECT().OnCloseDoAndReturn(func() error {
			return nil
		})
		conn.EXPECT().OnWriteDoAndReturn(
			gomock.Any(),
			func(arg0 []byte) (int, error) {
				return 0, io.EOF
			}).AnyTimes()

		instance, ctx, cancelFunc := sut.NewConnectionInstance("ABC", nil, model.ClientConnection, conn)
		assert.NotNil(t, ctx)
		assert.NotNil(t, cancelFunc)
		assert.NotNil(t, instance)
		assert.NoError(t, instance.Err())
		assert.NoError(t, instance.Start(context.Background()))
		close(dataChannel)
		err = instance.Stop(context.Background())
		assert.NoError(t, err)
	})
}
