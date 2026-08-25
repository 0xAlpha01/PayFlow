package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ngwebwork/payflow/internal/handlers"
	"github.com/ngwebwork/payflow/internal/middleware"
	"github.com/ngwebwork/payflow/internal/services"
)

func Setup(r *gin.Engine, h *handlers.AuthHandler, jwt *services.JWTService, clientURL string) {
	r.Use(middleware.CORS(clientURL))
	api := r.Group("/api/v1/auth")
	api.POST("/signup", h.Signup)
	api.POST("/login", h.Login)
	api.POST("/refresh", h.Refresh)
	api.POST("/forgot-password", h.ForgotPassword)
	api.POST("/reset-password", h.ResetPassword)
	api.POST("/verify-email", h.VerifyEmail)
	api.POST("/resend-email-verification", h.ResendEmail)
	api.POST("/send-phone-code", h.SendPhone)
	api.POST("/verify-phone", h.VerifyPhone)
	protected := api.Group("")
	protected.Use(middleware.Auth(jwt))
	protected.POST("/logout", h.Logout)
	protected.GET("/me", h.Me)
}
