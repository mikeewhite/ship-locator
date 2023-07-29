package clog

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.SugaredLogger
}

// newZapLogger creates a new logger
func newZapLogger() *zapLogger {
	cfg := zap.NewProductionConfig()
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeTime = utcFormat()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return &zapLogger{logger: logger.Sugar()}
}

func (z *zapLogger) Infof(msg string, args ...interface{}) {
	z.logger.Infof(msg, args...)
}

func (z *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	z.logger.Infow(msg, keysAndValues...)
}

func (z *zapLogger) Warnf(msg string, args ...interface{}) {
	z.logger.Warnf(msg, args...)
}

func (z *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	z.logger.Warnw(msg, keysAndValues...)
}

func (z *zapLogger) Errorf(msg string, args ...interface{}) {
	z.logger.Errorf(msg, args...)
}

func (z *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	z.logger.Errorw(msg, keysAndValues...)
}

func (z *zapLogger) Flush() {
	_ = z.logger.Sync()
}

// utcFormat defines a formatter that will output the date in UTC RFC339 format
func utcFormat() zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(time.RFC3339))
	}
}
