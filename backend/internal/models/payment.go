package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Payment struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID            primitive.ObjectID `bson:"userId" json:"-"`
	CustomerID        primitive.ObjectID `bson:"customerId,omitempty" json:"customerId,omitempty"`
	TransactionID     primitive.ObjectID `bson:"transactionId,omitempty" json:"transactionId,omitempty"`
	Amount            float64            `bson:"amount" json:"amount"`
	Currency          string             `bson:"currency" json:"currency"`
	Status            string             `bson:"status" json:"status"`
	PaystackReference string             `bson:"paystackReference,omitempty" json:"paystackReference,omitempty"`
	CreatedAt         time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updatedAt" json:"updatedAt"`
}
