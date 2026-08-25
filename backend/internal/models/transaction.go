package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Transaction struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID            primitive.ObjectID `bson:"userId" json:"-"`
	CustomerID        primitive.ObjectID `bson:"customerId,omitempty" json:"customerId,omitempty"`
	Amount            float64            `bson:"amount" json:"amount"`
	Currency          string             `bson:"currency" json:"currency"`
	PaymentMethod     string             `bson:"paymentMethod" json:"paymentMethod"`
	Status            string             `bson:"status" json:"status"`
	PaystackReference string             `bson:"paystackReference,omitempty" json:"paystackReference,omitempty"`
	CreatedAt         time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updatedAt" json:"updatedAt"`
}
