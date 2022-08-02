package Zap

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ProvideZapConfigForDev(outputPaths []string, errorOutputPaths []string) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: "DEV",
			Target: func(params struct {
				fx.In
				EncoderConfig *zapcore.EncoderConfig `name:"DEV"`
			}) *zap.Config {
				if outputPaths == nil {
					outputPaths = []string{"stderr"}
				}
				if errorOutputPaths == nil {
					errorOutputPaths = []string{"stderr"}
				}

				return &zap.Config{
					Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
					Development:       false,
					DisableCaller:     true,
					DisableStacktrace: true,
					Sampling:          nil,
					Encoding:          "json",
					EncoderConfig:     *params.EncoderConfig,
					OutputPaths:       outputPaths,
					ErrorOutputPaths:  errorOutputPaths,
					InitialFields:     nil,
				}
			},
		})
}
