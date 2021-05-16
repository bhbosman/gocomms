package messageCompressor

import (
	"compress/flate"
	"context"
	"encoding/binary"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
	"io"
	"sync"
)

func inbound(id uuid.UUID, stackCancelFunc internal.CancelFunc, opts ...rxgo.Option) func(inOutBoundParams internal.InOutBoundParams) (internal.IStackBoundDefinition, error) {
	return func(inOutBoundParams internal.InOutBoundParams) (internal.IStackBoundDefinition, error) {
		decompressorStream := gomessageblock.NewReaderWriter()
		decompressor := flate.NewReader(decompressorStream)
		// decompressorMutex is here to safe guard panics when trying to destroy,
		//and while still busy processing data
		decompressorMutex := sync.Mutex{}
		return &internal.StackBoundDefinition{
				PipeDefinition: func(stackData, pipeData interface{}, pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
					return id,
						pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
							inOutBoundParams.Index,
							pipeParams.ConnectionId,
							StackName,
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
				PipeState: &internal.PipeState{
					ID: id,
					Create: func(pipeData interface{}, ctx context.Context) (interface{}, error) {
						return internal.NewNoCloser(), nil
					},
					Destroy: func(stackData, pipeData interface{}) error {
						if closer, ok := pipeData.(io.Closer); ok {
							return closer.Close()
						}
						return nil
					},
					Start: func(stackData, pipeData interface{}, ctx context.Context) error {
						decompressorMutex.Lock()
						defer decompressorMutex.Unlock()
						return ctx.Err()
					},
					End: func(stackData, pipeData interface{}) error {
						decompressorMutex.Lock()
						defer decompressorMutex.Unlock()
						var err error = nil
						err = multierr.Append(err, decompressor.Close())
						if closer, ok := pipeData.(io.Closer); ok {
							err = multierr.Append(err, closer.Close())
						}
						return err
					},
				},
			},
			nil
	}
}
