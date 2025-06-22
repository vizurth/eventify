package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func AuthMiddleWareDefault(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: no token provided"})
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is invalid"})
			c.Abort()
			return
		}

		// Получаем user_id безопасно
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id is missing or invalid"})
			c.Abort()
			return
		}
		userID := int(userIDFloat)
		c.Set("userID", userID)

		// Получаем username
		username, ok := claims["username"].(string)
		if !ok || username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "username is missing or invalid"})
			c.Abort()
			return
		}
		c.Set("username", username)

		c.Next()
	}
}

func AuthMiddleWareOnAdmin(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// смотрим авторизованны ли мы
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		// получаем token через claims
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is invalid"})
		}
		// просматриваем role user чтобы узнать давать ему доступ к действию или нет
		role := claims["role"].(string)
		if role != "admin" {
			c.JSON(http.StatusLocked, gin.H{"error": "permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
