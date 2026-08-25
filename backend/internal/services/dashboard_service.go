package services

import (
	"context"
	"strings"
	"time"

	"github.com/ngwebwork/payflow/internal/models"
	"github.com/ngwebwork/payflow/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DashboardService struct {
	repo *repositories.DashboardRepository
}
type DashboardSummary struct {
	TotalRevenue           float64 `json:"totalRevenue"`
	AvailableBalance       float64 `json:"availableBalance"`
	PendingPayments        float64 `json:"pendingPayments"`
	TotalCustomers         int64   `json:"totalCustomers"`
	SuccessfulTransactions int64   `json:"successfulTransactions"`
	FailedTransactions     int64   `json:"failedTransactions"`
	ActiveSubscriptions    int64   `json:"activeSubscriptions"`
	OutstandingInvoices    int64   `json:"outstandingInvoices"`
}
type DashboardTransaction struct {
	ID            primitive.ObjectID `json:"id"`
	CustomerName  string             `json:"customerName"`
	Amount        float64            `json:"amount"`
	Currency      string             `json:"currency"`
	PaymentMethod string             `json:"paymentMethod"`
	Status        string             `json:"status"`
	CreatedAt     time.Time          `json:"createdAt"`
}
type DashboardCustomer struct {
	ID          primitive.ObjectID `json:"id"`
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	AmountSpent float64            `json:"amountSpent"`
	LastPayment *time.Time         `json:"lastPayment,omitempty"`
	Status      string             `json:"status"`
}

func NewDashboardService(repo *repositories.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}
func (s *DashboardService) Summary(ctx context.Context, userID primitive.ObjectID) (DashboardSummary, error) {
	successful, err := s.repo.Sum(ctx, "transactions", userID, bson.M{"status": "successful"}, "amount")
	if err != nil {
		return DashboardSummary{}, err
	}
	pending, err := s.repo.Sum(ctx, "transactions", userID, bson.M{"status": "pending"}, "amount")
	if err != nil {
		return DashboardSummary{}, err
	}
	customers, err := s.repo.Count(ctx, "customers", userID, bson.M{})
	if err != nil {
		return DashboardSummary{}, err
	}
	successCount, err := s.repo.Count(ctx, "transactions", userID, bson.M{"status": "successful"})
	if err != nil {
		return DashboardSummary{}, err
	}
	failedCount, err := s.repo.Count(ctx, "transactions", userID, bson.M{"status": "failed"})
	if err != nil {
		return DashboardSummary{}, err
	}
	subscriptions, err := s.repo.Count(ctx, "subscriptions", userID, bson.M{"status": "active"})
	if err != nil {
		return DashboardSummary{}, err
	}
	invoices, err := s.repo.Count(ctx, "invoices", userID, bson.M{"status": bson.M{"$in": []string{"pending", "overdue"}}})
	if err != nil {
		return DashboardSummary{}, err
	}
	return DashboardSummary{TotalRevenue: successful, AvailableBalance: successful - pending, PendingPayments: pending, TotalCustomers: customers, SuccessfulTransactions: successCount, FailedTransactions: failedCount, ActiveSubscriptions: subscriptions, OutstandingInvoices: invoices}, nil
}
func (s *DashboardService) Revenue(ctx context.Context, userID primitive.ObjectID, period string) ([]repositories.RevenuePoint, error) {
	duration := map[string]time.Duration{"7d": 7 * 24 * time.Hour, "30d": 30 * 24 * time.Hour, "3m": 90 * 24 * time.Hour, "6m": 180 * 24 * time.Hour, "1y": 365 * 24 * time.Hour}[period]
	if duration == 0 {
		duration = 30 * 24 * time.Hour
	}
	return s.repo.Revenue(ctx, userID, time.Now().UTC().Add(-duration))
}
func (s *DashboardService) Transactions(ctx context.Context, userID primitive.ObjectID) ([]models.Transaction, error) {
	return s.repo.Transactions(ctx, userID, time.Time{}, 10)
}
func (s *DashboardService) Customers(ctx context.Context, userID primitive.ObjectID) ([]models.Customer, error) {
	return s.repo.Customers(ctx, userID, 10)
}
func (s *DashboardService) RecentTransactions(ctx context.Context, userID primitive.ObjectID) ([]DashboardTransaction, error) {
	rows, err := s.Transactions(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]DashboardTransaction, 0, len(rows))
	for _, row := range rows {
		customerName := "Unknown customer"
		customer, lookupErr := s.repo.Customer(ctx, userID, row.CustomerID)
		if lookupErr != nil {
			return nil, lookupErr
		}
		if customer != nil {
			customerName = strings.TrimSpace(customer.FirstName + " " + customer.LastName)
		}
		result = append(result, DashboardTransaction{ID: row.ID, CustomerName: customerName, Amount: row.Amount, Currency: row.Currency, PaymentMethod: row.PaymentMethod, Status: row.Status, CreatedAt: row.CreatedAt})
	}
	return result, nil
}
func (s *DashboardService) RecentCustomers(ctx context.Context, userID primitive.ObjectID) ([]DashboardCustomer, error) {
	rows, err := s.Customers(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]DashboardCustomer, 0, len(rows))
	for _, row := range rows {
		result = append(result, DashboardCustomer{ID: row.ID, Name: strings.TrimSpace(row.FirstName + " " + row.LastName), Email: row.Email, AmountSpent: row.TotalSpent, Status: row.Status})
	}
	return result, nil
}
