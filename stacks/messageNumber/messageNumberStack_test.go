package messageNumber

import (
	"context"
	"encoding/binary"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gomessageblock"
	"github.com/reactivex/rxgo/v2"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestMessageNumberStackDefinition(t *testing.T) {
	t.Run("Inbound: Check error function ", func(t *testing.T) {
		_, err := StackDefinition(nil, nil)
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid param")
	})

	t.Run("Inbound001", func(t *testing.T) {
		stack, err := StackDefinition(
			nil,
			func(s string, inbound bool, err error) {

			})
		assert.NoError(t, err)
		ch := make(chan rxgo.Item)
		obs := rxgo.FromChannel(ch)
		obs2, err := stack.Inbound(1, context.Background()).PipeDefinition(
			internal.PipeDefinitionParams{
				CancelContext:   context.Background(),
				StackCancelFunc: func(context string, inbound bool, err error) {},
				Obs:             obs,
			})
		assert.NoError(t, err)
		wg := sync.WaitGroup{}

		doOnNext := obs2.DoOnNext(
			func(i interface{}) {
				switch v := i.(type) {
				case *gomessageblock.ReaderWriter:
					p := make([]byte, v.Size())
					_, _ = v.Read(p)
					assert.Equal(t, p, []byte{1, 2, 3, 4, 5, 6, 7, 8})
				}

				wg.Done()
			})
		go func() {
			<-doOnNext
			return
		}()
		wg.Add(1)
		item := rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8}))
		item.SendContext(context.Background(), ch)
		wg.Wait()
	})

	t.Run("Inbound1000", func(t *testing.T) {
		stack, err := StackDefinition(
			nil,
			func(s string, inbound bool, err error) {
				assert.NoError(t, err)
			},
			rxgo.WithBufferedChannel(1024))
		assert.NoError(t, err)
		ch := make(chan rxgo.Item)
		obs := rxgo.FromChannel(ch)
		obs2, err := stack.Inbound(2, context.Background()).PipeDefinition(
			internal.PipeDefinitionParams{

				CancelContext:   context.Background(),
				StackCancelFunc: func(context string, inbound bool, err error) {},
				Obs:             obs,
			})
		assert.NoError(t, err)
		wg := sync.WaitGroup{}

		_ = obs2.DoOnNext(
			func(i interface{}) {
				switch v := i.(type) {
				case *gomessageblock.ReaderWriter:
					p := make([]byte, v.Size())
					_, _ = v.Read(p)
					assert.Equal(t, p, []byte{1, 2, 3, 4, 5, 6, 7, 8})
				}
				wg.Done()
			}, rxgo.WithBufferedChannel(1024))
		var n uint64 = 1000
		wg.Add(int(n))

		var i uint64
		dd := [8]byte{1, 2, 3, 4, 5, 6, 7, 8}
		for i = 1; i <= n; i++ {
			data := make([]byte, 16)
			copy(data[8:], dd[:])

			binary.LittleEndian.PutUint64(data[0:8], i)

			item := rxgo.Of(gomessageblock.NewReaderWriterBlock(data))
			item.SendContext(context.Background(), ch)
		}
		wg.Wait()
	})

	t.Run("Outbound", func(t *testing.T) {
		stack, err := StackDefinition(
			nil,
			func(s string, inbound bool, err error) {

			})
		assert.NoError(t, err)
		ch := make(chan rxgo.Item)
		obs := rxgo.FromChannel(ch)
		obs2, err := stack.Outbound(9, context.Background()).PipeDefinition(
			internal.PipeDefinitionParams{
				CancelContext:   context.Background(),
				StackCancelFunc: func(context string, inbound bool, err error) {},
				Obs:             obs,
			})
		assert.NoError(t, err)
		wg := sync.WaitGroup{}
		dataIn := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		doOnNext := obs2.DoOnNext(
			func(i interface{}) {
				wg.Done()
			})
		go func() {
			<-doOnNext
			return
		}()
		wg.Add(1)
		ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock(dataIn[:]))
		wg.Wait()
	})

}
