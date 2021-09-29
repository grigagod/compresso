package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const skipCallers = 1

type Logger interface {
	Debugw(msg string, keyValues ...interface{})
	Infow(msg string, keyValues ...interface{})
	Warnw(msg string, keyValues ...interface{})
	Errorw(msg string, keyValues ...interface{})
	Panicw(msg string, keyValues ...interface{})
	Fatalw(msg string, keyValues ...interface{})

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

type WrappedLogger struct {
	*zap.SugaredLogger
}

// GetLoggerConfig returns example of configuration
func GetLoggerConfig() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return cfg
}

func NewWrappedLogger(cfg zap.Config, args ...interface{}) (*WrappedLogger, error) {
	l, err := cfg.Build(zap.AddCallerSkip(skipCallers))
	if err != nil {
		return nil, err
	}

	sl := l.Sugar().With(args...)
	return &WrappedLogger{sl}, nil
}
