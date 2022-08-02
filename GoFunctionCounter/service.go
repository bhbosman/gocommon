package GoFunctionCounter

import (
	"fmt"
	"github.com/bhbosman/goCommsDefinitions"
	"github.com/bhbosman/gocommon/ChannelHandler"
	"github.com/bhbosman/gocommon/Services/IFxService"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/services/ISendMessage"
	"github.com/bhbosman/goerrors"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"time"
)

type service struct {
	ctx                    context.Context
	cancelFunc             context.CancelFunc
	cmdChannel             chan interface{}
	onData                 func() (IData, error)
	Logger                 *zap.Logger
	state                  IFxService.State
	cb                     func(names []string)
	UniqueReferenceService interfaces.IUniqueReferenceService
	UniqueSessionNumber    interfaces.IUniqueSessionNumber
}

func (self *service) GoRun(s string, cb func()) error {
	if cb == nil {
		return goerrors.InvalidParam
	}
	go func() {
		functionName := self.CreateFunctionName(s)
		_ = self.Add(functionName)
		defer func(name string) {
			_ = self.Remove(name)
		}(functionName)
		cb()
	}()
	return nil
}

func (self *service) CreateFunctionName(s string) string {
	return fmt.Sprintf("%v.%012d.%v",
		time.Now().Format("20060102150405"),
		self.UniqueSessionNumber.Next(),
		s,
	)
}

func (self *service) SetConnectionListChange(cb func(names []string)) {
	self.cb = cb
}

func (self *service) Remove(name string) error {
	add, err := CallIGoFunctionCounterRemove(self.ctx, self.cmdChannel, true, name)
	if err != nil {
		return err
	}
	return add.Args0
}

func (self *service) Add(name string) error {
	add, err := CallIGoFunctionCounterAdd(self.ctx, self.cmdChannel, true, name)
	if err != nil {
		return err
	}
	return add.Args0
}

func (self *service) Send(message interface{}) error {
	send, err := ISendMessage.CallISendMessageSend(
		self.ctx,
		self.cmdChannel,
		true,
		message,
	)
	if err != nil {
		return err
	}
	return send.Args0
}

func (self *service) OnStart(ctx context.Context) error {
	err := self.start(ctx)
	if err != nil {
		return err
	}
	self.state = IFxService.Started
	return nil
}

func (self *service) start(_ context.Context) error {
	instanceData, err := self.onData()
	if err != nil {
		return err
	}

	// this function is part of the GoFunctionCounter count
	go self.goStart(instanceData)
	return nil
}

func (self *service) goStart(instanceData IData) {
	defer func(cmdChannel <-chan interface{}) {
		//flush
		for range cmdChannel {
		}
	}(self.cmdChannel)

	ticker := time.NewTicker(time.Second * 1)
	defer func() {
		ticker.Stop()
	}()

	instanceData.SetConnectionListChange(self.cb)
	instanceData.Add(self.CreateFunctionName("GoFunctionCounter.service.gostart"))
	channelHandlerCallback := ChannelHandler.CreateChannelHandlerCallback(
		self.ctx,
		instanceData,
		[]ChannelHandler.ChannelHandler{
			{
				//BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(IGoFunctionCounter); ok {
						return ChannelEventsForIGoFunctionCounter(unk, message)
					}
					return false, nil
				},
			},
			{
				//BreakOnSuccess: false,
				Cb: func(next interface{}, message interface{}) (bool, error) {
					if unk, ok := next.(ISendMessage.ISendMessage); ok {
						return true, unk.Send(message)
					}
					return false, nil
				},
			},
		},
		func() int {
			return len(self.cmdChannel)
		},
		goCommsDefinitions.CreateTryNextFunc(self.cmdChannel),
	)
loop:
	for {
		select {
		case <-self.ctx.Done():
			err := instanceData.ShutDown()
			if err != nil {
				self.Logger.Error(
					"error on done",
					zap.Error(err))
			}
			break loop
		case <-ticker.C:
			if unk, ok := instanceData.(ISendMessage.ISendMessage); ok {
				unk.Send(&timerReceived{})
			}
			break
		case event, ok := <-self.cmdChannel:
			if !ok {
				return
			}
			breakLoop, err := channelHandlerCallback(event)
			if err != nil || breakLoop {
				break loop
			}
		}
	}
}

func (self *service) OnStop(ctx context.Context) error {
	err := self.shutdown(ctx)
	close(self.cmdChannel)
	self.state = IFxService.Stopped
	return err
}

func (self *service) shutdown(_ context.Context) error {
	self.cancelFunc()
	return nil
}

func (self *service) State() IFxService.State {
	return self.state
}

func (self *service) ServiceName() string {
	return "Go Function Registration"
}

func newService(
	parentContext context.Context,
	onData func() (IData, error),
	logger *zap.Logger,
	UniqueReferenceService interfaces.IUniqueReferenceService,
	UniqueSessionNumber interfaces.IUniqueSessionNumber,
) (IService, error) {
	localCtx, localCancelFunc := context.WithCancel(parentContext)
	result := &service{
		ctx:                    localCtx,
		cancelFunc:             localCancelFunc,
		cmdChannel:             make(chan interface{}, 32),
		onData:                 onData,
		Logger:                 logger,
		UniqueReferenceService: UniqueReferenceService,
		UniqueSessionNumber:    UniqueSessionNumber,
	}

	return result, nil
}
