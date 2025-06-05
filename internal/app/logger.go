package app

import (
	"github.com/fire9900/auth/pkg/logger"
	"go.uber.org/zap"
	"time"
)

func LoggerRun() {
	if err := logger.InitLogger(); err != nil {
		panic(err)
	}
	logger.Logger.Info("Логгер запущен", zap.String("date", time.Now().String()))
	defer logger.Logger.Sync()
}
