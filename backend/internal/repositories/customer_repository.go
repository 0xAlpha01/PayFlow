package repositories

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/ngwebwork/payflow/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomerQuery struct {
	Page, Limit                       int
	Search, Status, SortBy, SortOrder string
	CreatedFrom, CreatedTo            time.Time
}
type CustomerPage struct {
	Items      []models.Customer `json:"items"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Total      int64             `json:"total"`
	TotalPages int               `json:"totalPages"`
}
type CustomerRepository struct {
	db        *mongo.Database
	customers *mongo.Collection
}

func NewCustomerRepository(db *mongo.Database) (*CustomerRepository, error) {
	r := &CustomerRepository{db: db, customers: db.Collection("customers")}
	_, err := r.customers.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "phone", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)},
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "createdAt", Value: -1}}}, {Keys: bson.D{{Key: "userId", Value: 1}, {Key: "status", Value: 1}}},
	})
	return r, err
}
func (r *CustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	_, err := r.customers.InsertOne(ctx, customer)
	return err
}
func (r *CustomerRepository) ByID(ctx context.Context, userID, id primitive.ObjectID) (*models.Customer, error) {
	var customer models.Customer
	err := r.customers.FindOne(ctx, bson.M{"_id": id, "userId": userID}).Decode(&customer)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &customer, err
}
func (r *CustomerRepository) Exists(ctx context.Context, userID primitive.ObjectID, field, value string, except primitive.ObjectID) (bool, error) {
	filter := bson.M{"userId": userID, field: value}
	if !except.IsZero() {
		filter["_id"] = bson.M{"$ne": except}
	}
	n, err := r.customers.CountDocuments(ctx, filter)
	return n > 0, err
}
func (r *CustomerRepository) List(ctx context.Context, userID primitive.ObjectID, query CustomerQuery) (CustomerPage, error) {
	filter := bson.M{"userId": userID}
	if query.Status != "" {
		filter["status"] = query.Status
	}
	if query.Search != "" {
		expression := bson.M{"$regex": regexp.QuoteMeta(query.Search), "$options": "i"}
		filter["$or"] = bson.A{bson.M{"firstName": expression}, bson.M{"lastName": expression}, bson.M{"email": expression}, bson.M{"phone": expression}}
	}
	if !query.CreatedFrom.IsZero() || !query.CreatedTo.IsZero() {
		dates := bson.M{}
		if !query.CreatedFrom.IsZero() {
			dates["$gte"] = query.CreatedFrom
		}
		if !query.CreatedTo.IsZero() {
			dates["$lte"] = query.CreatedTo
		}
		filter["createdAt"] = dates
	}
	total, err := r.customers.CountDocuments(ctx, filter)
	if err != nil {
		return CustomerPage{}, err
	}
	direction := -1
	if query.SortOrder == "asc" {
		direction = 1
	}
	cursor, err := r.customers.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: query.SortBy, Value: direction}}).SetSkip(int64((query.Page-1)*query.Limit)).SetLimit(int64(query.Limit)))
	if err != nil {
		return CustomerPage{}, err
	}
	defer cursor.Close(ctx)
	var items []models.Customer
	if err = cursor.All(ctx, &items); err != nil {
		return CustomerPage{}, err
	}
	pages := int((total + int64(query.Limit) - 1) / int64(query.Limit))
	return CustomerPage{Items: items, Page: query.Page, Limit: query.Limit, Total: total, TotalPages: pages}, nil
}
func (r *CustomerRepository) Update(ctx context.Context, userID, id primitive.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now().UTC()
	result, err := r.customers.UpdateOne(ctx, bson.M{"_id": id, "userId": userID}, bson.M{"$set": update})
	if err == nil && result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return err
}
func (r *CustomerRepository) Delete(ctx context.Context, userID, id primitive.ObjectID) error {
	result, err := r.customers.DeleteOne(ctx, bson.M{"_id": id, "userId": userID})
	if err == nil && result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return err
}
func (r *CustomerRepository) HasFinancialHistory(ctx context.Context, userID, customerID primitive.ObjectID) (bool, error) {
	filter := bson.M{"userId": userID, "customerId": customerID}
	payments, err := r.db.Collection("payments").CountDocuments(ctx, filter)
	if err != nil || payments > 0 {
		return payments > 0, err
	}
	transactions, err := r.db.Collection("transactions").CountDocuments(ctx, filter)
	return transactions > 0, err
}
func (r *CustomerRepository) Payments(ctx context.Context, userID, customerID primitive.ObjectID) ([]bson.M, error) {
	return r.history(ctx, "payments", userID, customerID)
}
func (r *CustomerRepository) Transactions(ctx context.Context, userID, customerID primitive.ObjectID) ([]bson.M, error) {
	return r.history(ctx, "transactions", userID, customerID)
}
func (r *CustomerRepository) history(ctx context.Context, collection string, userID, customerID primitive.ObjectID) ([]bson.M, error) {
	cursor, err := r.db.Collection(collection).Find(ctx, bson.M{"userId": userID, "customerId": customerID}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var items []bson.M
	return items, cursor.All(ctx, &items)
}
