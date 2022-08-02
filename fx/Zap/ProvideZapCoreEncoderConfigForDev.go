package Zap

import (
	"go.uber.org/fx"
	"go.uber.org/zap/zapcore"
)

func ProvideZapCoreEncoderConfigForDev() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: "DEV",
			Target: func() *zapcore.EncoderConfig {
				return &zapcore.EncoderConfig{
					TimeKey:        "T",
					LevelKey:       "L",
					NameKey:        "N",
					CallerKey:      "C",
					FunctionKey:    "F",
					MessageKey:     "M",
					StacktraceKey:  "S",
					LineEnding:     zapcore.DefaultLineEnding,
					EncodeLevel:    zapcore.CapitalLevelEncoder,
					EncodeTime:     zapcore.ISO8601TimeEncoder,
					EncodeDuration: zapcore.StringDurationEncoder,
					EncodeCaller:   zapcore.ShortCallerEncoder,
				}
			},
		})
}
