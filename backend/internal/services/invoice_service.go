package services

import (
	"context"
	"errors"
	"github.com/ngwebwork/payflow/internal/models"
	"github.com/ngwebwork/payflow/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"strings"
	"time"
)

type InvoiceInput struct {
	CustomerID    string               `json:"customerId"`
	InvoiceNumber string               `json:"invoiceNumber"`
	Items         []models.InvoiceItem `json:"items"`
	Tax           float64              `json:"tax"`
	Discount      float64              `json:"discount"`
	Currency      string               `json:"currency"`
	DueDate       string               `json:"dueDate"`
	Notes         string               `json:"notes"`
}
type InvoiceService struct {
	repo  *repositories.InvoiceRepository
	email *EmailService
}

func NewInvoiceService(repo *repositories.InvoiceRepository, email *EmailService) *InvoiceService {
	return &InvoiceService{repo, email}
}
func (s *InvoiceService) invoice(ctx context.Context, user primitive.ObjectID, input InvoiceInput, existing *models.Invoice) (*models.Invoice, error) {
	input.CustomerID = strings.TrimSpace(input.CustomerID)
	input.InvoiceNumber = strings.TrimSpace(input.InvoiceNumber)
	if input.CustomerID == "" || input.InvoiceNumber == "" || len(input.Items) == 0 {
		return nil, errors.New("customer, invoice number, and at least one item are required")
	}
	customerID, err := primitive.ObjectIDFromHex(input.CustomerID)
	if err != nil {
		return nil, errors.New("invalid customer id")
	}
	customer, err := s.repo.Customer(ctx, user, customerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, mongo.ErrNoDocuments
	}
	due, err := time.Parse("2006-01-02", input.DueDate)
	if err != nil {
		return nil, errors.New("invalid due date")
	}
	subtotal := 0.0
	for index := range input.Items {
		item := &input.Items[index]
		item.Description = strings.TrimSpace(item.Description)
		if item.Description == "" || item.Quantity <= 0 || item.Price < 0 {
			return nil, errors.New("each item requires a description, positive quantity, and valid unit price")
		}
		subtotal += float64(item.Quantity) * item.Price
	}
	if input.Tax < 0 || input.Discount < 0 || input.Tax > subtotal || input.Discount > subtotal+input.Tax {
		return nil, errors.New("invalid tax or discount")
	}
	currency := strings.ToUpper(strings.TrimSpace(input.Currency))
	if currency == "" {
		currency = "NGN"
	}
	if currency != "NGN" {
		return nil, errors.New("currency must be NGN")
	}
	now := time.Now().UTC()
	invoice := &models.Invoice{ID: primitive.NewObjectID(), UserID: user, CustomerID: customerID, InvoiceNumber: input.InvoiceNumber, Items: input.Items, Subtotal: subtotal, Tax: input.Tax, Discount: input.Discount, Total: math.Round((subtotal+input.Tax-input.Discount)*100) / 100, Currency: currency, Status: "draft", DueDate: due.UTC(), Notes: strings.TrimSpace(input.Notes), CreatedAt: now, UpdatedAt: now}
	if existing != nil {
		invoice.ID = existing.ID
		invoice.Status = existing.Status
		invoice.CreatedAt = existing.CreatedAt
	}
	return invoice, nil
}
func (s *InvoiceService) Create(ctx context.Context, user primitive.ObjectID, input InvoiceInput) (*models.Invoice, error) {
	i, err := s.invoice(ctx, user, input, nil)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Create(ctx, i); err != nil {
		return nil, err
	}
	return i, nil
}
func (s *InvoiceService) Update(ctx context.Context, user, id primitive.ObjectID, input InvoiceInput) (*models.Invoice, error) {
	current, err := s.repo.Get(ctx, user, id)
	if err != nil || current == nil {
		return current, err
	}
	i, err := s.invoice(ctx, user, input, current)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Update(ctx, user, id, i); err != nil {
		return nil, err
	}
	return i, nil
}
func (s *InvoiceService) List(ctx context.Context, user primitive.ObjectID, q repositories.InvoiceQuery) (repositories.InvoicePage, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	if q.Status != "" && q.Status != "draft" && q.Status != "pending" && q.Status != "paid" && q.Status != "overdue" && q.Status != "cancelled" {
		return repositories.InvoicePage{}, errors.New("invalid status")
	}
	return s.repo.List(ctx, user, q)
}
func (s *InvoiceService) Get(ctx context.Context, user, id primitive.ObjectID) (*models.Invoice, error) {
	return s.repo.Get(ctx, user, id)
}
func (s *InvoiceService) Delete(ctx context.Context, user, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, user, id)
}
func (s *InvoiceService) Send(ctx context.Context, user, id primitive.ObjectID) error {
	i, err := s.repo.Get(ctx, user, id)
	if err != nil || i == nil {
		return mongo.ErrNoDocuments
	}
	customer, err := s.repo.Customer(ctx, user, i.CustomerID)
	if err != nil || customer == nil {
		return mongo.ErrNoDocuments
	}
	if err = s.email.SendInvoice(customer.Email, i.InvoiceNumber, i.Total, i.Currency, i.DueDate.Format("2006-01-02")); err != nil {
		return err
	}
	return s.repo.SetStatus(ctx, user, id, "pending")
}
func (s *InvoiceService) MarkPaid(ctx context.Context, user, id primitive.ObjectID) error {
	return s.repo.SetStatus(ctx, user, id, "paid")
}
