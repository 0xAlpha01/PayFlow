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

type PaymentQuery struct {
	Page, Limit    int
	Search, Status string
	From, To       time.Time
}
type PaymentPage struct {
	Items      []models.Payment `json:"items"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	Total      int64            `json:"total"`
	TotalPages int              `json:"totalPages"`
}
type PaymentSummary struct {
	Successful       int64   `json:"successful"`
	Pending          int64   `json:"pending"`
	Failed           int64   `json:"failed"`
	Refunded         int64   `json:"refunded"`
	SuccessfulAmount float64 `json:"successfulAmount"`
	PendingAmount    float64 `json:"pendingAmount"`
	FailedAmount     float64 `json:"failedAmount"`
}
type PaymentRepository struct {
	db                                *mongo.Database
	payments, transactions, customers *mongo.Collection
}

func NewPaymentRepository(db *mongo.Database) (*PaymentRepository, error) {
	r := &PaymentRepository{db: db, payments: db.Collection("payments"), transactions: db.Collection("transactions"), customers: db.Collection("customers")}
	_, err := r.payments.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "createdAt", Value: -1}}},
		{Keys: bson.D{{Key: "customerId", Value: 1}}},
		{Keys: bson.D{{Key: "paystackReference", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)},
	})
	if err != nil {
		return nil, err
	}
	_, err = r.transactions.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bson.D{{Key: "paystackReference", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)})
	return r, err
}
func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	_, err := r.payments.InsertOne(ctx, payment)
	return err
}
func (r *PaymentRepository) Customer(ctx context.Context, userID, customerID primitive.ObjectID) (*models.Customer, error) {
	var c models.Customer
	err := r.customers.FindOne(ctx, bson.M{"_id": customerID, "userId": userID}).Decode(&c)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &c, err
}
func (r *PaymentRepository) ByID(ctx context.Context, userID, id primitive.ObjectID) (*models.Payment, error) {
	var p models.Payment
	err := r.payments.FindOne(ctx, bson.M{"_id": id, "userId": userID}).Decode(&p)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &p, err
}
func (r *PaymentRepository) ByReference(ctx context.Context, reference string) (*models.Payment, error) {
	var p models.Payment
	err := r.payments.FindOne(ctx, bson.M{"paystackReference": reference}).Decode(&p)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &p, err
}
func (r *PaymentRepository) List(ctx context.Context, userID primitive.ObjectID, q PaymentQuery) (PaymentPage, error) {
	filter := bson.M{"userId": userID}
	if q.Status != "" {
		filter["status"] = q.Status
	}
	if q.Search != "" {
		x := bson.M{"$regex": regexp.QuoteMeta(q.Search), "$options": "i"}
		filter["$or"] = bson.A{bson.M{"paystackReference": x}, bson.M{"currency": x}}
	}
	if !q.From.IsZero() || !q.To.IsZero() {
		dates := bson.M{}
		if !q.From.IsZero() {
			dates["$gte"] = q.From
		}
		if !q.To.IsZero() {
			dates["$lte"] = q.To
		}
		filter["createdAt"] = dates
	}
	total, err := r.payments.CountDocuments(ctx, filter)
	if err != nil {
		return PaymentPage{}, err
	}
	cur, err := r.payments.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetSkip(int64((q.Page-1)*q.Limit)).SetLimit(int64(q.Limit)))
	if err != nil {
		return PaymentPage{}, err
	}
	defer cur.Close(ctx)
	var items []models.Payment
	if err = cur.All(ctx, &items); err != nil {
		return PaymentPage{}, err
	}
	return PaymentPage{items, q.Page, q.Limit, total, int((total + int64(q.Limit) - 1) / int64(q.Limit))}, nil
}
func (r *PaymentRepository) Summary(ctx context.Context, userID primitive.ObjectID) (PaymentSummary, error) {
	cur, err := r.payments.Aggregate(ctx, mongo.Pipeline{{{Key: "$match", Value: bson.M{"userId": userID}}}, {{Key: "$group", Value: bson.M{"_id": "$status", "count": bson.M{"$sum": 1}, "amount": bson.M{"$sum": "$amount"}}}}})
	if err != nil {
		return PaymentSummary{}, err
	}
	defer cur.Close(ctx)
	var rows []struct {
		ID     string  `bson:"_id"`
		Count  int64   `bson:"count"`
		Amount float64 `bson:"amount"`
	}
	if err = cur.All(ctx, &rows); err != nil {
		return PaymentSummary{}, err
	}
	var s PaymentSummary
	for _, v := range rows {
		switch v.ID {
		case "successful":
			s.Successful, s.SuccessfulAmount = v.Count, v.Amount
		case "pending":
			s.Pending, s.PendingAmount = v.Count, v.Amount
		case "failed":
			s.Failed, s.FailedAmount = v.Count, v.Amount
		case "refunded":
			s.Refunded = v.Count
		}
	}
	return s, nil
}

// MarkSuccessful transitions a payment exactly once; only that transition creates a transaction and increments customer spending.
// The returned bool reports whether this call actually performed the transition, so callers
// (e.g. notifications) don't fire again on a repeated, already-successful verification.
func (r *PaymentRepository) MarkSuccessful(ctx context.Context, payment *models.Payment) (*models.Payment, bool, error) {
	if payment.Status == "successful" {
		return payment, false, nil
	}
	now := time.Now().UTC()
	result, err := r.payments.UpdateOne(ctx, bson.M{"_id": payment.ID, "status": bson.M{"$ne": "successful"}}, bson.M{"$set": bson.M{"status": "successful", "updatedAt": now}})
	if err != nil {
		return nil, false, err
	}
	if result.ModifiedCount == 0 {
		updated, err := r.ByReference(ctx, payment.PaystackReference)
		return updated, false, err
	}
	t := models.Transaction{ID: primitive.NewObjectID(), UserID: payment.UserID, CustomerID: payment.CustomerID, Amount: payment.Amount, Currency: payment.Currency, PaymentMethod: "paystack", Status: "successful", PaystackReference: payment.PaystackReference, CreatedAt: now, UpdatedAt: now}
	tx, err := r.transactions.UpdateOne(ctx, bson.M{"paystackReference": payment.PaystackReference}, bson.M{"$setOnInsert": t}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, false, err
	}
	var stored models.Transaction
	if err = r.transactions.FindOne(ctx, bson.M{"paystackReference": payment.PaystackReference}).Decode(&stored); err != nil {
		return nil, false, err
	}
	_, err = r.payments.UpdateOne(ctx, bson.M{"_id": payment.ID}, bson.M{"$set": bson.M{"transactionId": stored.ID, "updatedAt": now}})
	if err != nil {
		return nil, false, err
	}
	if tx.UpsertedCount == 1 {
		_, err = r.customers.UpdateOne(ctx, bson.M{"_id": payment.CustomerID, "userId": payment.UserID}, bson.M{"$inc": bson.M{"totalSpent": payment.Amount}, "$set": bson.M{"updatedAt": now}})
		if err != nil {
			return nil, false, err
		}
	}
	updated, err := r.ByReference(ctx, payment.PaystackReference)
	return updated, true, err
}
func (r *PaymentRepository) UpdateStatus(ctx context.Context, payment *models.Payment, status string) error {
	if payment.Status == "successful" {
		return nil
	}
	_, err := r.payments.UpdateOne(ctx, bson.M{"_id": payment.ID, "status": bson.M{"$ne": "successful"}}, bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now().UTC()}})
	return err
}
