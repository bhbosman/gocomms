package internal_test

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"testing"
)

//go:generate mockgen -package internal_test -destination stackDefinition_mock_test.go . IStackDefinition,IBoundResult
//go:generate mockgen -package internal_test -destination connectionManager_mock_test.go github.com/bhbosman/gocomms/connectionManager IConnectionManager
//go:generate mockgen -package internal_test -destination net_mock_test.go net Conn
func TestName(t *testing.T) {
	//id := uuid.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	stackDefinition := NewMockIStackDefinition(ctrl)

	inboundResult := NewMockIBoundResult(ctrl)
	inboundResult.EXPECT().GetBoundResult().Return(
		func(stackData, pipeData interface{}, params internal.InOutBoundParams) internal.IStackBoundDefinition {
			return internal.NewBoundDefinition(
				func(params internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
					return stackDefinition.GetId(), nil, nil
				},
				nil)
		}).Times(1)

	outboundResult := NewMockIBoundResult(ctrl)
	outboundResult.EXPECT().GetBoundResult().Return(
		func(stackData, pipeData interface{}, params internal.InOutBoundParams) internal.IStackBoundDefinition {
			return internal.NewBoundDefinition(
				func(params internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
					return stackDefinition.GetId(), nil, nil
				},
				nil)
		}).Times(1)

	stackDefinition.EXPECT().GetId().Return(uuid.New()).AnyTimes()
	//stackDefinition.EXPECT().GetName().Return("TestStack").Times(1)
	stackDefinition.EXPECT().GetInbound().Return(inboundResult).Times(1)
	stackDefinition.EXPECT().GetOutbound().Return(outboundResult).AnyTimes()
	stackDefinition.EXPECT().GetStackState().Return(internal.StackState{}).Times(1)

	def := internal.NewTwoWayPipeDefinition()
	def.AddStackDefinitionFunc(func() (internal.IStackDefinition, error) {
		return stackDefinition, nil
	})

	connectionManagerMock := NewMockIConnectionManager(ctrl)
	connectionMock := NewMockConn(ctrl)

	def.Build(
		"test",
		connectionManagerMock,
		connectionMock,
		context.Background(),
		func(context string, inbound bool, err error) {

		})
}
