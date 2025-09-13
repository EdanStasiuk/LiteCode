package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/EdanStasiuk/LiteCode/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Check Redis for blacklisted token
		val, err := redis.Rdb.Get(c, tokenString).Result()
		if err == nil && val == "blacklisted" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
			return
		}

		// Parse JWT
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("userID", fmt.Sprintf("%v", claims["user_id"]))
		c.Set("token", tokenString) // for logout
		c.Next()
	}
}
