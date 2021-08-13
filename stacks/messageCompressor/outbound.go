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

func outbound(id uuid.UUID, stackCancelFunc internal.CancelFunc, opts ...rxgo.Option) func(inOutBoundParams internal.InOutBoundParams) (internal.IStackBoundDefinition, error) {
	return func(inOutBoundParams internal.InOutBoundParams) (internal.IStackBoundDefinition, error) {
		compressionStream := gomessageblock.NewReaderWriter()
		compression, err := flate.NewWriter(compressionStream, flate.DefaultCompression)
		// compressorMutex is here to safe guard panics when trying to destroy,
		//and while still busy processing data
		compressionMutex := sync.Mutex{}
		return &internal.StackBoundDefinition{
				PipeDefinition: func(stackData, pipeData interface{}, pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
					//if _, ok := stackData.(*StackData); !ok {
					//	return uuid.Nil, nil, WrongStackDataError(connectionType, stackData)
					//}
					return id,
						pipeParams.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
							inOutBoundParams.Index,
							pipeParams.ConnectionId,
							StackName,
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
				PipeState: &internal.PipeState{
					Destroy: func(stackData, pipeData interface{}) error {
						if closer, ok := stackData.(io.Closer); ok {
							return closer.Close()
						}
						return nil
					},
					Start: func(stackData, pipeData interface{}, ctx context.Context) error {
						compressionMutex.Lock()
						defer compressionMutex.Unlock()
						return ctx.Err()
					},
					ID: id,
					Create: func(stackData interface{}, ctx context.Context) (interface{}, error) {
						return internal.NewNoCloser(), nil
					},
					End: func(stackData, pipeData interface{}) error {
						compressionMutex.Lock()
						defer compressionMutex.Unlock()
						var err error
						err = multierr.Append(err, compression.Close())
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
