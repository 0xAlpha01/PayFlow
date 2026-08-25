package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(clientURL string) gin.HandlerFunc {
	return cors.New(cors.Config{AllowOrigins: []string{clientURL}, AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, AllowHeaders: []string{"Origin", "Content-Type", "Authorization"}, AllowCredentials: true})
}
