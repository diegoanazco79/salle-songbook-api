package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OnlyAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Only admins can perform this action",
				"error":   "forbidden",
			})
			return
		}
		c.Next()
	}
}
