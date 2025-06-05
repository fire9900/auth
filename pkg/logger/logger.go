package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{"stdout", "./auth.log"}
	config.ErrorOutputPaths = []string{"stderr", "./auth-error.log"}

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(fmt.Sprintf("Ошибка запуска логгера: %w", err))
	}
	return nil
}
