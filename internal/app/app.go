package app

import (
	"github.com/fire9900/auth/internal/repository"
	"github.com/fire9900/auth/internal/transport/gin"
	"github.com/fire9900/auth/internal/usecase"
	"github.com/fire9900/auth/pkg/database"
	"github.com/fire9900/auth/pkg/logger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
)

func Run() {
	db, err := database.NewSQLiteConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)

	router := gin.SetupRouter(userUseCase)

	if err := router.Run(":8080"); err != nil {
		logger.Logger.Fatal("Ошибка запуска сервера на порту :8080",
			zap.Error(err),
			zap.String("app", "database"))
	}
	logger.Logger.Info("Микросервис стартует на порту :8080")
}
