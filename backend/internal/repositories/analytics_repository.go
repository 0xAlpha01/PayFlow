package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type AnalyticsRepository struct{ db *mongo.Database }

func NewAnalyticsRepository(db *mongo.Database) *AnalyticsRepository { return &AnalyticsRepository{db} }
func (r *AnalyticsRepository) Count(c context.Context, col string, u primitive.ObjectID, f bson.M) (int64, error) {
	f["userId"] = u
	return r.db.Collection(col).CountDocuments(c, f)
}
func (r *AnalyticsRepository) Sum(c context.Context, col string, u primitive.ObjectID, f bson.M) (float64, error) {
	f["userId"] = u
	cur, e := r.db.Collection(col).Aggregate(c, mongo.Pipeline{{{Key: "$match", Value: f}}, {{Key: "$group", Value: bson.M{"_id": nil, "amount": bson.M{"$sum": "$amount"}}}}})
	if e != nil {
		return 0, e
	}
	defer cur.Close(c)
	var x []struct {
		Amount float64 `bson:"amount"`
	}
	e = cur.All(c, &x)
	if len(x) > 0 {
		return x[0].Amount, e
	}
	return 0, e
}
func (r *AnalyticsRepository) Series(c context.Context, col string, u primitive.ObjectID, from time.Time) ([]bson.M, error) {
	cur, e := r.db.Collection(col).Aggregate(c, mongo.Pipeline{{{Key: "$match", Value: bson.M{"userId": u, "status": "successful", "createdAt": bson.M{"$gte": from}}}}, {{Key: "$group", Value: bson.M{"_id": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$createdAt"}}, "revenue": bson.M{"$sum": "$amount"}}}}, {{Key: "$sort", Value: bson.M{"_id": 1}}}})
	if e != nil {
		return nil, e
	}
	defer cur.Close(c)
	var x []bson.M
	e = cur.All(c, &x)
	for i := range x {
		x[i]["date"] = x[i]["_id"]
		delete(x[i], "_id")
	}
	return x, e
}
