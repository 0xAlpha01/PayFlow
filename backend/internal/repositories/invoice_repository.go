package repositories

import (
	"context"
	"errors"
	"github.com/ngwebwork/payflow/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"regexp"
	"time"
)

type InvoiceQuery struct {
	Page, Limit    int
	Search, Status string
}
type InvoicePage struct {
	Items      []models.Invoice `json:"items"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	Total      int64            `json:"total"`
	TotalPages int              `json:"totalPages"`
}
type InvoiceRepository struct{ invoices, customers *mongo.Collection }

func NewInvoiceRepository(db *mongo.Database) (*InvoiceRepository, error) {
	r := &InvoiceRepository{db.Collection("invoices"), db.Collection("customers")}
	_, err := r.invoices.Indexes().CreateMany(context.Background(), []mongo.IndexModel{{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "createdAt", Value: -1}}}, {Keys: bson.D{{Key: "customerId", Value: 1}}}, {Keys: bson.D{{Key: "userId", Value: 1}, {Key: "invoiceNumber", Value: 1}}, Options: options.Index().SetUnique(true)}})
	return r, err
}
func (r *InvoiceRepository) Customer(ctx context.Context, user, id primitive.ObjectID) (*models.Customer, error) {
	var c models.Customer
	err := r.customers.FindOne(ctx, bson.M{"_id": id, "userId": user}).Decode(&c)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &c, err
}
func (r *InvoiceRepository) Create(ctx context.Context, invoice *models.Invoice) error {
	_, err := r.invoices.InsertOne(ctx, invoice)
	return err
}
func (r *InvoiceRepository) Get(ctx context.Context, user, id primitive.ObjectID) (*models.Invoice, error) {
	var i models.Invoice
	err := r.invoices.FindOne(ctx, bson.M{"_id": id, "userId": user}).Decode(&i)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &i, err
}
func (r *InvoiceRepository) List(ctx context.Context, user primitive.ObjectID, q InvoiceQuery) (InvoicePage, error) {
	filter := bson.M{"userId": user}
	if q.Status != "" {
		filter["status"] = q.Status
	}
	if q.Search != "" {
		x := bson.M{"$regex": regexp.QuoteMeta(q.Search), "$options": "i"}
		filter["$or"] = bson.A{bson.M{"invoiceNumber": x}, bson.M{"currency": x}}
	}
	total, err := r.invoices.CountDocuments(ctx, filter)
	if err != nil {
		return InvoicePage{}, err
	}
	cur, err := r.invoices.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetSkip(int64((q.Page-1)*q.Limit)).SetLimit(int64(q.Limit)))
	if err != nil {
		return InvoicePage{}, err
	}
	defer cur.Close(ctx)
	var items []models.Invoice
	if err = cur.All(ctx, &items); err != nil {
		return InvoicePage{}, err
	}
	return InvoicePage{items, q.Page, q.Limit, total, int((total + int64(q.Limit) - 1) / int64(q.Limit))}, nil
}
func (r *InvoiceRepository) Update(ctx context.Context, user, id primitive.ObjectID, invoice *models.Invoice) error {
	invoice.UpdatedAt = time.Now().UTC()
	res, err := r.invoices.ReplaceOne(ctx, bson.M{"_id": id, "userId": user}, invoice)
	if err == nil && res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return err
}
func (r *InvoiceRepository) Delete(ctx context.Context, user, id primitive.ObjectID) error {
	res, err := r.invoices.DeleteOne(ctx, bson.M{"_id": id, "userId": user})
	if err == nil && res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return err
}
func (r *InvoiceRepository) SetStatus(ctx context.Context, user, id primitive.ObjectID, status string) error {
	res, err := r.invoices.UpdateOne(ctx, bson.M{"_id": id, "userId": user}, bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now().UTC()}})
	if err == nil && res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return err
}
