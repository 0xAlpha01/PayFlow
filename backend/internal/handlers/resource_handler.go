package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ngwebwork/payflow/internal/repositories"
	"github.com/ngwebwork/payflow/internal/services"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResourceHandler struct{ service *services.ResourceService }

func NewResourceHandler(service *services.ResourceService) *ResourceHandler {
	return &ResourceHandler{service: service}
}
func (h *ResourceHandler) id(c *gin.Context) (primitive.ObjectID, bool) {
	value, ok := c.Get("userID")
	if !ok {
		return primitive.NilObjectID, false
	}
	id, err := primitive.ObjectIDFromHex(value.(string))
	return id, err == nil
}
func (h *ResourceHandler) collection(c *gin.Context) string {
	return strings.Trim(c.Param("resource"), " /")
}
func (h *ResourceHandler) List(c *gin.Context) {
	userID, ok := h.id(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	query := map[string]string{"page": c.Query("page"), "limit": c.Query("limit"), "search": c.Query("search"), "status": c.Query("status"), "customerId": c.Query("customerId")}
	result, err := h.service.List(c, h.collection(c), userID, query)
	if err != nil {
		utils.Failure(c, 500, "unable to load records", nil)
		return
	}
	utils.Success(c, 200, "Records loaded", result)
}
func (h *ResourceHandler) Get(c *gin.Context) {
	userID, ok := h.id(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	id, err := repositories.ParseID(c.Param("id"))
	if err != nil {
		utils.Failure(c, 400, "invalid record id", nil)
		return
	}
	item, err := h.service.Get(c, h.collection(c), userID, id)
	if err != nil {
		utils.Failure(c, 500, "unable to load record", nil)
		return
	}
	if item == nil {
		utils.Failure(c, 404, "record not found", nil)
		return
	}
	utils.Success(c, 200, "Record loaded", item)
}
func (h *ResourceHandler) Create(c *gin.Context) {
	userID, ok := h.id(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	var payload bson.M
	if c.ShouldBindJSON(&payload) != nil {
		utils.Failure(c, 400, "invalid request body", nil)
		return
	}
	delete(payload, "userId")
	item, err := h.service.Create(c, h.collection(c), userID, payload)
	if err != nil {
		utils.Failure(c, 400, err.Error(), nil)
		return
	}
	utils.Success(c, 201, "Record created", item)
}
func (h *ResourceHandler) Update(c *gin.Context) {
	userID, ok := h.id(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	id, err := repositories.ParseID(c.Param("id"))
	if err != nil {
		utils.Failure(c, 400, "invalid record id", nil)
		return
	}
	var payload bson.M
	if c.ShouldBindJSON(&payload) != nil {
		utils.Failure(c, 400, "invalid request body", nil)
		return
	}
	if err := h.service.Update(c, h.collection(c), userID, id, payload); err != nil {
		code := 500
		if strings.Contains(err.Error(), "no documents") {
			code = 404
		}
		utils.Failure(c, code, "unable to update record", nil)
		return
	}
	utils.Success(c, 200, "Record updated", nil)
}
func (h *ResourceHandler) Delete(c *gin.Context) {
	userID, ok := h.id(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	id, err := repositories.ParseID(c.Param("id"))
	if err != nil {
		utils.Failure(c, 400, "invalid record id", nil)
		return
	}
	if err := h.service.Delete(c, h.collection(c), userID, id); err != nil {
		utils.Failure(c, 404, "record not found", nil)
		return
	}
	utils.Success(c, 200, "Record deleted", nil)
}
func (h *ResourceHandler) Action(c *gin.Context) {
	userID, ok := h.id(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	id, err := repositories.ParseID(c.Param("id"))
	if err != nil {
		utils.Failure(c, 400, "invalid record id", nil)
		return
	}
	action := c.Param("action")
	status := map[string]string{"send": "pending", "mark-paid": "paid", "cancel": "cancelled"}[action]
	if status == "" {
		utils.Failure(c, 400, "unsupported action", nil)
		return
	}
	if err := h.service.Update(c, h.collection(c), userID, id, bson.M{"status": status}); err != nil {
		utils.Failure(c, 404, "record not found", nil)
		return
	}
	utils.Success(c, 200, "Record updated", nil)
}

var _ = http.StatusOK
