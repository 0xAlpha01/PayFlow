package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ngwebwork/payflow/internal/services"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationHandler struct{ service *services.NotificationService }

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}
func notificationUser(c *gin.Context) (primitive.ObjectID, bool) {
	value, ok := c.Get("userID")
	if !ok {
		return primitive.NilObjectID, false
	}
	id, err := primitive.ObjectIDFromHex(value.(string))
	return id, err == nil
}
func (h *NotificationHandler) List(c *gin.Context) {
	userID, ok := notificationUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	result, err := h.service.List(c, userID, page, limit)
	if err != nil {
		utils.Failure(c, 500, "unable to load notifications", nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Notifications loaded", "data": result.Items, "pagination": gin.H{"page": result.Page, "limit": result.Limit, "total": result.Total, "totalPages": result.TotalPages}})
}
func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID, ok := notificationUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	count, err := h.service.UnreadCount(c, userID)
	if err != nil {
		utils.Failure(c, 500, "unable to load unread count", nil)
		return
	}
	utils.Success(c, 200, "Unread count loaded", gin.H{"count": count})
}
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID, ok := notificationUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.Failure(c, 400, "invalid notification id", nil)
		return
	}
	if err := h.service.MarkRead(c, userID, id); err != nil {
		utils.Failure(c, 404, "notification not found", nil)
		return
	}
	utils.Success(c, 200, "Notification marked as read", nil)
}
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID, ok := notificationUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	if err := h.service.MarkAllRead(c, userID); err != nil {
		utils.Failure(c, 500, "unable to update notifications", nil)
		return
	}
	utils.Success(c, 200, "All notifications marked as read", nil)
}
