// Package logger provides logger interface.
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const skipCallers = 1

// Logger defines logger inteface.
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

// WrappedLogger imlements Logger interface unsing zap.SugaredLogger.
type WrappedLogger struct {
	*zap.SugaredLogger
}

// GetLoggerConfig return example of configuration.
func GetLoggerConfig() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return cfg
}

// NewWrappedLogger return new WrappedLogger with given config and additional args.
func NewWrappedLogger(cfg zap.Config, args ...interface{}) (*WrappedLogger, error) {
	l, err := cfg.Build(zap.AddCallerSkip(skipCallers))
	if err != nil {
		return nil, err
	}

	sl := l.Sugar().With(args...)
	return &WrappedLogger{sl}, nil
}
