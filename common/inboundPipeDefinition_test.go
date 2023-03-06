package common_test

import (
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func TestCurrentBehaviorWithNoDefinedStacks(t *testing.T) {
	def := common.NewInboundPipeDefinition(nil)
	items := make(chan rxgo.Item)
	defer close(items)
	m := make(map[string]*common.StackDataContainer)

	obs, err := def.BuildIncomingObs(items, m, context.Background())
	if !assert.NoError(t, err) {
		return
	}
	obs.InboundObservable.Observe()

}
