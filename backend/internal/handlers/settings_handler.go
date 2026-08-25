// package handlers

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/ngwebwork/payflow/internal/services"
// 	"github.com/ngwebwork/payflow/internal/utils"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type SettingsHandler struct { service *services.SettingsService }
// func NewSettingsHandler(service *services.SettingsService) *SettingsHandler { return &SettingsHandler{service: service} }
// func settingsUser(c *gin.Context) (primitive.ObjectID, bool) { value, ok := c.Get("userID"); if !ok { return primitive.NilObjectID, false }; id, err := primitive.ObjectIDFromHex(value.(string)); return id, err == nil }
// func (h *SettingsHandler) Profile(c *gin.Context) { h.read(c, h.service.Profile) }
// func (h *SettingsHandler) Business(c *gin.Context) { h.read(c, h.service.Business) }
// func (h *SettingsHandler) Notifications(c *gin.Context) { h.read(c, h.service.Notifications) }
// func (h *SettingsHandler) read(c *gin.Context, fn func(*gin.Context, primitive.ObjectID) (bson.M, error)) { id, ok := settingsUser(c); if !ok { utils.Failure(c, 401, "unauthorized", nil); return }; data, err := fn(c, id); if err != nil { utils.Failure(c, 404, "settings not found", nil); return }; utils.Success(c, 200, "Settings loaded", data) }
// func (h *SettingsHandler) update(c *gin.Context, fn func(*gin.Context, primitive.ObjectID, bson.M) error) { id, ok := settingsUser(c); if !ok { utils.Failure(c, 401, "unauthorized", nil); return }; var payload bson.M; if c.ShouldBindJSON(&payload) != nil { utils.Failure(c, 400, "invalid request body", nil); return }; if err := fn(c, id, payload); err != nil { utils.Failure(c, 400, "unable to update settings", nil); return }; utils.Success(c, 200, "Settings updated", nil) }
// func (h *SettingsHandler) UpdateProfile(c *gin.Context) { h.update(c, h.service.UpdateProfile) }
// func (h *SettingsHandler) UpdateBusiness(c *gin.Context) { h.update(c, h.service.UpdateBusiness) }
// func (h *SettingsHandler) UpdateNotifications(c *gin.Context) { h.update(c, h.service.UpdateNotifications) }
// func (h *SettingsHandler) ChangePassword(c *gin.Context) { id, ok := settingsUser(c); if !ok { utils.Failure(c, 401, "unauthorized", nil); return }; var body struct { CurrentPassword string `json:"currentPassword"`; NewPassword string `json:"newPassword"` }; if c.ShouldBindJSON(&body) != nil { utils.Failure(c, 400, "invalid request body", nil); return }; if err := h.service.ChangePassword(c, id, body.CurrentPassword, body.NewPassword); err != nil { utils.Failure(c, 400, err.Error(), nil); return }; utils.Success(c, 200, "Password changed", nil) }

package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ngwebwork/payflow/internal/services"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SettingsHandler struct {
	service *services.SettingsService
}

func NewSettingsHandler(service *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		service: service,
	}
}

func settingsUser(c *gin.Context) (primitive.ObjectID, bool) {
	value, ok := c.Get("userID")
	if !ok {
		return primitive.NilObjectID, false
	}

	userID, ok := value.(string)
	if !ok {
		return primitive.NilObjectID, false
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return primitive.NilObjectID, false
	}

	return id, true
}

func (h *SettingsHandler) Profile(c *gin.Context) {
	id, ok := settingsUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}

	data, err := h.service.Profile(c.Request.Context(), id)
	if err != nil {
		utils.Failure(c, 404, "settings not found", nil)
		return
	}

	utils.Success(c, 200, "Settings loaded", data)
}

func (h *SettingsHandler) Business(c *gin.Context) {
	id, ok := settingsUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}

	data, err := h.service.Business(c.Request.Context(), id)
	if err != nil {
		utils.Failure(c, 404, "settings not found", nil)
		return
	}

	utils.Success(c, 200, "Settings loaded", data)
}

func (h *SettingsHandler) Notifications(c *gin.Context) {
	id, ok := settingsUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}

	data, err := h.service.Notifications(c.Request.Context(), id)
	if err != nil {
		utils.Failure(c, 404, "settings not found", nil)
		return
	}

	utils.Success(c, 200, "Settings loaded", data)
}

func (h *SettingsHandler) UpdateProfile(c *gin.Context) {
	h.updateProfile(c, h.service.UpdateProfile)
}

func (h *SettingsHandler) UpdateBusiness(c *gin.Context) {
	h.updateProfile(c, h.service.UpdateBusiness)
}

func (h *SettingsHandler) UpdateNotifications(c *gin.Context) {
	h.updateProfile(c, h.service.UpdateNotifications)
}

func (h *SettingsHandler) updateProfile(
	c *gin.Context,
	fn func(ctx context.Context, id primitive.ObjectID, update bson.M) error,
) {
	id, ok := settingsUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}

	var payload bson.M

	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.Failure(c, 400, "invalid request body", nil)
		return
	}

	if err := fn(c.Request.Context(), id, payload); err != nil {
		utils.Failure(c, 400, "unable to update settings", nil)
		return
	}

	utils.Success(c, 200, "Settings updated", nil)
}

func (h *SettingsHandler) ChangePassword(c *gin.Context) {
	id, ok := settingsUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}

	var body struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Failure(c, 400, "invalid request body", nil)
		return
	}

	if err := h.service.ChangePassword(
		c.Request.Context(),
		id,
		body.CurrentPassword,
		body.NewPassword,
	); err != nil {
		utils.Failure(c, 400, err.Error(), nil)
		return
	}

	utils.Success(c, 200, "Password changed", nil)
}