package application

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(errors.Join(err, errInitLoggerFailed))
	}

	return logger.
		WithOptions(
			zap.AddStacktrace(zapcore.FatalLevel),
			zap.WithCaller(false),
		)
}

var errInitLoggerFailed = errors.New("failed to init logger")

func loggerInterceptor(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)
		for i := 0; i < len(fields); i += 2 {
			f = append(f, zap.Any(fields[i].(string), fields[i+1]))
		}
		l = l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg)
		case logging.LevelInfo:
			l.Info(msg)
		case logging.LevelWarn:
			l.Warn(msg)
		case logging.LevelError:
			l.Error(msg)
		}
	})
}
