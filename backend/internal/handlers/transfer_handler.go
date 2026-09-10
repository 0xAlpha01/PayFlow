package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ngwebwork/payflow/internal/services"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransferHandler struct{ service *services.TransferService }

func NewTransferHandler(service *services.TransferService) *TransferHandler {
	return &TransferHandler{service: service}
}
func transferUser(c *gin.Context) (primitive.ObjectID, bool) {
	value, ok := c.Get("userID")
	if !ok {
		return primitive.NilObjectID, false
	}
	id, err := primitive.ObjectIDFromHex(value.(string))
	return id, err == nil
}
func (h *TransferHandler) Balance(c *gin.Context) {
	userID, ok := transferUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	balance, err := h.service.Balance(c, userID)
	if err != nil {
		utils.Failure(c, 500, "unable to load balance", nil)
		return
	}
	utils.Success(c, 200, "Balance loaded", gin.H{"balance": balance, "currency": "NGN"})
}
func (h *TransferHandler) Send(c *gin.Context) {
	userID, ok := transferUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	var input services.SendTransferInput
	if c.ShouldBindJSON(&input) != nil {
		utils.Failure(c, 400, "invalid request body", nil)
		return
	}
	transfer, err := h.service.Send(c, userID, input)
	if err != nil {
		utils.Failure(c, 400, err.Error(), nil)
		return
	}
	utils.Success(c, 201, "Transfer completed successfully.", transfer)
}
func (h *TransferHandler) List(c *gin.Context) {
	userID, ok := transferUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	result, err := h.service.List(c, userID, page, limit)
	if err != nil {
		utils.Failure(c, 500, "unable to load transfers", nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Transfers loaded", "data": result.Items, "pagination": gin.H{"page": result.Page, "limit": result.Limit, "total": result.Total, "totalPages": result.TotalPages}})
}
