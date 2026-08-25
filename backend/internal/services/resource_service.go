package services

import (
	"context"
	"errors"
	"github.com/ngwebwork/payflow/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResourceService struct {
	repo *repositories.ResourceRepository
}

func NewResourceService(repo *repositories.ResourceRepository) *ResourceService {
	return &ResourceService{repo: repo}
}
func (s *ResourceService) List(ctx context.Context, collection string, userID primitive.ObjectID, query map[string]string) (repositories.PageResult, error) {
	return s.repo.List(ctx, collection, userID, query)
}
func (s *ResourceService) Get(ctx context.Context, collection string, userID, id primitive.ObjectID) (bson.M, error) {
	return s.repo.Get(ctx, collection, userID, id)
}
func (s *ResourceService) Create(ctx context.Context, collection string, userID primitive.ObjectID, payload bson.M) (bson.M, error) {
	if collection == "invoices" {
		if err := calculateInvoice(payload); err != nil {
			return nil, err
		}
	}
	return s.repo.Insert(ctx, collection, userID, payload)
}
func (s *ResourceService) Update(ctx context.Context, collection string, userID, id primitive.ObjectID, payload bson.M) error {
	if collection == "invoices" {
		if err := calculateInvoice(payload); err != nil {
			return err
		}
	}
	return s.repo.Update(ctx, collection, userID, id, payload)
}
func (s *ResourceService) Delete(ctx context.Context, collection string, userID, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, collection, userID, id)
}
func calculateInvoice(payload bson.M) error {
	items, ok := payload["items"].(primitive.A)
	if !ok {
		if values, valid := payload["items"].([]interface{}); valid {
			items = values
		} else {
			return errors.New("invoice items are required")
		}
	}
	var subtotal float64
	for _, value := range items {
		item, ok := value.(map[string]interface{})
		if !ok {
			if bsonItem, valid := value.(bson.M); valid {
				item = bsonItem
			} else {
				continue
			}
		}
		quantity, _ := item["quantity"].(float64)
		price, _ := item["price"].(float64)
		subtotal += quantity * price
	}
	tax, _ := payload["tax"].(float64)
	discount, _ := payload["discount"].(float64)
	payload["subtotal"] = subtotal
	payload["total"] = subtotal + tax - discount
	return nil
}
