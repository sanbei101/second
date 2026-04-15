package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/phuslu/log"
)

var JWTSecret = []byte("campus-secondhand-secret-key")

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn().Msg("missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}
		if authHeader[:7] != "Bearer " {
			log.Warn().Str("token", authHeader).Msg("invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}
		tokenStr := authHeader[7:]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			return JWTSecret, nil
		})

		if err != nil || !token.Valid {
			log.Warn().Str("token", tokenStr).Msg("invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) uint {
	id, exist := c.Get("userID")
	if !exist {
		log.Error().Msg("userID not found in context")
		return 0
	}
	return id.(uint)
}
