package rxHelper

import (
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/messages"
	"github.com/reactivex/rxgo/v2"
)

type pushBlankMessage struct {
}

func HandleMessage(
	nextFunc rxgo.NextFunc,
	tryNextFunc goCommsDefinitions.TryNextFunc,
	isActive func() bool,
	primaryFeedCount func() int,
	listenChannels ...func() int,
) func(i interface{}) {
	messageCount := 0
	deltaCount := 0
	emptyErrorCount := 0

	// making sure we can deal with incoming nil pointers
	localIsActive := func(isActive func() bool) func() bool {
		if isActive != nil {
			return isActive
		}
		return func() bool {
			return true
		}
	}(isActive)

	// making sure we can deal with incoming nil pointers
	localPrimaryFeedCount := func(primaryFeedCount func() int) func() int {
		if primaryFeedCount != nil {
			return primaryFeedCount
		}
		return func() int {
			return 0
		}
	}(primaryFeedCount)

	return func(i interface{}) {
		if !localIsActive() {
			return
		}
		messageCount++
		deltaCount++
		switch i.(type) {
		case *messages.EmptyQueue:
			// swallow this message
			break
		case *pushBlankMessage:
			// swallow this message
			break
		default:
			if localIsActive() {
				nextFunc(i)
			} else {
				return
			}
			break
		}
		queueCount := localPrimaryFeedCount()
		if queueCount == 0 {
			if listenChannels != nil && len(listenChannels) > 0 {
				for _, ch := range listenChannels {
					if ch != nil {
						queueCount += ch()
						if queueCount > 0 {
							return
						}

					}
				}
			}
			if queueCount == 0 {
				emptyQueue := &messages.EmptyQueue{
					Count:           deltaCount,
					OverallCount:    messageCount,
					EmptyErrorCount: emptyErrorCount,
					ErrorHappen:     false,
				}
				if localIsActive() {
					nextFunc(emptyQueue)
					if emptyQueue.ErrorHappen {
						emptyErrorCount++
						if tryNextFunc != nil {
							tryNextFunc(&pushBlankMessage{})
						}
					}
				}
			}
		}
	}
}
