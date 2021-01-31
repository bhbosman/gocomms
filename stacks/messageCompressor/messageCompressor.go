package messageCompressor

import (
	"compress/flate"
	"context"
	"encoding/binary"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"sync"

	"io"
	"net"
)

func StackDefinition(
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	if stackCancelFunc == nil {
		return nil, goerrors.InvalidParam
	}
	const stackName = "Compression"
	return &internal.StackDefinition{
		Name: stackName,
		Inbound: func(inOutBoundParams internal.InOutBoundParams) internal.BoundDefinition {
			decompressorStream := gomessageblock.NewReaderWriter()
			decompressor := flate.NewReader(decompressorStream)
			// decompressorMutex is here to safe guard panics when trying to destroy,
			//and while still busy processing data
			decompressorMutex := sync.Mutex{}
			return internal.BoundDefinition{
				PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, goerrors.InvalidParam
					}
					return pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						inOutBoundParams.Index,
						pipeParams.ConnectionId,
						stackName,
						rxgo.StreamDirectionInbound,
						pipeParams.ConnectionManager,
						func(ctx context.Context, incomingBlock goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
							decompressorMutex.Lock()
							defer decompressorMutex.Unlock()
							b := [8]byte{}
							_, err := incomingBlock.Read(b[:])
							if err != nil {
								stackCancelFunc("trying to read uncompressed length", true, err)
								return nil, err
							}
							uncompressedLength := int64(binary.LittleEndian.Uint64(b[:]))
							_, err = io.Copy(decompressorStream, incomingBlock)
							if err != nil {
								stackCancelFunc("trying to copy incoming data to pipeWriter", true, err)
								return nil, err
							}

							_, err = io.CopyN(incomingBlock, decompressor, uncompressedLength)
							if err != nil {
								stackCancelFunc("trying to copy uncompressed data to rws", true, err)
								return nil, err
							}

							return incomingBlock, nil
						}, opts...), nil
				},
				PipeState: internal.PipeState{
					Start: func(ctx context.Context) error {
						decompressorMutex.Lock()
						defer decompressorMutex.Unlock()
						return ctx.Err()
					},
					End: func() error {
						decompressorMutex.Lock()
						defer decompressorMutex.Unlock()
						return decompressor.Close()
					},
				},
			}
		},
		Outbound: func(inOutBoundParams internal.InOutBoundParams) internal.BoundDefinition {
			compressionStream := gomessageblock.NewReaderWriter()
			compression, err := flate.NewWriter(compressionStream, flate.DefaultCompression)
			// compressorMutex is here to safe guard panics when trying to destroy,
			//and while still busy processing data
			compressionMutex := sync.Mutex{}
			return internal.BoundDefinition{
				PipeDefinition: func(pipeParams internal.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, goerrors.InvalidParam
					}
					if err != nil {
						return nil, err
					}
					return pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						inOutBoundParams.Index,
						pipeParams.ConnectionId,
						stackName,
						rxgo.StreamDirectionOutbound,
						pipeParams.ConnectionManager,
						func(ctx context.Context, size goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
							compressionMutex.Lock()
							defer compressionMutex.Unlock()
							if ctx.Err() != nil {
								return nil, err
							}
							uncompressedSize, err := io.Copy(compression, size)
							if err != nil {
								return nil, err
							}

							if ctx.Err() != nil {
								return nil, err
							}
							err = compression.Flush()
							if err != nil {
								return nil, err
							}

							if ctx.Err() != nil {
								return nil, err
							}
							b := [8]byte{}
							binary.LittleEndian.PutUint64(b[:], uint64(uncompressedSize))

							if ctx.Err() != nil {
								return nil, err
							}
							_, err = size.Write(b[:])
							if err != nil {
								return nil, err
							}

							if ctx.Err() != nil {
								return nil, err
							}
							_, err = io.Copy(size, compressionStream)
							if err != nil {
								return nil, err
							}

							return size, nil
						},
						opts...), nil
				},
				PipeState: internal.PipeState{
					Start: func(ctx context.Context) error {
						compressionMutex.Lock()
						defer compressionMutex.Unlock()
						return ctx.Err()
					},
					End: func() error {
						compressionMutex.Lock()
						defer compressionMutex.Unlock()
						return compression.Close()
					},
				},
			}
		},
		StackState: internal.StackState{
			Start: func(startParams internal.StackStartStateParams) (net.Conn, error) {
				return startParams.Conn, startParams.Ctx.Err()
			},
			End: func(endParams internal.StackEndStateParams) error {
				return nil
			},
		},
	}, nil
}
