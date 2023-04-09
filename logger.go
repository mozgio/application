package Application

import (
	"errors"

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
