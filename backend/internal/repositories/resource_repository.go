package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
)

type ResourceRepository struct{ db *mongo.Database }
type PageResult struct {
	Items      []bson.M `json:"items"`
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	Total      int64    `json:"total"`
	TotalPages int      `json:"totalPages"`
}

func NewResourceRepository(db *mongo.Database) *ResourceRepository {
	return &ResourceRepository{db: db}
}
func (r *ResourceRepository) List(ctx context.Context, collection string, userID primitive.ObjectID, query map[string]string) (PageResult, error) {
	page, _ := strconv.Atoi(query["page"])
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(query["limit"])
	if limit < 1 || limit > 100 {
		limit = 20
	}
	filter := bson.M{"userId": userID}
	if query["status"] != "" {
		filter["status"] = query["status"]
	}
	if query["customerId"] != "" {
		id, err := primitive.ObjectIDFromHex(query["customerId"])
		if err == nil {
			filter["customerId"] = id
		}
	}
	if query["search"] != "" {
		expression := bson.M{"$regex": query["search"], "$options": "i"}
		filter["$or"] = bson.A{bson.M{"email": expression}, bson.M{"firstName": expression}, bson.M{"lastName": expression}, bson.M{"invoiceNumber": expression}, bson.M{"paystackReference": expression}}
	}
	collectionRef := r.db.Collection(collection)
	total, err := collectionRef.CountDocuments(ctx, filter)
	if err != nil {
		return PageResult{}, err
	}
	cursor, err := collectionRef.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)))
	if err != nil {
		return PageResult{}, err
	}
	defer cursor.Close(ctx)
	var items []bson.M
	if err := cursor.All(ctx, &items); err != nil {
		return PageResult{}, err
	}
	pages := int(total / int64(limit))
	if total%int64(limit) > 0 {
		pages++
	}
	return PageResult{Items: items, Page: page, Limit: limit, Total: total, TotalPages: pages}, nil
}
func (r *ResourceRepository) Get(ctx context.Context, collection string, userID, id primitive.ObjectID) (bson.M, error) {
	var item bson.M
	err := r.db.Collection(collection).FindOne(ctx, bson.M{"_id": id, "userId": userID}).Decode(&item)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return item, err
}
func (r *ResourceRepository) Insert(ctx context.Context, collection string, userID primitive.ObjectID, payload bson.M) (bson.M, error) {
	payload["_id"] = primitive.NewObjectID()
	payload["userId"] = userID
	now := time.Now().UTC()
	payload["createdAt"] = now
	payload["updatedAt"] = now
	_, err := r.db.Collection(collection).InsertOne(ctx, payload)
	return payload, err
}
func (r *ResourceRepository) Update(ctx context.Context, collection string, userID, id primitive.ObjectID, payload bson.M) error {
	delete(payload, "_id")
	delete(payload, "userId")
	delete(payload, "createdAt")
	payload["updatedAt"] = time.Now().UTC()
	result, err := r.db.Collection(collection).UpdateOne(ctx, bson.M{"_id": id, "userId": userID}, bson.M{"$set": payload})
	if err == nil && result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return err
}
func (r *ResourceRepository) Delete(ctx context.Context, collection string, userID, id primitive.ObjectID) error {
	result, err := r.db.Collection(collection).DeleteOne(ctx, bson.M{"_id": id, "userId": userID})
	if err == nil && result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return err
}
func ParseID(value string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(strings.TrimSpace(value))
}
