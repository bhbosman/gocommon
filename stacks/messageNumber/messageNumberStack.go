package messageNumber

import (
	"context"
	"encoding/binary"
	"fmt"
	rxgo "github.com/ReactiveX/RxGo"
	"github.com/bhbosman/gocommon/constants"
	"github.com/bhbosman/gocommon/multiBlock"
	"github.com/bhbosman/gocommon/stacks/defs"
	"github.com/bhbosman/gocommon/stream"
)

func StackDefinition(
	userContext interface{},
	stackCancelFunc defs.CancelFunc,
	opts ...rxgo.Option) (*defs.StackDefinition, error) {
	if stackCancelFunc == nil {
		return nil, constants.InvalidParam
	}
	const stackName = "MessageNumber"

	return &defs.StackDefinition{
		Name: stackName,
		Inbound: func(index int, ctx context.Context) defs.BoundDefinition {
			return defs.BoundDefinition{
				PipeDefinition: func(params defs.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, constants.InvalidParam
					}
					errorState := false
					var number uint64 = 0
					return params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						index,
						params.ConnectionId,
						stackName,
						rxgo.StreamDirectionInbound,
						params.ConnectionManager,
						func(ctx context.Context, i stream.ReadWriterSize) (stream.ReadWriterSize, error) {
							if errorState {
								stackCancelFunc("In error state", true, constants.InvalidState)
								return nil, constants.InvalidState
							}
							buffer := [8]byte{0, 0, 0, 0, 0, 0, 0, 0}
							_, err := i.Read(buffer[:])
							if err != nil {
								stackCancelFunc("Could not read message number", true, err)
								errorState = true
								return nil, err
							}
							newNumber := binary.LittleEndian.Uint64(buffer[:])
							number++
							if newNumber != number {
								stackCancelFunc(
									fmt.Sprintf("Invalid number. Expected: %v, Received: %v", number, newNumber),
									true,
									err)
								errorState = true
								return nil, constants.InvalidSequenceNumber
							}
							return i, nil
						},
						opts...), nil
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
					var number uint64 = 0
					return params.Obs.(rxgo.InOutBoundObservable).MapInOutBound(
						index,
						params.ConnectionId,
						stackName,
						rxgo.StreamDirectionOutbound,
						params.ConnectionManager,
						func(ctx context.Context, i stream.ReadWriterSize) (stream.ReadWriterSize, error) {
							if errorState {
								stackCancelFunc("In error state", true, constants.InvalidState)
								return nil, constants.InvalidState
							}
							number++
							buffer := [8]byte{}
							binary.LittleEndian.PutUint64(buffer[:], number)
							rw := multiBlock.NewReaderWriterBlock(buffer[:])
							_ = rw.SetNext(i)
							return rw, nil
						},
						opts...), nil
				},
			}
		},
	}, nil
}
