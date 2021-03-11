package crMessageBreaker_test

import (
	"context"
	internal2 "github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gomessageblock"
	"github.com/reactivex/rxgo/v2"
	"github.com/stretchr/testify/assert"
	"sync"

	"github.com/bhbosman/gocomms/stacks/crMessageBreaker"
	"testing"
)

func TestMessageBreaker(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		definition, err := crMessageBreaker.StackDefinition(
			internal2.ServerConnection,
			"",
			func(context string, inbound bool, err error) {

			},
			nil)
		assert.NoError(t, err)
		cancelContext := context.Background()
		stackCancelFunc := func(context string, inbound bool, err error) {

		}
		bottomChannel := make(chan rxgo.Item)
		defer close(bottomChannel)
		bottomObservable := rxgo.FromChannel(bottomChannel)

		inbound := definition.Inbound(
			internal2.InOutBoundParams{
				Index:   0,
				Context: cancelContext,
			})

		crObservable, err := inbound.PipeDefinition(
			internal2.PipeDefinitionParams{
				ConnectionId:      "",
				ConnectionManager: nil,
				CancelContext:     cancelContext,
				StackCancelFunc:   stackCancelFunc,
				Obs:               bottomObservable,
			})
		wg := sync.WaitGroup{}
		assert.NoError(t, err)
		crObservable.DoOnNext(
			func(i interface{}) {
				wg.Done()
			})

		wg.Add(1)
		bottomChannel <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte("HelloWorld\nHelloWorld\nHelloWorld\nHelloWorld\n")))
		wg.Wait()
	})
}
