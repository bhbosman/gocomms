package common_test

import (
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"testing"
)

func TestCurrentBehaviorWithNoDefinedStacks(t *testing.T) {
	def := common.NewInboundPipeDefinition(nil)

	items := make(chan rxgo.Item)

	def.BuildIncomingObs(items)

}
