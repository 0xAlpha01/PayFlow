package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ngwebwork/payflow/internal/services"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DashboardHandler struct{ service *services.DashboardService }

func NewDashboardHandler(service *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}
func (h *DashboardHandler) userID(c *gin.Context) (primitive.ObjectID, bool) {
	value, ok := c.Get("userID")
	if !ok {
		return primitive.NilObjectID, false
	}
	id, err := primitive.ObjectIDFromHex(value.(string))
	return id, err == nil
}
func (h *DashboardHandler) Summary(c *gin.Context) {
	id, ok := h.userID(c)
	if !ok {
		utils.Failure(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	summary, err := h.service.Summary(c.Request.Context(), id)
	if err != nil {
		utils.Failure(c, http.StatusInternalServerError, "unable to load dashboard summary", nil)
		return
	}
	utils.Success(c, http.StatusOK, "Dashboard summary", summary)
}
func (h *DashboardHandler) Revenue(c *gin.Context) {
	id, ok := h.userID(c)
	if !ok {
		utils.Failure(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	data, err := h.service.Revenue(c.Request.Context(), id, c.DefaultQuery("period", "30d"))
	if err != nil {
		utils.Failure(c, http.StatusInternalServerError, "unable to load revenue data", nil)
		return
	}
	utils.Success(c, http.StatusOK, "Revenue data", data)
}
func (h *DashboardHandler) Transactions(c *gin.Context) {
	id, ok := h.userID(c)
	if !ok {
		utils.Failure(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	data, err := h.service.RecentTransactions(c.Request.Context(), id)
	if err != nil {
		utils.Failure(c, http.StatusInternalServerError, "unable to load transactions", nil)
		return
	}
	utils.Success(c, http.StatusOK, "Recent transactions", data)
}
func (h *DashboardHandler) Customers(c *gin.Context) {
	id, ok := h.userID(c)
	if !ok {
		utils.Failure(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	data, err := h.service.RecentCustomers(c.Request.Context(), id)
	if err != nil {
		utils.Failure(c, http.StatusInternalServerError, "unable to load customers", nil)
		return
	}
	utils.Success(c, http.StatusOK, "Recent customers", data)
}
func (h *DashboardHandler) AnalyticsRevenue(c *gin.Context) { h.Revenue(c) }
func (h *DashboardHandler) AnalyticsTransactions(c *gin.Context) { h.Summary(c) }
func (h *DashboardHandler) AnalyticsCustomers(c *gin.Context) { h.Summary(c) }
func (h *DashboardHandler) AnalyticsPayments(c *gin.Context) { h.Summary(c) }
