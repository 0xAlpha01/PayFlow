package handlers

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ngwebwork/payflow/internal/repositories"
	"github.com/ngwebwork/payflow/internal/services"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentHandler struct {
	service        *services.PaymentService
	paystackSecret string
}

func NewPaymentHandler(service *services.PaymentService, paystackSecret string) *PaymentHandler {
	return &PaymentHandler{service: service, paystackSecret: paystackSecret}
}
func paymentUserID(c *gin.Context) (primitive.ObjectID, bool) {
	value, ok := c.Get("userID")
	if !ok {
		return primitive.NilObjectID, false
	}
	id, err := primitive.ObjectIDFromHex(value.(string))
	return id, err == nil
}
func (h *PaymentHandler) Initialize(c *gin.Context) {
	userID, ok := paymentUserID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	var input services.InitializePaymentInput
	if c.ShouldBindJSON(&input) != nil {
		utils.Failure(c, 400, "invalid request body", nil)
		return
	}
	_, result, err := h.service.Initialize(c, userID, input)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.Failure(c, 404, "customer not found", nil)
		} else {
			utils.Failure(c, 400, err.Error(), nil)
		}
		return
	}
	utils.Success(c, 200, "Payment initialized successfully.", gin.H{"authorizationUrl": result.AuthorizationURL, "reference": result.Reference})
}
func (h *PaymentHandler) Verify(c *gin.Context) {
	userID, ok := paymentUserID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	reference := strings.TrimSpace(c.Param("reference"))
	if reference == "" {
		utils.Failure(c, 400, "invalid payment reference", nil)
		return
	}
	payment, err := h.service.Verify(c, userID, reference)
	if err == mongo.ErrNoDocuments || payment == nil {
		utils.Failure(c, 404, "payment not found", nil)
		return
	}
	if err != nil {
		utils.Failure(c, 502, err.Error(), nil)
		return
	}
	utils.Success(c, 200, "Payment verified successfully.", payment)
}
func (h *PaymentHandler) List(c *gin.Context) {
	userID, ok := paymentUserID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	q := repositories.PaymentQuery{Page: page, Limit: limit, Search: c.Query("search"), Status: c.Query("status")}
	if value := c.Query("from"); value != "" {
		q.From, _ = time.Parse("2006-01-02", value)
	}
	if value := c.Query("to"); value != "" {
		q.To, _ = time.Parse("2006-01-02", value)
		if !q.To.IsZero() {
			q.To = q.To.Add(24*time.Hour - time.Nanosecond)
		}
	}
	result, err := h.service.List(c, userID, q)
	if err != nil {
		utils.Failure(c, 400, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Payments retrieved successfully.", "data": result.Items, "pagination": gin.H{"page": result.Page, "limit": result.Limit, "total": result.Total, "totalPages": result.TotalPages}})
}
func (h *PaymentHandler) Get(c *gin.Context) {
	userID, ok := paymentUserID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		utils.Failure(c, 400, "invalid payment id", nil)
		return
	}
	payment, err := h.service.Get(c, userID, id)
	if err != nil {
		utils.Failure(c, 500, "unable to retrieve payment", nil)
		return
	}
	if payment == nil {
		utils.Failure(c, 404, "payment not found", nil)
		return
	}
	utils.Success(c, 200, "Payment retrieved successfully.", payment)
}
func (h *PaymentHandler) Summary(c *gin.Context) {
	userID, ok := paymentUserID(c)
	if !ok {
		utils.Failure(c, 401, "unauthorized", nil)
		return
	}
	summary, err := h.service.Summary(c, userID)
	if err != nil {
		utils.Failure(c, 500, "unable to retrieve payment summary", nil)
		return
	}
	utils.Success(c, 200, "Payment summary retrieved successfully.", summary)
}
func (h *PaymentHandler) Webhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.Failure(c, 400, "invalid webhook payload", nil)
		return
	}
	mac := hmac.New(sha512.New, []byte(h.paystackSecret))
	_, _ = mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	if h.paystackSecret == "" || !hmac.Equal([]byte(strings.ToLower(c.GetHeader("x-paystack-signature"))), []byte(expected)) {
		utils.Failure(c, 401, "invalid webhook signature", nil)
		return
	}
	var event struct {
		Event string `json:"event"`
		Data  struct {
			Reference string `json:"reference"`
			Status    string `json:"status"`
		} `json:"data"`
	}
	if json.Unmarshal(body, &event) != nil || strings.TrimSpace(event.Data.Reference) == "" {
		utils.Failure(c, 400, "invalid webhook payload", nil)
		return
	}
	payment, err := h.service.Webhook(c, event.Data.Reference, event.Data.Status)
	if err != nil {
		utils.Failure(c, 500, "unable to process webhook", nil)
		return
	}
	if payment == nil {
		c.Status(http.StatusOK)
		return
	}
	c.Status(http.StatusOK)
}
