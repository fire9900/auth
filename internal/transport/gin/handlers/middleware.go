package handlers

import (
	"fmt"
	"github.com/fire9900/auth/pkg/auth"
	"github.com/fire9900/auth/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Logger.Error(fmt.Sprintf("Отсутствие хедера %s", authHeader))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Отсутствие хедера"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("Невалидный токен %s", tokenString),
				zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Невалидный токен"})
			return
		}
		logger.Logger.Info("Успешная проверка авторизации пользователя")
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
