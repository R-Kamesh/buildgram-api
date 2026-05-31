package utils

import "github.com/gin-gonic/gin"

// SendSuccess formats standard success payloads
func SendSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status": "success",
		"data":   data,
	})
}

// SendError formats standardized human-readable error payloads
func SendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": message,
	})
}