package gin

import (
	_ "github.com/fire9900/auth/docs"
	"github.com/fire9900/auth/internal/transport/gin/handlers"
	"github.com/fire9900/auth/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

func SetupRouter(userUseCase usecase.UseCase) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	userHandler := handlers.NewUserHandler(userUseCase)
	api := router.Group("/api/v1")
	{
		api.POST("/login", userHandler.Login)
		api.POST("/refresh", userHandler.RefreshToken)
		api.POST("/users", userHandler.Create)
		api.GET("/users", userHandler.GetAll)
		api.GET("/user/:id", userHandler.GetByID)
		auth := api.Group("/")
		auth.Use(handlers.AuthMiddleware())
		{
			auth.GET("/users/:email", userHandler.GetByEmail)
			auth.PUT("/users/:id", userHandler.UpdatePassword)
			auth.DELETE("/users/:id", userHandler.Delete)
			auth.POST("/user/:id", userHandler.CheckPassword)
			auth.GET("/logout", userHandler.Logout)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
