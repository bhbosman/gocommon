package messageCompressor

import (
	"compress/flate"
	"context"
	"encoding/binary"
	rxgo "github.com/ReactiveX/RxGo"
	"sync"

	"github.com/bhbosman/gocommon/constants"
	"github.com/bhbosman/gocommon/multiBlock"
	"github.com/bhbosman/gocommon/stacks/defs"
	"github.com/bhbosman/gocommon/stream"
	"io"
	"net"
	"net/url"
)

func StackDefinition(
	stackCancelFunc defs.CancelFunc,
	opts ...rxgo.Option) (*defs.StackDefinition, error) {
	if stackCancelFunc == nil {
		return nil, constants.InvalidParam
	}
	const stackName = "Compression"
	return &defs.StackDefinition{
		Name: stackName,
		Inbound: func(index int, ctx context.Context) defs.BoundDefinition {
			decompressorStream := multiBlock.NewReaderWriter()
			decompressor := flate.NewReader(decompressorStream)
			mutex := sync.Mutex{}
			return defs.BoundDefinition{
				PipeDefinition: func(params defs.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, constants.InvalidParam
					}
					return params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						index,
						params.ConnectionId,
						stackName,
						rxgo.StreamDirectionInbound,
						params.ConnectionManager,
						func(ctx context.Context, incomingBlock stream.ReadWriterSize) (stream.ReadWriterSize, error) {
							mutex.Lock()
							defer mutex.Unlock()
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
				PipeState: defs.PipeState{
					Start: func(ctx context.Context) error {
						mutex.Lock()
						defer mutex.Unlock()
						return ctx.Err()
					},
					End: func() error {
						mutex.Lock()
						defer mutex.Unlock()
						return decompressor.Close()
					},
				},
			}
		},
		Outbound: func(index int, ctx context.Context) defs.BoundDefinition {
			compressionStream := multiBlock.NewReaderWriter()
			compression, err := flate.NewWriter(compressionStream, flate.DefaultCompression)
			mutex := sync.Mutex{}
			return defs.BoundDefinition{
				PipeDefinition: func(params defs.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, constants.InvalidParam
					}
					if err != nil {
						return nil, err
					}
					return params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						index,
						params.ConnectionId,
						stackName,
						rxgo.StreamDirectionOutbound,
						params.ConnectionManager,
						func(ctx context.Context, size stream.ReadWriterSize) (stream.ReadWriterSize, error) {
							mutex.Lock()
							defer mutex.Unlock()
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
				PipeState: defs.PipeState{
					Start: func(ctx context.Context) error {
						mutex.Lock()
						defer mutex.Unlock()
						return ctx.Err()
					},
					End: func() error {
						mutex.Lock()
						defer mutex.Unlock()
						return compression.Close()
					},
				},
			}
		},
		StackState: defs.StackState{
			Start: func(conn net.Conn, url *url.URL, ctx context.Context, cancelFunc defs.CancelFunc) (net.Conn, error) {
				return conn, ctx.Err()
			},
			End: func() error {
				return nil
			},
		},
	}, nil
}
