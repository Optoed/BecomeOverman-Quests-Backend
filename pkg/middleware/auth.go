package middleware

import (
	"BecomeOverMan/internal/services"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := services.ValidateJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// сохраняем user_id в контекст
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// ==== Helping Functions ====

func GetUserID(c *gin.Context) (int, error) {
	userIDKey, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("Can't get user_id from context")
	}

	userID, ok := userIDKey.(int)
	if !ok {
		return 0, errors.New("user ID is not integer")
	}

	return userID, nil
}
