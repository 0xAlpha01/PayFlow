package repositories

import (
	"context"
	"errors"
	"github.com/ngwebwork/payflow/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type UserRepository struct{ collection *mongo.Collection }

func NewUserRepository(db *mongo.Database) (*UserRepository, error) {
	r := &UserRepository{collection: db.Collection("users")}
	_, err := r.collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)}, {Keys: bson.D{{Key: "phone", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)}})
	return r, err
}
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}
func (r *UserRepository) ByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": strings.ToLower(email)}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &user, err
}
func (r *UserRepository) ByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &user, err
}
func (r *UserRepository) ByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &user, err
}
func (r *UserRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now().UTC()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}
