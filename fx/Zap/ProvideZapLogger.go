package Zap

import (
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ProvideZapLogger() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotated{Target: func(params struct {
				fx.In
				DevConfig  *zap.Config `name:"DEV"`
				ProdConfig *zap.Config `name:"PROD"`
				Which      string      `name:"LogConfiguration" optional:"true"`
			}) (*zap.Config, error) {
				var config *zap.Config
				switch params.Which {
				case "":
					config = params.DevConfig
					break
				case "DEV":
					config = params.DevConfig
				case "PROD":
					config = params.ProdConfig
				default:
					return nil, fmt.Errorf("invalid LogConfiguration. Value is: %v", params.Which)
				}
				return config, nil
			}}),
		fx.Provide(
			fx.Annotated{
				Target: func(params struct {
					fx.In
					Config         *zap.Config
					ZapLoggerCores []zapcore.Core        `group:"ZapCore.Core.Loggers"`
					ZapErrorCores  []zapcore.WriteSyncer `group:"ZapCore.Core.Errors"`
				}) (*zap.Logger, error) {
					if len(params.ZapLoggerCores) == 0 {
						return params.Config.Build()
					}
					zapCores := zapcore.NewTee(params.ZapLoggerCores...)
					var options []zap.Option
					if len(params.ZapErrorCores) > 0 {
						options = append(
							options,
							zap.ErrorOutput(zapcore.NewMultiWriteSyncer(params.ZapErrorCores...)))
					}
					logger := zap.New(zapCores, options...)
					return logger, nil
				},
			},
		),
	)
}
