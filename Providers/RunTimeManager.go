package Providers

import (
	"context"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"go.uber.org/zap"
	"runtime"
	"time"
)

type RunTimeManager struct {
	logger            *zap.Logger
	cancelCtx         context.Context
	cancelFunc        context.CancelFunc
	goFunctionCounter GoFunctionCounter.IService
}

func NewRunTimeManager(
	logger *zap.Logger,
	cancelContext context.Context,
	goFunctionCounter GoFunctionCounter.IService,
) *RunTimeManager {
	ctx, f := context.WithCancel(cancelContext)
	return &RunTimeManager{
		logger:            logger.Named("RunTimeManager"),
		cancelCtx:         ctx,
		cancelFunc:        f,
		goFunctionCounter: goFunctionCounter,
	}
}

func (self *RunTimeManager) Stop(ctx context.Context) error {
	self.cancelFunc()
	return nil
}

func (self *RunTimeManager) Start(ctx context.Context) error {
	return self.goFunctionCounter.GoRun(
		"RunTimeManager.Start",
		func() {
			ticker := time.NewTicker(time.Second * 5)
			defer func() {
				ticker.Stop()
			}()
		loop:
			for {
				select {
				case <-self.cancelCtx.Done():
					break loop
				case _, ok := <-ticker.C:
					if !ok {
						break loop
					}
					self.logger.Info(
						"NumGoroutine",
						zap.Int("NumGoroutine", runtime.NumGoroutine()))
				}
			}
		},
	)
}
