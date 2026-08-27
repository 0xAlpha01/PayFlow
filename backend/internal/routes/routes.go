package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ngwebwork/payflow/internal/handlers"
	"github.com/ngwebwork/payflow/internal/middleware"
	"github.com/ngwebwork/payflow/internal/services"
)

func Setup(r *gin.Engine, h *handlers.AuthHandler, dashboard *handlers.DashboardHandler, resources *handlers.ResourceHandler, settings *handlers.SettingsHandler, jwt *services.JWTService, clientURL string) {
	r.Use(middleware.CORS(clientURL))
	api := r.Group("/api/v1/auth")
	api.POST("/signup", h.Signup)
	api.POST("/login", h.Login)
	api.POST("/refresh", h.Refresh)
	api.POST("/forgot-password", h.ForgotPassword)
	api.POST("/reset-password", h.ResetPassword)
	api.POST("/verify-email", h.VerifyEmail)
	api.POST("/resend-email-verification", h.ResendEmail)
	protected := api.Group("")
	protected.Use(middleware.Auth(jwt))
	protected.POST("/logout", h.Logout)
	protected.GET("/me", h.Me)
	protected.POST("/send-phone-code", h.SendPhone)
	protected.POST("/resend-phone-code", h.ResendPhone)
	protected.POST("/verify-phone", h.VerifyPhone)
	dashboardAPI := r.Group("/api/v1/dashboard")
	dashboardAPI.Use(middleware.Auth(jwt))
	dashboardAPI.GET("/summary", dashboard.Summary)
	dashboardAPI.GET("/revenue", dashboard.Revenue)
	dashboardAPI.GET("/transactions", dashboard.Transactions)
	dashboardAPI.GET("/customers", dashboard.Customers)
	analyticsAPI := r.Group("/api/v1/analytics")
	analyticsAPI.Use(middleware.Auth(jwt))
	analyticsAPI.GET("/revenue", dashboard.AnalyticsRevenue)
	analyticsAPI.GET("/transactions", dashboard.AnalyticsTransactions)
	analyticsAPI.GET("/customers", dashboard.AnalyticsCustomers)
	analyticsAPI.GET("/payments", dashboard.AnalyticsPayments)
	for _, collection := range []string{"customers", "payments", "transactions", "invoices", "subscriptions"} {
		resourceAPI := r.Group("/api/v1/" + collection)
		resourceAPI.Use(middleware.Auth(jwt))
		resourceAPI.GET("", func(c *gin.Context) {
			c.Params = append(c.Params, gin.Param{Key: "resource", Value: collection})
			resources.List(c)
		})
		resourceAPI.POST("", func(c *gin.Context) {
			c.Params = append(c.Params, gin.Param{Key: "resource", Value: collection})
			resources.Create(c)
		})
		resourceAPI.GET("/:id", func(c *gin.Context) {
			c.Params = append(c.Params, gin.Param{Key: "resource", Value: collection})
			resources.Get(c)
		})
		resourceAPI.PUT("/:id", func(c *gin.Context) {
			c.Params = append(c.Params, gin.Param{Key: "resource", Value: collection})
			resources.Update(c)
		})
		resourceAPI.DELETE("/:id", func(c *gin.Context) {
			c.Params = append(c.Params, gin.Param{Key: "resource", Value: collection})
			resources.Delete(c)
		})
	}
	for _, action := range []struct{ collection, name string }{{"invoices", "send"}, {"invoices", "mark-paid"}, {"subscriptions", "cancel"}} {
		resourceAPI := r.Group("/api/v1/" + action.collection)
		resourceAPI.Use(middleware.Auth(jwt))
		resourceAPI.POST("/:id/"+action.name, func(c *gin.Context) {
			c.Params = append(c.Params, gin.Param{Key: "resource", Value: action.collection}, gin.Param{Key: "action", Value: action.name})
			resources.Action(c)
		})
	}
	settingsAPI := r.Group("/api/v1/settings")
	settingsAPI.Use(middleware.Auth(jwt))
	settingsAPI.GET("/profile", settings.Profile)
	settingsAPI.PUT("/profile", settings.UpdateProfile)
	settingsAPI.GET("/business", settings.Business)
	settingsAPI.PUT("/business", settings.UpdateBusiness)
	settingsAPI.GET("/notifications", settings.Notifications)
	settingsAPI.PUT("/notifications", settings.UpdateNotifications)
	settingsAPI.POST("/change-password", settings.ChangePassword)
}
