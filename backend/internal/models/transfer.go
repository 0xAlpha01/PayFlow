package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Transfer struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SenderID        primitive.ObjectID `bson:"senderId" json:"senderId"`
	ReceiverID      primitive.ObjectID `bson:"receiverId" json:"receiverId"`
	Amount          float64            `bson:"amount" json:"amount"`
	Currency        string             `bson:"currency" json:"currency"`
	Status          string             `bson:"status" json:"status"`
	Note            string             `bson:"note,omitempty" json:"note,omitempty"`
	ClientReference string             `bson:"clientReference,omitempty" json:"clientReference,omitempty"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
}
