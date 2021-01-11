package gocommon

import (
	"context"
	"fmt"
	"github.com/bhbosman/gologging"
	log2 "log"
	"runtime"
	"time"
)

type RunTimeManager struct {
	logger        *gologging.SubSystemLogger
	ticker        *time.Ticker
	CancelContext context.Context
}

func NewRunTimeManager(logger *gologging.Factory, cancelContext context.Context) *RunTimeManager {
	return &RunTimeManager{
		logger:        logger.Create("RunTimeManager"),
		CancelContext: cancelContext,
	}
}

func (self *RunTimeManager) Stop(ctx context.Context) error {
	self.ticker.Stop()
	return nil
}

func (self *RunTimeManager) Start(ctx context.Context) error {
	self.ticker = time.NewTicker(time.Second * 5)
	go func() {
		for {
			select {
			case _, ok := <-self.ticker.C:
				if !ok {
					return
				}
				self.logger.LogWithLevel(0, func(logger *log2.Logger) {
					logger.Printf(fmt.Sprintf("NumGoroutine: %v", runtime.NumGoroutine()))
				})
			}
		}
	}()
	return nil
}
