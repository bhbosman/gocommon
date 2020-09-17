package messageBreaker

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	rxgo "github.com/ReactiveX/RxGo"
	"github.com/bhbosman/gocommon/constants"
	"github.com/bhbosman/gocommon/multiBlock"
	"github.com/bhbosman/gocommon/stacks/defs"
	"github.com/bhbosman/gocommon/stacks/messageBreaker/internal"
	"github.com/bhbosman/goprotoextra"
	"reflect"
)

func StackDefinition(
	stackCancelFunc defs.CancelFunc,
	stateFunc func(stateFrom, stateTo internal.BuildMessageState, length uint32),
	connectionManager rxgo.IPublishToConnectionManager,
	opts ...rxgo.Option) (*defs.StackDefinition, error) {
	if stackCancelFunc == nil {
		return nil, constants.InvalidParam
	}

	marker := [4]byte{'B', 'V', 'I', 'S'}
	markerAsUInt32 := binary.LittleEndian.Uint32(marker[:])
	const stackName = "MessageBreaker"
	return &defs.StackDefinition{
		Name: stackName,
		Inbound: func(index int, ctx context.Context) defs.BoundDefinition {
			var nextChannel chan rxgo.Item
			return defs.BoundDefinition{
				PipeDefinition: func(params defs.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, constants.InvalidParam
					}
					rw := multiBlock.NewReaderWriter()
					state := internal.BuildMessageStateReadMessageSignature
					var length uint32 = 0

					errorState := false
					inboundState := func(onNext func(data []byte)) {
						var p [4]byte
						canContinue := true
						for canContinue {
							switch state {
							case 0:
								if rw.Size() >= 4 {
									_, err := rw.Read(p[:])
									if err != nil {
										stackCancelFunc("Could not read signature", true, err)
										errorState = true
										return
									}
									c := bytes.Compare(p[:], marker[:])
									if c != 0 {
										stackCancelFunc("Signature incorrect", true, constants.InvalidSignature)
										errorState = true
										return
									}
									prev := state
									state = internal.BuildMessageStateReadMessageLength
									if stateFunc != nil {
										stateFunc(prev, state, length)
									}
									break
								} else {
									canContinue = false
								}
							case 1:
								if rw.Size() >= 4 {
									_, err := rw.Read(p[:])
									if err != nil {
										stackCancelFunc("Could not read length", true, err)
										errorState = true
										return
									}
									length = binary.LittleEndian.Uint32(p[:])
									prev := state
									state = internal.BuildMessageStateReadMessageData
									if stateFunc != nil {
										stateFunc(prev, state, length)
									}
									break
								} else {
									canContinue = false
								}
							case 2:
								if uint32(rw.Size()) >= length {
									dataBlock := make([]byte, length)
									_, err := rw.Read(dataBlock)
									if err != nil {
										stackCancelFunc("Could not read data block", true, err)
										errorState = true
										return
									}
									onNext(dataBlock)

									length = 0
									prev := state
									state = internal.BuildMessageStateReadMessageSignature
									if stateFunc != nil {
										stateFunc(prev, state, length)
									}
									break
								} else {
									canContinue = false
								}
							}
						}
					}
					nextChannel = make(chan rxgo.Item)
					_ = params.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
						index,
						params.ConnectionId,
						stackName,
						rxgo.StreamDirectionInbound,
						connectionManager,
						func(ctx context.Context, i goprotoextra.ReadWriterSize) {
							if errorState {
								stackCancelFunc("In error state", true, constants.InvalidState)
								return
							}
							switch v := i.(type) {
							case *multiBlock.ReaderWriter:
								err := rw.SetNext(v)
								if err != nil {
									params.StackCancelFunc("rw.SetNext()", true, err)
									return
								}
								inboundState(func(dataBlock []byte) {
									item := rxgo.Of(multiBlock.NewReaderWriterBlock(dataBlock))
									item.SendContext(ctx, nextChannel)
								})
							default:
								stackCancelFunc(
									fmt.Sprintf("Invalid type(%v) received", reflect.TypeOf(i).String()),
									true,
									constants.InvalidType)
								errorState = true
								return
							}
						},
						opts...)
					return rxgo.FromChannel(nextChannel, opts...), nil
				},
				PipeState: defs.PipeState{
					Start: func(ctx context.Context) error {
						return ctx.Err()
					},
					End: func() error {
						close(nextChannel)
						return nil
					},
				},
			}
		},
		Outbound: func(index int, ctx context.Context) defs.BoundDefinition {
			return defs.BoundDefinition{
				PipeDefinition: func(params defs.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, constants.InvalidParam
					}
					errorState := false
					return params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						index,
						params.ConnectionId,
						stackName,
						rxgo.StreamDirectionOutbound,
						params.ConnectionManager,
						func(ctx context.Context, i goprotoextra.ReadWriterSize) (goprotoextra.ReadWriterSize, error) {
							if errorState {
								stackCancelFunc("In error state", false, constants.InvalidState)
								return nil, constants.InvalidState
							}
							block := make([]byte, 8)
							binary.LittleEndian.PutUint32(block[0:4], markerAsUInt32)
							binary.LittleEndian.PutUint32(block[4:8], uint32(i.Size()))
							result := multiBlock.NewReaderWriterBlock(block)
							err := result.SetNext(i)
							if err != nil {
								params.StackCancelFunc("rw.SetNext()", false, err)
								return nil, err
							}
							return result, nil
						},
						opts...), nil
				},
			}
		},
	}, nil
}
