package middleware

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// AllowedOrigins expands the configured CLIENT_URL into the full set of origins the API
// accepts, shared by the CORS middleware and the websocket upgrader's origin check.
// CLIENT_URL may hold a single origin or a comma-separated list (e.g. the production
// frontend plus a local dev URL), so both can be allowed at once without swapping env
// vars between environments. Each entry is also expanded to its localhost/127.0.0.1
// equivalent for convenience during local development.
func AllowedOrigins(clientURL string) []string {
	var origins []string
	for _, raw := range strings.Split(clientURL, ",") {
		origin := strings.TrimSpace(strings.TrimSuffix(raw, "/"))
		if origin == "" {
			continue
		}
		origins = append(origins, origin)
		if strings.Contains(origin, "localhost") {
			origins = append(origins, strings.Replace(origin, "localhost", "127.0.0.1", 1))
		}
		if strings.Contains(origin, "127.0.0.1") {
			origins = append(origins, strings.Replace(origin, "127.0.0.1", "localhost", 1))
		}
	}
	return origins
}
func CORS(clientURL string) gin.HandlerFunc {
	return cors.New(cors.Config{AllowOrigins: AllowedOrigins(clientURL), AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, AllowHeaders: []string{"Origin", "Content-Type", "Authorization"}, AllowCredentials: true})
}
