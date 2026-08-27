package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ngwebwork/payflow/internal/models"
	"github.com/ngwebwork/payflow/internal/repositories"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
}
type CustomerService struct {
	repo *repositories.CustomerRepository
}

func NewCustomerService(repo *repositories.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}
func (s *CustomerService) validate(input *CustomerInput) error {
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if err := utils.Required(input.FirstName, input.LastName, input.Email, input.Phone); err != nil {
		return err
	}
	if err := utils.Email(input.Email); err != nil {
		return err
	}
	phone, err := utils.NormalizeNigerianPhone(input.Phone)
	if err != nil {
		return err
	}
	input.Phone = phone
	if input.Status != "" && input.Status != "active" && input.Status != "inactive" {
		return errors.New("status must be active or inactive")
	}
	return nil
}
func (s *CustomerService) Create(ctx context.Context, userID primitive.ObjectID, input CustomerInput) (*models.Customer, error) {
	if err := s.validate(&input); err != nil {
		return nil, err
	}
	duplicate, err := s.repo.Exists(ctx, userID, "email", input.Email, primitive.NilObjectID)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, errors.New("a customer with this email already exists")
	}
	duplicate, err = s.repo.Exists(ctx, userID, "phone", input.Phone, primitive.NilObjectID)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, errors.New("a customer with this phone already exists")
	}
	now := time.Now().UTC()
	customer := &models.Customer{ID: primitive.NewObjectID(), UserID: userID, FirstName: input.FirstName, LastName: input.LastName, Email: input.Email, Phone: input.Phone, TotalSpent: 0, Status: "active", CreatedAt: now, UpdatedAt: now}
	if err := s.repo.Create(ctx, customer); err != nil {
		return nil, err
	}
	return customer, nil
}
func (s *CustomerService) List(ctx context.Context, userID primitive.ObjectID, query repositories.CustomerQuery) (repositories.CustomerPage, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 20
	}
	if query.Status != "" && query.Status != "active" && query.Status != "inactive" {
		return repositories.CustomerPage{}, errors.New("invalid status")
	}
	allowed := map[string]bool{"createdAt": true, "firstName": true, "totalSpent": true}
	if !allowed[query.SortBy] {
		query.SortBy = "createdAt"
	}
	if query.SortOrder != "asc" && query.SortOrder != "desc" {
		query.SortOrder = "desc"
	}
	return s.repo.List(ctx, userID, query)
}
func (s *CustomerService) Get(ctx context.Context, userID, id primitive.ObjectID) (*models.Customer, error) {
	return s.repo.ByID(ctx, userID, id)
}
func (s *CustomerService) Update(ctx context.Context, userID, id primitive.ObjectID, input CustomerInput) (*models.Customer, error) {
	current, err := s.repo.ByID(ctx, userID, id)
	if err != nil || current == nil {
		return nil, err
	}
	if input.FirstName == "" {
		input.FirstName = current.FirstName
	}
	if input.LastName == "" {
		input.LastName = current.LastName
	}
	if input.Email == "" {
		input.Email = current.Email
	}
	if input.Phone == "" {
		input.Phone = current.Phone
	}
	if input.Status == "" {
		input.Status = current.Status
	}
	if err := s.validate(&input); err != nil {
		return nil, err
	}
	duplicate, err := s.repo.Exists(ctx, userID, "email", input.Email, id)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, errors.New("a customer with this email already exists")
	}
	duplicate, err = s.repo.Exists(ctx, userID, "phone", input.Phone, id)
	if err != nil {
		return nil, err
	}
	if duplicate {
		return nil, errors.New("a customer with this phone already exists")
	}
	if err := s.repo.Update(ctx, userID, id, bson.M{"firstName": input.FirstName, "lastName": input.LastName, "email": input.Email, "phone": input.Phone, "status": input.Status}); err != nil {
		return nil, err
	}
	return s.repo.ByID(ctx, userID, id)
}
func (s *CustomerService) Delete(ctx context.Context, userID, id primitive.ObjectID) error {
	customer, err := s.repo.ByID(ctx, userID, id)
	if err != nil {
		return err
	}
	if customer == nil {
		return mongo.ErrNoDocuments
	}
	history, err := s.repo.HasFinancialHistory(ctx, userID, id)
	if err != nil {
		return err
	}
	if history {
		return errors.New("customer cannot be deleted because they have financial history")
	}
	return s.repo.Delete(ctx, userID, id)
}
func (s *CustomerService) Payments(ctx context.Context, userID, id primitive.ObjectID) ([]bson.M, error) {
	customer, err := s.Get(ctx, userID, id)
	if err != nil || customer == nil {
		return nil, mongo.ErrNoDocuments
	}
	return s.repo.Payments(ctx, userID, id)
}
func (s *CustomerService) Transactions(ctx context.Context, userID, id primitive.ObjectID) ([]bson.M, error) {
	customer, err := s.Get(ctx, userID, id)
	if err != nil || customer == nil {
		return nil, mongo.ErrNoDocuments
	}
	return s.repo.Transactions(ctx, userID, id)
}
