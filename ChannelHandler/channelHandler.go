package ChannelHandler

import (
	"context"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/messages"
	"github.com/bhbosman/gocommon/services/ISendMessage"
)

type Callback func(message interface{}) (bool, error)
type ChannelHandler struct {
	Cb func(next interface{}, message interface{}) (bool, error)
}

type pushBlankMessage struct {
}

func CreateChannelHandlerCallback(
	ctx context.Context,
	next interface{},
	handlers []ChannelHandler,
	queueCount func() int,
	pushMessage goCommsDefinitions.TryNextFunc,
) Callback {
	handleEmptyQueue := func() func() {
		if queueCount == nil {
			return func() {
			}
		}
		deltaCount := 0
		messageCount := 0
		return func() {
			deltaCount++
			messageCount++
			if queueCount() == 0 {
				if sm, ok := next.(ISendMessage.ISendMessage); ok {
					emptyQueue := &messages.EmptyQueue{
						Count:        deltaCount,
						OverallCount: messageCount,
					}
					_ = sm.Send(emptyQueue)
					if emptyQueue.ErrorHappen {
						if pushMessage != nil {
							pushMessage(&pushBlankMessage{})
						}
					}
				}
				deltaCount = 0
			}
		}
	}()

	return func(message interface{}) (bool, error) {
		if ctx.Err() != nil {
			return true, ctx.Err()
		}

		if _, ok := message.(*pushBlankMessage); !ok {
			for _, handler := range handlers {
				if ctx.Err() != nil {
					return true, ctx.Err()
				}
				if handler.Cb != nil {
					success, err := handler.Cb(next, message)
					if err != nil {
						return true, err
					}
					if !success {
						continue
					}
					break
				}
			}
		}
		handleEmptyQueue()

		return false, nil
	}
}
