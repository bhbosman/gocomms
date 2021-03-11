package bvisMessageBreaker_test

import (
	"context"
	"crypto/sha1"
	"encoding/binary"
	internal2 "github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/stacks/bvisMessageBreaker"
	"github.com/bhbosman/gocomms/stacks/bvisMessageBreaker/internal"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/reactivex/rxgo/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"sync"
	"testing"
)

func TestBuildBlocksInbound(t *testing.T) {
	buildMessage := func(s string, hash io.Writer) []byte {
		stringAsBytes := []byte(s)
		if hash != nil {
			_, _ = hash.Write(stringAsBytes)
		}
		result := make([]byte, len(stringAsBytes)+8)

		copy(result[0:4], []byte{'B', 'V', 'I', 'S'})
		binary.LittleEndian.PutUint32(result[4:8], uint32(len(stringAsBytes)))
		copy(result[8:len(stringAsBytes)+8], stringAsBytes)
		return result
	}
	t.Run("Inbound: Check error function ", func(t *testing.T) {
		_, err := bvisMessageBreaker.StackDefinition(internal2.ServerConnection, "", nil, nil, nil)
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid param")
	})
	t.Run("Inbound, using multiBlock.NewReaderWriterBlock", func(t *testing.T) {
		t.Run("Inbound: Invalid signature", func(t *testing.T) {
			wg := sync.WaitGroup{}
			var blockError error
			blocks, err := bvisMessageBreaker.StackDefinition(
				internal2.ServerConnection,
				"",
				func(s string, inbound bool, err error) {
					blockError = err
					wg.Done()
				},
				nil,
				nil)
			assert.NoError(t, err)

			ch := make(chan rxgo.Item)
			defer close(ch)
			obs := rxgo.FromChannel(ch)
			inbound, err := blocks.Inbound(internal2.NewInOutBoundParams(0, context.TODO())).PipeDefinition(
				internal2.PipeDefinitionParams{
					ConnectionId:      "",
					ConnectionManager: nil,
					CancelContext:     context.Background(),
					StackCancelFunc:   func(context string, inbound bool, err error) {},
					Obs:               obs,
				})
			assert.NoError(t, err)
			inbound.DoOnNext(func(i interface{}) {
			})
			wg.Add(1)
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}))
			wg.Wait()
			assert.Error(t, blockError)
			assert.Equal(t, goerrors.InvalidSignature, blockError)
		})
		t.Run("Inbound: Valid Signature", func(t *testing.T) {
			wg := sync.WaitGroup{}
			var blockError error
			var buildState internal.BuildMessageState
			blocks, err := bvisMessageBreaker.StackDefinition(
				internal2.ServerConnection,
				"",
				func(s string, inbound bool, err error) {
					blockError = err
				},
				func(from, to internal.BuildMessageState, length uint32) {
					buildState = to
					wg.Done()
				},
				nil)
			assert.NoError(t, err)

			ch := make(chan rxgo.Item)
			defer close(ch)
			obs := rxgo.FromChannel(ch)

			inbound, err := blocks.Inbound(internal2.NewInOutBoundParams(0, context.TODO())).PipeDefinition(
				internal2.PipeDefinitionParams{
					CancelContext:   context.Background(),
					StackCancelFunc: func(context string, inbound bool, err error) {},
					Obs:             obs,
				})
			assert.NoError(t, err)

			inbound.DoOnNext(func(i interface{}) {
			})
			assert.Equal(t, internal.BuildMessageStateReadMessageSignature, buildState)
			wg.Add(1)
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{'B'}))
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{'V'}))
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{'I'}))
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{'S'}))
			wg.Wait()
			assert.NoError(t, blockError)
			assert.Equal(t, internal.BuildMessageStateReadMessageLength, buildState)
		})
		t.Run("Inbound: Valid Signature and length of 4", func(t *testing.T) {
			wg := sync.WaitGroup{}
			var blockError error
			var buildState internal.BuildMessageState
			var buildLength uint32
			blocks, err := bvisMessageBreaker.StackDefinition(
				internal2.ServerConnection,
				"",
				func(s string, inbound bool, err error) {
					blockError = err
				},
				func(from, to internal.BuildMessageState, length uint32) {
					buildState = to
					buildLength = length
					wg.Done()
				},
				nil)
			assert.NoError(t, err)

			ch := make(chan rxgo.Item)
			defer close(ch)
			obs := rxgo.FromChannel(ch)

			inbound, err := blocks.Inbound(internal2.NewInOutBoundParams(0, context.TODO())).PipeDefinition(
				internal2.PipeDefinitionParams{
					CancelContext:   context.Background(),
					StackCancelFunc: func(context string, inbound bool, err error) {},
					Obs:             obs,
				})
			assert.NoError(t, err)

			inbound.DoOnNext(func(i interface{}) {
			})
			assert.Equal(t, internal.BuildMessageStateReadMessageSignature, buildState)
			wg.Add(1)
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{'B'}))
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{'V'}))
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{'I'}))
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{'S'}))
			wg.Wait()
			assert.NoError(t, blockError)
			assert.Equal(t, internal.BuildMessageStateReadMessageLength, buildState)

			wg.Add(1)
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{4}))
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{0}))
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{0}))
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{0}))
			wg.Wait()
			assert.NoError(t, blockError)
			assert.Equal(t, internal.BuildMessageStateReadMessageData, buildState)
			assert.Equal(t, uint32(4), buildLength)
		})
		t.Run("Inbound: Valid Signature and length of 4 with data", func(t *testing.T) {
			wg := sync.WaitGroup{}
			var blockError error
			var buildState internal.BuildMessageState
			blocks, err := bvisMessageBreaker.StackDefinition(
				internal2.ServerConnection,
				"",
				func(s string, inbound bool, err error) {
					blockError = err
				},
				func(from, to internal.BuildMessageState, length uint32) {
					buildState = to
					wg.Done()
				},
				nil)
			assert.NoError(t, err)
			ch := make(chan rxgo.Item)
			defer close(ch)
			obs := rxgo.FromChannel(ch)
			inbound, err := blocks.Inbound(internal2.NewInOutBoundParams(0, context.TODO())).PipeDefinition(
				internal2.PipeDefinitionParams{
					CancelContext:   context.Background(),
					StackCancelFunc: func(context string, inbound bool, err error) {},
					Obs:             obs,
				})
			assert.NoError(t, err)
			inbound.DoOnNext(func(i interface{}) {
				wg.Done()
			})
			assert.Equal(t, internal.BuildMessageStateReadMessageSignature, buildState)
			wg.Add(4)
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{'B', 'V', 'I', 'S'}))
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{4, 0, 0, 0}))
			ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{4, 3, 2, 1}))
			wg.Wait()
			assert.NoError(t, blockError)
			assert.Equal(t, internal.BuildMessageStateReadMessageSignature, buildState)
		})
		t.Run("Inbound: Blocks of data", func(t *testing.T) {
			blocks, err := bvisMessageBreaker.StackDefinition(
				internal2.ServerConnection,
				"",
				func(s string, inbound bool, err error) {
				},
				func(from, to internal.BuildMessageState, length uint32) {
				},
				nil)
			assert.NoError(t, err)
			ch := make(chan rxgo.Item)
			defer close(ch)
			obs := rxgo.FromChannel(ch)
			inbound, err := blocks.Inbound(internal2.NewInOutBoundParams(0, context.TODO())).PipeDefinition(
				internal2.PipeDefinitionParams{
					CancelContext:   context.Background(),
					StackCancelFunc: func(context string, inbound bool, err error) {},
					Obs:             obs,
				})
			assert.NoError(t, err)
			wg := sync.WaitGroup{}
			receivedHash := sha1.New()
			inbound.DoOnNext(
				func(i interface{}) {
					switch v := i.(type) {
					case *gomessageblock.ReaderWriter:
						_, _ = io.Copy(receivedHash, v)
						wg.Done()
					default:
						assert.Fail(t, "wrong type")
					}
				})

			mutex := sync.Mutex{}
			t.Run("one message", func(t *testing.T) {
				mutex.Lock()
				defer mutex.Unlock()
				wg.Add(1)
				ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock(buildMessage("Hello world", nil)))
				wg.Wait()
			})
			t.Run("two message as one block", func(t *testing.T) {
				mutex.Lock()
				defer mutex.Unlock()

				data01 := buildMessage("Hello world0001", nil)
				data02 := buildMessage("Hello world0001", nil)

				block01 := gomessageblock.NewReaderWriterBlock(data01)
				block02 := gomessageblock.NewReaderWriterBlock(data02)
				_ = block01.SetNext(block02)
				wg.Add(2)

				ch <- rxgo.Of(block01)
				wg.Wait()
				wg.Add(1)
				ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock(buildMessage("Hello world", nil)))
				wg.Wait()

			})

			t.Run("500 message", func(t *testing.T) {
				mutex.Lock()
				defer mutex.Unlock()
				receivedHash.Reset()
				sendHash := sha1.New()
				n := 500
				wg.Add(n)
				for i := 0; i < n; i++ {
					ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock(buildMessage("Hello world0001", sendHash)))
				}
				wg.Wait()
				a := receivedHash.Sum(nil)
				b := sendHash.Sum(nil)
				assert.Equal(t, a, b)
			})
			t.Run("500 message broken up", func(t *testing.T) {
				mutex.Lock()
				defer mutex.Unlock()
				receivedHash.Reset()
				sendHash := sha1.New()
				n := 500
				wg.Add(n)
				for i := 0; i < n; i++ {
					data := buildMessage("Hello world0001", sendHash)
					for _, c := range data {
						ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{c}))
					}
				}
				wg.Wait()
				a := receivedHash.Sum(nil)
				b := sendHash.Sum(nil)
				assert.Equal(t, a, b)
			})
			t.Run("two message", func(t *testing.T) {
				mutex.Lock()
				defer mutex.Unlock()
				wg.Add(1)
				ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock(buildMessage("Hello world0001", nil)))
				wg.Wait()
				wg.Add(1)
				ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock(buildMessage("Hello world0002", nil)))
				wg.Wait()
			})

		})
	})
	t.Run("Outbound", func(t *testing.T) {
		blocks, err := bvisMessageBreaker.StackDefinition(
			internal2.ServerConnection,
			"",
			func(s string, inbound bool, err error) {
			},
			func(stateFrom, stateTo internal.BuildMessageState, length uint32) {

			},
			nil)
		assert.NoError(t, err)

		ch := make(chan rxgo.Item)
		defer close(ch)
		obs := rxgo.FromChannel(ch)
		outbound, err := blocks.Outbound(internal2.NewInOutBoundParams(0, context.TODO())).PipeDefinition(
			internal2.PipeDefinitionParams{
				CancelContext:   context.Background(),
				StackCancelFunc: func(context string, inbound bool, err error) {},
				Obs:             obs,
			})
		assert.NoError(t, err)
		wg := sync.WaitGroup{}
		outbound.DoOnNext(
			func(i interface{}) {
				switch v := i.(type) {
				case *gomessageblock.ReaderWriter:
					p := make([]byte, v.Size())
					_, err := v.Read(p)
					assert.NoError(t, err)
					mustBe := []byte{'B', 'V', 'I', 'S', 4, 0, 0, 0, 0, 0, 0, 0}
					assert.Equal(t, mustBe, p)
					wg.Done()
				}
			})
		wg.Add(1)
		ch <- rxgo.Of(gomessageblock.NewReaderWriterBlock([]byte{0, 0, 0, 0}))
		wg.Wait()
	})

}
