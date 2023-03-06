package common_test

import (
	"github.com/bhbosman/gocomms/common"
	"github.com/golang/mock/gomock"
	"github.com/reactivex/rxgo/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func TestCurrentBehaviorWithNoDefinedStacks(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	bottomStack := common.NewMockIInboundData(mockController)
	bottomStack.EXPECT().Name().Return("BottomStack").AnyTimes()
	def := common.NewInboundPipeDefinition(
		[]common.IInboundData{
			bottomStack,
		},
	)
	items := make(chan rxgo.Item)
	defer close(items)
	m := make(map[string]*common.StackDataContainer)

	obs, err := def.BuildIncomingObs(items, m, context.Background())
	if !assert.NoError(t, err) {
		return
	}
	obs.Observe()
}
