package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/ngwebwork/payflow/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrInsufficientBalance = errors.New("insufficient balance")

type TransferPage struct {
	Items      []models.Transfer `json:"items"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Total      int64             `json:"total"`
	TotalPages int               `json:"totalPages"`
}
type TransferRepository struct {
	client    *mongo.Client
	users     *mongo.Collection
	transfers *mongo.Collection
}

func NewTransferRepository(client *mongo.Client, db *mongo.Database) (*TransferRepository, error) {
	r := &TransferRepository{client: client, users: db.Collection("users"), transfers: db.Collection("transfers")}
	_, err := r.transfers.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "senderId", Value: 1}, {Key: "createdAt", Value: -1}}},
		{Keys: bson.D{{Key: "receiverId", Value: 1}, {Key: "createdAt", Value: -1}}},
		{Keys: bson.D{{Key: "senderId", Value: 1}, {Key: "clientReference", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)},
	})
	return r, err
}

// ByClientReference supports idempotent retries: a repeated submission with the same
// clientReference from the same sender returns the transfer already created for it.
func (r *TransferRepository) ByClientReference(ctx context.Context, senderID primitive.ObjectID, clientReference string) (*models.Transfer, error) {
	if clientReference == "" {
		return nil, nil
	}
	var t models.Transfer
	err := r.transfers.FindOne(ctx, bson.M{"senderId": senderID, "clientReference": clientReference}).Decode(&t)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &t, err
}

// Balance reads a user's current wallet balance.
func (r *TransferRepository) Balance(ctx context.Context, userID primitive.ObjectID) (float64, error) {
	var doc struct {
		Balance float64 `bson:"balance"`
	}
	if err := r.users.FindOne(ctx, bson.M{"_id": userID}, options.FindOne().SetProjection(bson.M{"balance": 1})).Decode(&doc); err != nil {
		return 0, err
	}
	return doc.Balance, nil
}

// Execute atomically debits the sender, credits the receiver, and records the transfer.
// The sender-side update filters on balance >= amount so concurrent transfers can never
// push a balance negative, and the whole operation commits or rolls back together.
func (r *TransferRepository) Execute(ctx context.Context, transfer *models.Transfer) error {
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (any, error) {
		debit, err := r.users.UpdateOne(sc, bson.M{"_id": transfer.SenderID, "balance": bson.M{"$gte": transfer.Amount}}, bson.M{"$set": bson.M{"updatedAt": time.Now().UTC()}, "$inc": bson.M{"balance": -transfer.Amount}})
		if err != nil {
			return nil, err
		}
		if debit.MatchedCount == 0 {
			return nil, ErrInsufficientBalance
		}
		if _, err := r.users.UpdateOne(sc, bson.M{"_id": transfer.ReceiverID}, bson.M{"$set": bson.M{"updatedAt": time.Now().UTC()}, "$inc": bson.M{"balance": transfer.Amount}}); err != nil {
			return nil, err
		}
		if _, err := r.transfers.InsertOne(sc, transfer); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}
func (r *TransferRepository) List(ctx context.Context, userID primitive.ObjectID, page, limit int) (TransferPage, error) {
	filter := bson.M{"$or": bson.A{bson.M{"senderId": userID}, bson.M{"receiverId": userID}}}
	total, err := r.transfers.CountDocuments(ctx, filter)
	if err != nil {
		return TransferPage{}, err
	}
	cursor, err := r.transfers.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)))
	if err != nil {
		return TransferPage{}, err
	}
	defer cursor.Close(ctx)
	items := []models.Transfer{}
	if err := cursor.All(ctx, &items); err != nil {
		return TransferPage{}, err
	}
	pages := int((total + int64(limit) - 1) / int64(limit))
	return TransferPage{Items: items, Page: page, Limit: limit, Total: total, TotalPages: pages}, nil
}
