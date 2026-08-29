package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ngwebwork/payflow/internal/repositories"
	"github.com/ngwebwork/payflow/internal/services"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

type InvoiceHandler struct{ service *services.InvoiceService }

func NewInvoiceHandler(s *services.InvoiceService) *InvoiceHandler { return &InvoiceHandler{s} }
func invoiceUser(c *gin.Context) (primitive.ObjectID, bool) {
	v, ok := c.Get("userID")
	if !ok {
		return primitive.NilObjectID, false
	}
	id, e := primitive.ObjectIDFromHex(v.(string))
	return id, e == nil
}
func invoiceID(c *gin.Context) (primitive.ObjectID, bool) {
	id, e := primitive.ObjectIDFromHex(c.Param("id"))
	return id, e == nil
}
func (h *InvoiceHandler) List(c *gin.Context) {
	u, ok := invoiceUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	r, e := h.service.List(c, u, repositories.InvoiceQuery{Page: p, Limit: l, Search: c.Query("search"), Status: c.Query("status")})
	if e != nil {
		utils.Failure(c, 400, e.Error(), nil)
		return
	}
	c.JSON(200, gin.H{"success": true, "data": r.Items, "pagination": gin.H{"page": r.Page, "limit": r.Limit, "total": r.Total, "totalPages": r.TotalPages}})
}
func (h *InvoiceHandler) Get(c *gin.Context) {
	u, ok := invoiceUser(c)
	id, valid := invoiceID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	if !valid {
		utils.Failure(c, 400, "invalid invoice id", nil)
		return
	}
	i, e := h.service.Get(c, u, id)
	if e != nil {
		utils.Failure(c, 500, "unable to retrieve invoice", nil)
		return
	}
	if i == nil {
		utils.Failure(c, 404, "invoice not found", nil)
		return
	}
	utils.Success(c, 200, "Invoice retrieved successfully.", i)
}
func (h *InvoiceHandler) Create(c *gin.Context) {
	u, ok := invoiceUser(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	var in services.InvoiceInput
	if c.ShouldBindJSON(&in) != nil {
		utils.Failure(c, 400, "invalid request body", nil)
		return
	}
	i, e := h.service.Create(c, u, in)
	if e != nil {
		code := 400
		if e == mongo.ErrNoDocuments {
			code = 404
		}
		utils.Failure(c, code, e.Error(), nil)
		return
	}
	utils.Success(c, 201, "Invoice created successfully.", i)
}
func (h *InvoiceHandler) Update(c *gin.Context) {
	u, ok := invoiceUser(c)
	id, valid := invoiceID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	if !valid {
		utils.Failure(c, 400, "invalid invoice id", nil)
		return
	}
	var in services.InvoiceInput
	if c.ShouldBindJSON(&in) != nil {
		utils.Failure(c, 400, "invalid request body", nil)
		return
	}
	i, e := h.service.Update(c, u, id, in)
	if e == mongo.ErrNoDocuments || i == nil {
		utils.Failure(c, 404, "invoice not found", nil)
		return
	}
	if e != nil {
		utils.Failure(c, 400, e.Error(), nil)
		return
	}
	utils.Success(c, 200, "Invoice updated successfully.", i)
}
func (h *InvoiceHandler) Delete(c *gin.Context) {
	u, ok := invoiceUser(c)
	id, valid := invoiceID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	if !valid {
		utils.Failure(c, 400, "invalid invoice id", nil)
		return
	}
	if e := h.service.Delete(c, u, id); e != nil {
		utils.Failure(c, 404, "invoice not found", nil)
		return
	}
	utils.Success(c, 200, "Invoice deleted successfully.", nil)
}
func (h *InvoiceHandler) Send(c *gin.Context)     { h.action(c, false) }
func (h *InvoiceHandler) MarkPaid(c *gin.Context) { h.action(c, true) }
func (h *InvoiceHandler) action(c *gin.Context, paid bool) {
	u, ok := invoiceUser(c)
	id, valid := invoiceID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	if !valid {
		utils.Failure(c, 400, "invalid invoice id", nil)
		return
	}
	var e error
	if paid {
		e = h.service.MarkPaid(c, u, id)
	} else {
		e = h.service.Send(c, u, id)
	}
	if e == mongo.ErrNoDocuments {
		utils.Failure(c, 404, "invoice not found", nil)
		return
	}
	if e != nil {
		utils.Failure(c, 502, "unable to send invoice", nil)
		return
	}
	utils.Success(c, http.StatusOK, "Invoice updated successfully.", nil)
}
