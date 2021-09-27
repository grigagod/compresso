package logger

import (
	"github.com/jackc/pgx"
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
}

type LogLevel = pgx.LogLevel

type WrappedLogger struct {
	*zap.SugaredLogger
}

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
