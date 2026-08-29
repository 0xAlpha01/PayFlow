package services

import (
	"context"
	"errors"
	"github.com/ngwebwork/payflow/internal/models"
	"github.com/ngwebwork/payflow/internal/repositories"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"strings"
	"time"
)

type InitializePaymentInput struct {
	CustomerID string  `json:"customerId"`
	Amount     float64 `json:"amount"`
	Email      string  `json:"email"`
}
type PaymentService struct {
	repo     *repositories.PaymentRepository
	paystack *PaystackService
}

func NewPaymentService(repo *repositories.PaymentRepository, paystack *PaystackService) *PaymentService {
	return &PaymentService{repo, paystack}
}
func (s *PaymentService) Initialize(ctx context.Context, userID primitive.ObjectID, input InitializePaymentInput) (*models.Payment, PaystackInitializeResponse, error) {
	input.CustomerID = strings.TrimSpace(input.CustomerID)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if err := utils.Required(input.CustomerID, input.Email); err != nil {
		return nil, PaystackInitializeResponse{}, err
	}
	if err := utils.Email(input.Email); err != nil {
		return nil, PaystackInitializeResponse{}, err
	}
	if input.Amount <= 0 || math.IsNaN(input.Amount) || math.IsInf(input.Amount, 0) {
		return nil, PaystackInitializeResponse{}, errors.New("amount must be greater than zero")
	}
	customerID, err := primitive.ObjectIDFromHex(input.CustomerID)
	if err != nil {
		return nil, PaystackInitializeResponse{}, errors.New("invalid customer id")
	}
	customer, err := s.repo.Customer(ctx, userID, customerID)
	if err != nil {
		return nil, PaystackInitializeResponse{}, err
	}
	if customer == nil {
		return nil, PaystackInitializeResponse{}, mongo.ErrNoDocuments
	}
	now := time.Now().UTC()
	reference := "pf_" + primitive.NewObjectID().Hex()
	p := &models.Payment{ID: primitive.NewObjectID(), UserID: userID, CustomerID: customerID, Amount: input.Amount, Currency: "NGN", Status: "pending", PaystackReference: reference, CreatedAt: now, UpdatedAt: now}
	if err = s.repo.Create(ctx, p); err != nil {
		return nil, PaystackInitializeResponse{}, err
	}
	amountKobo := int64(math.Round(input.Amount * 100))
	result, err := s.paystack.Initialize(ctx, input.Email, reference, amountKobo)
	if err != nil {
		return p, PaystackInitializeResponse{}, err
	}
	return p, result, nil
}
func (s *PaymentService) Verify(ctx context.Context, userID primitive.ObjectID, reference string) (*models.Payment, error) {
	p, err := s.repo.ByReference(ctx, reference)
	if err != nil || p == nil {
		return p, err
	}
	if p.UserID != userID {
		return nil, mongo.ErrNoDocuments
	}
	if p.Status == "successful" {
		return p, nil
	}
	status, err := s.paystack.Verify(ctx, reference)
	if err != nil {
		return nil, err
	}
	if status == "success" {
		return s.repo.MarkSuccessful(ctx, p)
	}
	local := "pending"
	if status == "failed" || status == "abandoned" {
		local = "failed"
	}
	if err = s.repo.UpdateStatus(ctx, p, local); err != nil {
		return nil, err
	}
	return s.repo.ByReference(ctx, reference)
}
func (s *PaymentService) Webhook(ctx context.Context, reference, status string) (*models.Payment, error) {
	p, err := s.repo.ByReference(ctx, reference)
	if err != nil || p == nil {
		return p, err
	}
	if status == "success" {
		return s.repo.MarkSuccessful(ctx, p)
	}
	local := "pending"
	if status == "failed" || status == "abandoned" {
		local = "failed"
	}
	return p, s.repo.UpdateStatus(ctx, p, local)
}
func (s *PaymentService) List(ctx context.Context, userID primitive.ObjectID, q repositories.PaymentQuery) (repositories.PaymentPage, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	if q.Status != "" && q.Status != "pending" && q.Status != "successful" && q.Status != "failed" && q.Status != "refunded" {
		return repositories.PaymentPage{}, errors.New("invalid status")
	}
	return s.repo.List(ctx, userID, q)
}
func (s *PaymentService) Get(ctx context.Context, userID, id primitive.ObjectID) (*models.Payment, error) {
	return s.repo.ByID(ctx, userID, id)
}
func (s *PaymentService) Summary(ctx context.Context, userID primitive.ObjectID) (repositories.PaymentSummary, error) {
	return s.repo.Summary(ctx, userID)
}
