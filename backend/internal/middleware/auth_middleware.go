package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ngwebwork/payflow/internal/services"
	"net/http"
	"strings"
)

func Auth(jwt *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "missing or invalid authorization header"})
			return
		}
		claims, err := jwt.Parse(parts[1], "access")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid or expired access token"})
			return
		}
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
