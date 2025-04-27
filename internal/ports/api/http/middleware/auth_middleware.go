package middleware

import (
	"net/http"
	"strings"

	"salle-songbook-api/pkg/token"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		claims, err := token.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Guardamos claims en contexto para que lo use el handler
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func AdminOrComposerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "role not found"})
			return
		}

		r, _ := role.(string)
		if r != "admin" && r != "composer" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized role"})
			return
		}

		c.Next()
	}
}
