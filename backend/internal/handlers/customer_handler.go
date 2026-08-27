package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ngwebwork/payflow/internal/repositories"
	"github.com/ngwebwork/payflow/internal/services"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerHandler struct{ service *services.CustomerService }

func NewCustomerHandler(service *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}
func customerUserID(c *gin.Context) (primitive.ObjectID, bool) {
	value, ok := c.Get("userID")
	if !ok {
		return primitive.NilObjectID, false
	}
	id, err := primitive.ObjectIDFromHex(value.(string))
	return id, err == nil
}
func customerID(c *gin.Context) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	return id, err == nil
}
func (h *CustomerHandler) List(c *gin.Context) {
	userID, ok := customerUserID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	query := repositories.CustomerQuery{Page: page, Limit: limit, Search: c.Query("search"), Status: c.Query("status"), SortBy: c.Query("sortBy"), SortOrder: c.Query("sortOrder")}
	if value := c.Query("createdFrom"); value != "" {
		query.CreatedFrom, _ = time.Parse("2006-01-02", value)
	}
	if value := c.Query("createdTo"); value != "" {
		query.CreatedTo, _ = time.Parse("2006-01-02", value)
	}
	result, err := h.service.List(c, userID, query)
	if err != nil {
		utils.Failure(c, 400, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Customers retrieved successfully.", "data": result.Items, "pagination": gin.H{"page": result.Page, "limit": result.Limit, "total": result.Total, "totalPages": result.TotalPages}})
}
func (h *CustomerHandler) Create(c *gin.Context) {
	userID, ok := customerUserID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	var input services.CustomerInput
	if c.ShouldBindJSON(&input) != nil {
		utils.Failure(c, 400, "invalid request body", nil)
		return
	}
	customer, err := h.service.Create(c, userID, input)
	if err != nil {
		status := 400
		if err.Error() == "a customer with this email already exists" || err.Error() == "a customer with this phone already exists" {
			status = 409
		}
		utils.Failure(c, status, err.Error(), nil)
		return
	}
	utils.Success(c, 201, "Customer created successfully.", customer)
}
func (h *CustomerHandler) Get(c *gin.Context) {
	userID, ok := customerUserID(c)
	id, valid := customerID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	if !valid {
		utils.Failure(c, 400, "invalid customer id", nil)
		return
	}
	customer, err := h.service.Get(c, userID, id)
	if err != nil {
		utils.Failure(c, 500, "unable to retrieve customer", nil)
		return
	}
	if customer == nil {
		utils.Failure(c, 404, "customer not found", nil)
		return
	}
	utils.Success(c, 200, "Customer retrieved successfully.", customer)
}
func (h *CustomerHandler) Update(c *gin.Context) {
	userID, ok := customerUserID(c)
	id, valid := customerID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	if !valid {
		utils.Failure(c, 400, "invalid customer id", nil)
		return
	}
	var input services.CustomerInput
	if c.ShouldBindJSON(&input) != nil {
		utils.Failure(c, 400, "invalid request body", nil)
		return
	}
	customer, err := h.service.Update(c, userID, id, input)
	if err == mongo.ErrNoDocuments || customer == nil {
		utils.Failure(c, 404, "customer not found", nil)
		return
	}
	if err != nil {
		status := 400
		if err.Error() == "a customer with this email already exists" || err.Error() == "a customer with this phone already exists" {
			status = 409
		}
		utils.Failure(c, status, err.Error(), nil)
		return
	}
	utils.Success(c, 200, "Customer updated successfully.", customer)
}
func (h *CustomerHandler) Delete(c *gin.Context) {
	userID, ok := customerUserID(c)
	id, valid := customerID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	if !valid {
		utils.Failure(c, 400, "invalid customer id", nil)
		return
	}
	err := h.service.Delete(c, userID, id)
	if err == mongo.ErrNoDocuments {
		utils.Failure(c, 404, "customer not found", nil)
		return
	}
	if err != nil {
		utils.Failure(c, 409, err.Error(), nil)
		return
	}
	utils.Success(c, 200, "Customer deleted successfully.", nil)
}
func (h *CustomerHandler) Payments(c *gin.Context)     { h.history(c, true) }
func (h *CustomerHandler) Transactions(c *gin.Context) { h.history(c, false) }
func (h *CustomerHandler) history(c *gin.Context, payments bool) {
	userID, ok := customerUserID(c)
	id, valid := customerID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	if !valid {
		utils.Failure(c, 400, "invalid customer id", nil)
		return
	}
	var data any
	var err error
	if payments {
		data, err = h.service.Payments(c, userID, id)
	} else {
		data, err = h.service.Transactions(c, userID, id)
	}
	if err == mongo.ErrNoDocuments {
		utils.Failure(c, 404, "customer not found", nil)
		return
	}
	if err != nil {
		utils.Failure(c, 500, "unable to retrieve customer history", nil)
		return
	}
	utils.Success(c, 200, "Customer history retrieved successfully.", data)
}
