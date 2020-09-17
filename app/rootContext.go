package app

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocommon/log"
	"go.uber.org/fx"
	log2 "log"
	"runtime"
	"time"
)

type RunTimeManager struct {
	logger        *log.SubSystemLogger
	ticker        *time.Ticker
	CancelContext context.Context
}

func NewRunTimeManager(logger *log.Factory, cancelContext context.Context) *RunTimeManager {
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

func RegisterRootContext() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{
				Target: func(params struct {
					fx.In
					Logger        *log.Factory
					Lifecycle     fx.Lifecycle
					CancelContext context.Context `name:"Application"`
				}) (*RunTimeManager, error) {
					result := NewRunTimeManager(params.Logger, params.CancelContext)
					params.Lifecycle.Append(fx.Hook{
						OnStart: func(ctx context.Context) error {
							return result.Start(ctx)
						},
						OnStop: func(ctx context.Context) error {
							return result.Stop(ctx)
						},
					})
					return result, nil
				}}),
		fx.Provide(
			fx.Annotated{
				Name: "Application",
				Target: func(params struct {
					fx.In
					Lifecycle fx.Lifecycle
				}) (context.Context, context.CancelFunc) {
					cancel, cancelFunc := context.WithCancel(context.Background())
					params.Lifecycle.Append(fx.Hook{
						OnStart: nil,
						OnStop: func(ctx context.Context) error {
							cancelFunc()
							return nil
						},
					})
					return cancel, cancelFunc
				}}),
		ProvidePubSub("Application"))
}
