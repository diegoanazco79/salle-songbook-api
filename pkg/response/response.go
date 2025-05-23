package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
		"error":   nil,
	})
}

func Created(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": message,
		"data":    data,
		"error":   nil,
	})
}

func Error(c *gin.Context, status int, message string, errDetail string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
		"data":    nil,
		"error":   errDetail,
	})
}
