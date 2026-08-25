package utils

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, code int, message string, data any) {
	c.JSON(code, gin.H{"success": true, "message": message, "data": data})
}
func Failure(c *gin.Context, code int, message string, details any) {
	body := gin.H{"success": false, "message": message}
	if details != nil {
		body["error"] = details
	}
	c.JSON(code, body)
}
