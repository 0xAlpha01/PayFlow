package repositories

import (
	"context"

	"github.com/ngwebwork/payflow/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationPage struct {
	Items      []models.Notification `json:"items"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	Total      int64                 `json:"total"`
	TotalPages int                   `json:"totalPages"`
}
type NotificationRepository struct {
	notifications *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) (*NotificationRepository, error) {
	r := &NotificationRepository{notifications: db.Collection("notifications")}
	_, err := r.notifications.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "createdAt", Value: -1}}},
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "read", Value: 1}}},
	})
	return r, err
}
func (r *NotificationRepository) Create(ctx context.Context, n *models.Notification) error {
	_, err := r.notifications.InsertOne(ctx, n)
	return err
}
func (r *NotificationRepository) List(ctx context.Context, userID primitive.ObjectID, page, limit int) (NotificationPage, error) {
	filter := bson.M{"userId": userID}
	total, err := r.notifications.CountDocuments(ctx, filter)
	if err != nil {
		return NotificationPage{}, err
	}
	cursor, err := r.notifications.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)))
	if err != nil {
		return NotificationPage{}, err
	}
	defer cursor.Close(ctx)
	items := []models.Notification{}
	if err := cursor.All(ctx, &items); err != nil {
		return NotificationPage{}, err
	}
	pages := int((total + int64(limit) - 1) / int64(limit))
	return NotificationPage{Items: items, Page: page, Limit: limit, Total: total, TotalPages: pages}, nil
}
func (r *NotificationRepository) UnreadCount(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	return r.notifications.CountDocuments(ctx, bson.M{"userId": userID, "read": false})
}
func (r *NotificationRepository) MarkRead(ctx context.Context, userID, id primitive.ObjectID) error {
	result, err := r.notifications.UpdateOne(ctx, bson.M{"_id": id, "userId": userID}, bson.M{"$set": bson.M{"read": true}})
	if err == nil && result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return err
}
func (r *NotificationRepository) MarkAllRead(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.notifications.UpdateMany(ctx, bson.M{"userId": userID, "read": false}, bson.M{"$set": bson.M{"read": true}})
	return err
}
