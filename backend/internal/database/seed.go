package database

import (
	"context"
	"time"

	"github.com/ngwebwork/payflow/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SeedCollections(ctx context.Context, db *mongo.Database, userID primitive.ObjectID) error {
	if count, err := db.Collection("customers").CountDocuments(ctx, bson.M{"userId": userID}); err != nil {
		return err
	} else if count > 0 {
		return nil
	}
	now := time.Now().UTC()
	customerID := primitive.NewObjectID()
	transactionID := primitive.NewObjectID()
	_, err := db.Collection("customers").InsertOne(ctx, models.Customer{ID: customerID, UserID: userID, FirstName: "Maya", LastName: "Okafor", Email: "maya@northstar.io", Phone: "+2348012345678", TotalSpent: 124500, Status: "active", CreatedAt: now.Add(-30 * 24 * time.Hour), UpdatedAt: now})
	if err != nil {
		return err
	}
	_, err = db.Collection("transactions").InsertMany(ctx, []interface{}{models.Transaction{ID: transactionID, UserID: userID, CustomerID: customerID, Amount: 85000, Currency: "NGN", PaymentMethod: "Card", Status: "successful", PaystackReference: "seed_trx_001", CreatedAt: now.Add(-2 * 24 * time.Hour), UpdatedAt: now}, models.Transaction{ID: primitive.NewObjectID(), UserID: userID, CustomerID: customerID, Amount: 39500, Currency: "NGN", PaymentMethod: "Bank transfer", Status: "pending", PaystackReference: "seed_trx_002", CreatedAt: now.Add(-24 * time.Hour), UpdatedAt: now}})
	if err != nil {
		return err
	}
	_, err = db.Collection("payments").InsertOne(ctx, models.Payment{ID: primitive.NewObjectID(), UserID: userID, CustomerID: customerID, TransactionID: transactionID, Amount: 85000, Currency: "NGN", Status: "successful", PaystackReference: "seed_pay_001", CreatedAt: now.Add(-2 * 24 * time.Hour), UpdatedAt: now})
	if err != nil {
		return err
	}
	_, err = db.Collection("invoices").InsertOne(ctx, models.Invoice{ID: primitive.NewObjectID(), UserID: userID, CustomerID: customerID, InvoiceNumber: "INV-SEED-001", Items: []models.InvoiceItem{{Description: "Growth plan", Quantity: 1, Price: 124500}}, Subtotal: 124500, Total: 124500, Currency: "NGN", Status: "paid", DueDate: now.Add(15 * 24 * time.Hour), CreatedAt: now.Add(-3 * 24 * time.Hour), UpdatedAt: now})
	if err != nil {
		return err
	}
	_, err = db.Collection("subscriptions").InsertOne(ctx, models.Subscription{ID: primitive.NewObjectID(), UserID: userID, CustomerID: customerID, Plan: "Growth", Amount: 25000, Currency: "NGN", BillingCycle: "monthly", Status: "active", NextPaymentDate: now.Add(28 * 24 * time.Hour), CreatedAt: now.Add(-20 * 24 * time.Hour), UpdatedAt: now})
	return err
}
