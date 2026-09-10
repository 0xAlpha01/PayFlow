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

type ConversationSummary struct {
	Conversation models.Conversation `json:"conversation"`
	Peer         models.User         `json:"peer"`
	UnreadCount  int64               `json:"unreadCount"`
}
type MessagePage struct {
	Items      []models.ChatMessage `json:"items"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	Total      int64                `json:"total"`
	TotalPages int                  `json:"totalPages"`
}
type ChatRepository struct {
	users         *mongo.Collection
	conversations *mongo.Collection
	messages      *mongo.Collection
}

func NewChatRepository(db *mongo.Database) (*ChatRepository, error) {
	r := &ChatRepository{users: db.Collection("users"), conversations: db.Collection("conversations"), messages: db.Collection("messages")}
	_, err := r.conversations.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "pairKey", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "participantIds", Value: 1}, {Key: "lastMessageAt", Value: -1}}},
	})
	if err != nil {
		return nil, err
	}
	_, err = r.messages.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "conversationId", Value: 1}, {Key: "createdAt", Value: 1}}},
		{Keys: bson.D{{Key: "receiverId", Value: 1}, {Key: "read", Value: 1}}},
	})
	return r, err
}
func pairKey(a, b primitive.ObjectID) string {
	if a.Hex() <= b.Hex() {
		return a.Hex() + "_" + b.Hex()
	}
	return b.Hex() + "_" + a.Hex()
}

// SearchUsers returns other registered accounts matching a name/email prefix, safe fields only.
func (r *ChatRepository) SearchUsers(ctx context.Context, currentUserID primitive.ObjectID, query string, limit int) ([]models.User, error) {
	filter := bson.M{"_id": bson.M{"$ne": currentUserID}}
	if query != "" {
		expression := bson.M{"$regex": regexp.QuoteMeta(query), "$options": "i"}
		filter["$or"] = bson.A{bson.M{"firstName": expression}, bson.M{"lastName": expression}, bson.M{"email": expression}}
	}
	projection := bson.M{"firstName": 1, "lastName": 1, "email": 1, "businessName": 1}
	cursor, err := r.users.Find(ctx, filter, options.Find().SetProjection(projection).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	items := []models.User{}
	return items, cursor.All(ctx, &items)
}
func (r *ChatRepository) User(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var u models.User
	err := r.users.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"firstName": 1, "lastName": 1, "email": 1, "businessName": 1})).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &u, err
}
func (r *ChatRepository) GetOrCreateConversation(ctx context.Context, userA, userB primitive.ObjectID) (*models.Conversation, error) {
	key := pairKey(userA, userB)
	now := time.Now().UTC()
	after := options.After
	result := r.conversations.FindOneAndUpdate(ctx,
		bson.M{"pairKey": key},
		bson.M{"$setOnInsert": bson.M{"_id": primitive.NewObjectID(), "participantIds": []primitive.ObjectID{userA, userB}, "pairKey": key, "createdAt": now, "updatedAt": now}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(after),
	)
	var conversation models.Conversation
	if err := result.Decode(&conversation); err != nil {
		return nil, err
	}
	return &conversation, nil
}
func (r *ChatRepository) ConversationByID(ctx context.Context, id primitive.ObjectID) (*models.Conversation, error) {
	var c models.Conversation
	err := r.conversations.FindOne(ctx, bson.M{"_id": id}).Decode(&c)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &c, err
}
func (r *ChatRepository) ListConversations(ctx context.Context, userID primitive.ObjectID) ([]ConversationSummary, error) {
	cursor, err := r.conversations.Find(ctx, bson.M{"participantIds": userID}, options.Find().SetSort(bson.D{{Key: "lastMessageAt", Value: -1}, {Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var conversations []models.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		return nil, err
	}
	summaries := make([]ConversationSummary, 0, len(conversations))
	for _, conversation := range conversations {
		var peerID primitive.ObjectID
		for _, participant := range conversation.ParticipantIDs {
			if participant != userID {
				peerID = participant
			}
		}
		peer, err := r.User(ctx, peerID)
		if err != nil {
			return nil, err
		}
		if peer == nil {
			continue
		}
		unread, err := r.messages.CountDocuments(ctx, bson.M{"conversationId": conversation.ID, "receiverId": userID, "read": false})
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, ConversationSummary{Conversation: conversation, Peer: *peer, UnreadCount: unread})
	}
	return summaries, nil
}
func (r *ChatRepository) ListMessages(ctx context.Context, conversationID primitive.ObjectID, page, limit int) (MessagePage, error) {
	filter := bson.M{"conversationId": conversationID}
	total, err := r.messages.CountDocuments(ctx, filter)
	if err != nil {
		return MessagePage{}, err
	}
	skip := total - int64(page*limit)
	fetchLimit := int64(limit)
	if skip < 0 {
		fetchLimit += skip
		skip = 0
	}
	items := []models.ChatMessage{}
	if fetchLimit > 0 {
		cursor, err := r.messages.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: 1}}).SetSkip(skip).SetLimit(fetchLimit))
		if err != nil {
			return MessagePage{}, err
		}
		defer cursor.Close(ctx)
		if err := cursor.All(ctx, &items); err != nil {
			return MessagePage{}, err
		}
	}
	pages := int((total + int64(limit) - 1) / int64(limit))
	return MessagePage{Items: items, Page: page, Limit: limit, Total: total, TotalPages: pages}, nil
}
func (r *ChatRepository) SendMessage(ctx context.Context, message *models.ChatMessage) error {
	if _, err := r.messages.InsertOne(ctx, message); err != nil {
		return err
	}
	_, err := r.conversations.UpdateOne(ctx, bson.M{"_id": message.ConversationID}, bson.M{"$set": bson.M{"lastMessage": message.Content, "lastMessageAt": message.CreatedAt, "lastSenderId": message.SenderID, "updatedAt": message.CreatedAt}})
	return err
}
func (r *ChatRepository) MarkRead(ctx context.Context, conversationID, userID primitive.ObjectID) error {
	_, err := r.messages.UpdateMany(ctx, bson.M{"conversationId": conversationID, "receiverId": userID, "read": false}, bson.M{"$set": bson.M{"read": true}})
	return err
}
func (r *ChatRepository) UnreadCount(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	return r.messages.CountDocuments(ctx, bson.M{"receiverId": userID, "read": false})
}
