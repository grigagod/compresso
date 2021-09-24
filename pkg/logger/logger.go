package logger

import (
	"context"

	"github.com/jackc/pgx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const skipCallers = 1

type Logger interface {
	pgx.Logger
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

func (wl *WrappedLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	fields := make([]zapcore.Field, len(data))
	i := 0
	for k, v := range data {
		fields[i] = zap.Any(k, v)
		i++
	}

	logger := wl.SugaredLogger.Desugar()

	switch level {
	case pgx.LogLevelTrace:
		logger.Debug(msg, append(fields, zap.Stringer("PGX_LOG_LEVEL", level))...)
	case pgx.LogLevelDebug:
		logger.Debug(msg, fields...)
	case pgx.LogLevelInfo:
		logger.Info(msg, fields...)
	case pgx.LogLevelWarn:
		logger.Warn(msg, fields...)
	case pgx.LogLevelError:
		logger.Error(msg, fields...)
	default:
		logger.Error(msg, append(fields, zap.Stringer("PGX_LOG_LEVEL", level))...)
	}
}
