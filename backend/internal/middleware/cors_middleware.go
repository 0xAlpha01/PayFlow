package middleware

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(clientURL string) gin.HandlerFunc {
	origins := []string{clientURL}
	if strings.Contains(clientURL, "localhost") {
		origins = append(origins, strings.Replace(clientURL, "localhost", "127.0.0.1", 1))
	}
	if strings.Contains(clientURL, "127.0.0.1") {
		origins = append(origins, strings.Replace(clientURL, "127.0.0.1", "localhost", 1))
	}
	return cors.New(cors.Config{AllowOrigins: origins, AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, AllowHeaders: []string{"Origin", "Content-Type", "Authorization"}, AllowCredentials: true})
}
