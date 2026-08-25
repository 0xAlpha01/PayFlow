package repositories

import (
	"context"
	"time"

	"github.com/ngwebwork/payflow/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DashboardRepository struct{ db *mongo.Database }
type RevenuePoint struct {
	Date    string  `bson:"_id" json:"date"`
	Revenue float64 `bson:"revenue" json:"revenue"`
}

func NewDashboardRepository(db *mongo.Database) *DashboardRepository {
	return &DashboardRepository{db: db}
}
func (r *DashboardRepository) Transactions(ctx context.Context, userID primitive.ObjectID, since time.Time, limit int) ([]models.Transaction, error) {
	cursor, err := r.db.Collection("transactions").Find(ctx, bson.M{"userId": userID, "createdAt": bson.M{"$gte": since}}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var result []models.Transaction
	return result, cursor.All(ctx, &result)
}
func (r *DashboardRepository) AllTransactions(ctx context.Context, userID primitive.ObjectID) ([]models.Transaction, error) {
	cursor, err := r.db.Collection("transactions").Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var result []models.Transaction
	return result, cursor.All(ctx, &result)
}
func (r *DashboardRepository) Customers(ctx context.Context, userID primitive.ObjectID, limit int) ([]models.Customer, error) {
	cursor, err := r.db.Collection("customers").Find(ctx, bson.M{"userId": userID}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var result []models.Customer
	return result, cursor.All(ctx, &result)
}
func (r *DashboardRepository) Customer(ctx context.Context, userID, customerID primitive.ObjectID) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Collection("customers").FindOne(ctx, bson.M{"_id": customerID, "userId": userID}).Decode(&customer)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &customer, err
}
func (r *DashboardRepository) Count(ctx context.Context, collection string, userID primitive.ObjectID, filter bson.M) (int64, error) {
	filter["userId"] = userID
	return r.db.Collection(collection).CountDocuments(ctx, filter)
}
func (r *DashboardRepository) Sum(ctx context.Context, collection string, userID primitive.ObjectID, filter bson.M, field string) (float64, error) {
	filter["userId"] = userID
	rows, err := r.db.Collection(collection).Aggregate(ctx, mongo.Pipeline{{{Key: "$match", Value: filter}}, {{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$" + field}}}}})
	if err != nil {
		return 0, err
	}
	defer rows.Close(ctx)
	var result []struct {
		Total float64 `bson:"total"`
	}
	if err := rows.All(ctx, &result); err != nil || len(result) == 0 {
		return 0, err
	}
	return result[0].Total, nil
}
func (r *DashboardRepository) Revenue(ctx context.Context, userID primitive.ObjectID, since time.Time) ([]RevenuePoint, error) {
	rows, err := r.db.Collection("transactions").Aggregate(ctx, mongo.Pipeline{{{Key: "$match", Value: bson.M{"userId": userID, "status": "successful", "createdAt": bson.M{"$gte": since}}}}, {{Key: "$group", Value: bson.M{"_id": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$createdAt"}}, "revenue": bson.M{"$sum": "$amount"}}}}, {{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}}})
	if err != nil {
		return nil, err
	}
	defer rows.Close(ctx)
	var result []RevenuePoint
	return result, rows.All(ctx, &result)
}
