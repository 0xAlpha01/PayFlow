package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ngwebwork/payflow/internal/services"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnalyticsHandler struct{ s *services.AnalyticsService }

func NewAnalyticsHandler(s *services.AnalyticsService) *AnalyticsHandler { return &AnalyticsHandler{s} }
func (a *AnalyticsHandler) user(c *gin.Context) (primitive.ObjectID, bool) {
	v, ok := c.Get("userID")
	if !ok {
		return primitive.NilObjectID, false
	}
	id, e := primitive.ObjectIDFromHex(v.(string))
	return id, e == nil
}
func (a *AnalyticsHandler) serve(c *gin.Context, fn func(primitive.ObjectID) (any, error)) {
	u, ok := a.user(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	x, e := fn(u)
	if e != nil {
		utils.Failure(c, 500, "unable to load analytics", nil)
		return
	}
	utils.Success(c, 200, "Analytics retrieved.", x)
}
func (a *AnalyticsHandler) Revenue(c *gin.Context) {
	a.serve(c, func(u primitive.ObjectID) (any, error) { return a.s.Revenue(c, u, c.DefaultQuery("period", "30d")) })
}
func (a *AnalyticsHandler) Transactions(c *gin.Context) {
	a.serve(c, func(u primitive.ObjectID) (any, error) {
		return a.s.Transactions(c, u, c.DefaultQuery("period", "30d"))
	})
}
func (a *AnalyticsHandler) Customers(c *gin.Context) {
	a.serve(c, func(u primitive.ObjectID) (any, error) { return a.s.Customers(c, u, c.DefaultQuery("period", "30d")) })
}
func (a *AnalyticsHandler) Payments(c *gin.Context) {
	a.serve(c, func(u primitive.ObjectID) (any, error) { return a.s.Payments(c, u, c.DefaultQuery("period", "30d")) })
}
func (a *AnalyticsHandler) Summary(c *gin.Context) {
	a.serve(c, func(u primitive.ObjectID) (any, error) { return a.s.Payments(c, u, "30d") })
}
